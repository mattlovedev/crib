package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"mattlove.dev/crib/game"
	"mattlove.dev/crib/strategy"
	"mattlove.dev/crib/strategy/discard"
	"mattlove.dev/crib/strategy/peg"
)

const humanIdx = 0
const aiIdx = 1
const winScore = 121

type Phase string

const (
	PhaseDiscard  Phase = "discard"
	PhasePeg      Phase = "peg"
	PhaseScore    Phase = "score"
	PhaseGameOver Phase = "game_over"
)

type Session struct {
	rng *rand.Rand
	ai  strategy.Player

	Phase  Phase
	Dealer int    // 0=human, 1=ai
	Scores [2]int // [human, ai]

	HumanHand game.Cards // 6 during discard, 4 after
	AIHand    game.Cards // 4 cards after discard (hidden until score)

	Crib   game.Cards
	Cut    game.Card
	CutSet bool

	HumanInHand game.Cards
	AIInHand    game.Cards
	PegSeries   game.Cards
	PegCount    int
	PegCurrent  int // 0=human, 1=ai
	LastToPlay  int // -1=none
	HumanPlayed game.Cards
	AIPlayed    game.Cards

	HumanHandScore int
	AIHandScore    int
	CribScore      int

	Log    []string
	Winner int // -1=none, 0=human, 1=ai
}

type StateView struct {
	Phase       string `json:"phase"`
	Dealer      int    `json:"dealer"`
	Scores      [2]int `json:"scores"`
	HumanHand   []int  `json:"human_hand"`
	AIHandCount int    `json:"ai_hand_count"`
	AIHand      []int  `json:"ai_hand"`
	Crib        []int  `json:"crib"`
	Cut         int    `json:"cut"` // -1 if not revealed yet

	PegCount    int    `json:"peg_count"`
	PegSeries   []int  `json:"peg_series"`
	HumanPlayed []int  `json:"human_played"`
	AIPlayed    []int  `json:"ai_played"`
	WhoseTurn   string `json:"whose_turn"`
	CanPlay     bool   `json:"can_play"`

	HumanHandScore int `json:"human_hand_score"`
	AIHandScore    int `json:"ai_hand_score"`
	CribScore      int `json:"crib_score"`

	Log    []string `json:"log"`
	Winner string   `json:"winner"`
	Error  string   `json:"error,omitempty"`
}

var (
	globalSession *Session
	sessionMu     sync.Mutex
	aiPlayer      strategy.Player
)

func ids(cards game.Cards) []int {
	if cards == nil {
		return []int{}
	}
	out := make([]int, len(cards))
	for i, c := range cards {
		out[i] = c.Id
	}
	return out
}

func (s *Session) view() StateView {
	v := StateView{
		Phase:       string(s.Phase),
		Dealer:      s.Dealer,
		Scores:      s.Scores,
		HumanHand:   ids(s.HumanHand),
		AIHandCount: len(s.AIHand),
		AIHand:      []int{},
		Crib:        []int{},
		PegCount:    s.PegCount,
		PegSeries:   ids(s.PegSeries),
		HumanPlayed: ids(s.HumanPlayed),
		AIPlayed:    ids(s.AIPlayed),
		Log:         s.Log,
		Cut:         -1,
	}
	if s.Log == nil {
		v.Log = []string{}
	}
	if s.PegSeries == nil {
		v.PegSeries = []int{}
	}

	if s.CutSet {
		v.Cut = s.Cut.Id
	}

	if s.Phase == PhaseScore || s.Phase == PhaseGameOver {
		v.AIHand = ids(s.AIHand)
		v.Crib = ids(s.Crib)
		v.HumanHandScore = s.HumanHandScore
		v.AIHandScore = s.AIHandScore
		v.CribScore = s.CribScore
	}

	if s.Phase == PhasePeg {
		// Show remaining pegging cards, not full kept hand
		v.HumanHand = ids(s.HumanInHand)
		v.AIHandCount = len(s.AIInHand)
		v.WhoseTurn = "human"
		if s.PegCurrent == aiIdx {
			v.WhoseTurn = "ai"
		}
		v.CanPlay = canPlayAny(s.HumanInHand, s.PegCount)
	}

	switch s.Winner {
	case humanIdx:
		v.Winner = "human"
	case aiIdx:
		v.Winner = "ai"
	default:
		v.Winner = ""
	}

	return v
}

func (s *Session) dealRound() {
	deck := game.NewDeck()
	s.rng.Shuffle(len(deck.Cards), func(i, j int) {
		deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
	})
	s.HumanHand = deck.Cards[0:6].Copy()
	s.AIHand = deck.Cards[6:12].Copy()
	s.Cut = deck.Cards[12]
	s.CutSet = false
	s.Crib = nil
	s.HumanPlayed = nil
	s.AIPlayed = nil
	s.HumanInHand = nil
	s.AIInHand = nil
	s.PegSeries = nil
	s.PegCount = 0
	s.LastToPlay = -1
	s.Log = nil
	s.Phase = PhaseDiscard
	s.HumanHandScore = 0
	s.AIHandScore = 0
	s.CribScore = 0
}

func (s *Session) pegPoints(player, pts int, reason string) bool {
	s.Scores[player] += pts
	name := "You"
	if player == aiIdx {
		name = "AI"
	}
	msg := fmt.Sprintf("%s pegged %d (%s) → %d", name, pts, reason, s.Scores[player])
	s.Log = append(s.Log, msg)
	if s.Scores[player] >= winScore {
		s.Winner = player
		s.Phase = PhaseGameOver
		return true
	}
	return false
}

func (s *Session) aiPlayOnePeg() bool {
	card := s.ai.Play(s.AIInHand, strategy.PeggingState{Count: s.PegCount, Series: s.PegSeries})
	s.AIInHand = removeCard(s.AIInHand, card)
	s.AIPlayed = append(s.AIPlayed, card)
	s.PegSeries = append(s.PegSeries, card)
	s.PegCount += card.Value
	s.LastToPlay = aiIdx

	pts := game.ScorePeggingPlay(s.PegCount, s.PegSeries)
	msg := fmt.Sprintf("AI plays %s (count: %d)", card.String(), s.PegCount)
	if pts > 0 {
		msg += fmt.Sprintf(" +%d", pts)
	}
	s.Log = append(s.Log, msg)

	if pts > 0 {
		if s.pegPoints(aiIdx, pts, scoreReason(s.PegCount)) {
			return true
		}
	}

	if s.PegCount == 31 {
		s.PegSeries = nil
		s.PegCount = 0
		s.LastToPlay = -1
		s.PegCurrent = humanIdx
	} else {
		s.PegCurrent = humanIdx
	}
	return false
}

func scoreReason(count int) string {
	if count == 15 {
		return "fifteen"
	}
	if count == 31 {
		return "31"
	}
	return "peg"
}

// advancePegging runs AI turns until it's human's turn or the round ends.
func (s *Session) advancePegging() {
	for len(s.HumanInHand) > 0 || len(s.AIInHand) > 0 {
		humanCanPlay := canPlayAny(s.HumanInHand, s.PegCount)
		aiCanPlay := canPlayAny(s.AIInHand, s.PegCount)

		if !humanCanPlay && !aiCanPlay {
			if s.LastToPlay >= 0 && s.PegCount != 31 {
				if s.pegPoints(s.LastToPlay, 1, "go") {
					return
				}
			}
			s.PegSeries = nil
			s.PegCount = 0
			if s.LastToPlay >= 0 {
				s.PegCurrent = 1 - s.LastToPlay
			}
			s.LastToPlay = -1
			continue
		}

		if s.PegCurrent == humanIdx && !humanCanPlay {
			s.PegCurrent = aiIdx
		} else if s.PegCurrent == aiIdx && !aiCanPlay {
			s.PegCurrent = humanIdx
		}

		if s.PegCurrent == humanIdx {
			return
		}

		if s.aiPlayOnePeg() {
			return
		}
	}

	// Last card: 1 point to whoever played it (a 31 already scored 2 and reset LastToPlay).
	if s.LastToPlay >= 0 && s.PegCount != 31 {
		if s.pegPoints(s.LastToPlay, 1, "last card") {
			return
		}
	}

	s.enterScoringPhase()
}

func (s *Session) score(player, pts int, label string) bool {
	s.Scores[player] += pts
	name := "You"
	if player == aiIdx {
		name = "AI"
	}
	s.Log = append(s.Log, fmt.Sprintf("%s %s: %d pts (score: %d)", name, label, pts, s.Scores[player]))
	if s.Scores[player] >= winScore {
		s.Winner = player
		s.Phase = PhaseGameOver
		return true
	}
	return false
}

func (s *Session) enterScoringPhase() {
	humanPts := game.CountCards(s.HumanHand, &s.Cut, false)
	aiPts := game.CountCards(s.AIHand, &s.Cut, false)
	cribPts := game.CountCards(s.Crib, &s.Cut, true)

	s.HumanHandScore = humanPts
	s.AIHandScore = aiPts
	s.CribScore = cribPts

	nonDealer := 1 - s.Dealer
	if nonDealer == humanIdx {
		if s.score(humanIdx, humanPts, "hand") {
			return
		}
		if s.score(aiIdx, aiPts, "hand") {
			return
		}
	} else {
		if s.score(aiIdx, aiPts, "hand") {
			return
		}
		if s.score(humanIdx, humanPts, "hand") {
			return
		}
	}
	s.score(s.Dealer, cribPts, "crib")
	if s.Phase != PhaseGameOver {
		s.Phase = PhaseScore
	}
}

func canPlayAny(hand game.Cards, count int) bool {
	for _, c := range hand {
		if count+c.Value <= 31 {
			return true
		}
	}
	return false
}

func removeCard(hand game.Cards, card game.Card) game.Cards {
	out := make(game.Cards, 0, len(hand)-1)
	removed := false
	for _, c := range hand {
		if !removed && c.Id == card.Id {
			removed = true
			continue
		}
		out = append(out, c)
	}
	return out
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func errView(msg string) StateView {
	return StateView{Error: msg, HumanHand: []int{}, AIHand: []int{}, Crib: []int{}, PegSeries: []int{}, HumanPlayed: []int{}, AIPlayed: []int{}, Log: []string{}}
}

func handleNew(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	sessionMu.Lock()
	defer sessionMu.Unlock()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := &Session{
		rng:    rng,
		ai:     aiPlayer,
		Dealer: rng.Intn(2),
		Winner: -1,
	}
	s.dealRound()
	globalSession = s
	writeJSON(w, s.view())
}

func handleDiscard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	sessionMu.Lock()
	defer sessionMu.Unlock()

	if globalSession == nil {
		writeJSON(w, errView("no active game"))
		return
	}
	s := globalSession

	var req struct {
		Cards [2]int `json:"cards"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, errView("bad request"))
		return
	}

	if s.Phase != PhaseDiscard {
		writeJSON(w, errView("not in discard phase"))
		return
	}

	// Validate cards are in human hand
	var discardCards game.Cards
	hand := s.HumanHand
	for _, id := range req.Cards {
		found := false
		for _, c := range hand {
			if c.Id == id {
				found = true
				discardCards = append(discardCards, c)
				break
			}
		}
		if !found {
			writeJSON(w, errView(fmt.Sprintf("card %d not in hand", id)))
			return
		}
	}
	if discardCards[0].Id == discardCards[1].Id {
		writeJSON(w, errView("cannot discard the same card twice"))
		return
	}

	// Human keeps the other 4
	var kept game.Cards
	for _, c := range hand {
		if c.Id != discardCards[0].Id && c.Id != discardCards[1].Id {
			kept = append(kept, c)
		}
	}
	s.HumanHand = kept

	// AI discards
	aiKept, aiDiscard := s.ai.Discard(s.AIHand, s.Dealer == aiIdx)
	s.AIHand = aiKept

	// Build crib
	s.Crib = append(discardCards, aiDiscard...)
	s.CutSet = true

	// Nibs
	if s.Cut.Face == game.Jack {
		s.Log = append(s.Log, fmt.Sprintf("Nibs! %s cut — dealer gets 2", s.Cut.String()))
		if s.Scores[s.Dealer]+2 >= winScore {
			s.Scores[s.Dealer] += 2
			s.Winner = s.Dealer
			s.Phase = PhaseGameOver
			writeJSON(w, s.view())
			return
		}
		s.Scores[s.Dealer] += 2
	}

	// Setup pegging
	s.HumanInHand = s.HumanHand.Copy()
	s.AIInHand = s.AIHand.Copy()
	s.HumanPlayed = game.Cards{}
	s.AIPlayed = game.Cards{}
	s.PegSeries = nil
	s.PegCount = 0
	s.LastToPlay = -1
	s.PegCurrent = 1 - s.Dealer // non-dealer leads

	s.Phase = PhasePeg

	// If AI leads, let it play
	if s.PegCurrent == aiIdx {
		s.advancePegging()
	}

	writeJSON(w, s.view())
}

func handlePeg(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	sessionMu.Lock()
	defer sessionMu.Unlock()

	if globalSession == nil {
		writeJSON(w, errView("no active game"))
		return
	}
	s := globalSession

	if s.Phase != PhasePeg {
		writeJSON(w, errView("not in peg phase"))
		return
	}

	var req struct {
		Card int `json:"card"` // -1 = go
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, errView("bad request"))
		return
	}

	if req.Card == -1 {
		// Human says go
		if canPlayAny(s.HumanInHand, s.PegCount) {
			writeJSON(w, errView("you have a legal play"))
			return
		}
		s.PegCurrent = aiIdx
		s.advancePegging()
		writeJSON(w, s.view())
		return
	}

	// Find card in human's remaining hand
	var played game.Card
	found := false
	for _, c := range s.HumanInHand {
		if c.Id == req.Card {
			played = c
			found = true
			break
		}
	}
	if !found {
		writeJSON(w, errView("card not in hand"))
		return
	}
	if s.PegCount+played.Value > 31 {
		writeJSON(w, errView(fmt.Sprintf("would exceed 31 (count: %d)", s.PegCount)))
		return
	}

	s.HumanInHand = removeCard(s.HumanInHand, played)
	s.HumanPlayed = append(s.HumanPlayed, played)
	s.PegSeries = append(s.PegSeries, played)
	s.PegCount += played.Value
	s.LastToPlay = humanIdx

	pts := game.ScorePeggingPlay(s.PegCount, s.PegSeries)
	msg := fmt.Sprintf("You play %s (count: %d)", played.String(), s.PegCount)
	if pts > 0 {
		msg += fmt.Sprintf(" +%d", pts)
	}
	s.Log = append(s.Log, msg)

	if pts > 0 {
		if s.pegPoints(humanIdx, pts, scoreReason(s.PegCount)) {
			writeJSON(w, s.view())
			return
		}
	}

	if s.PegCount == 31 {
		s.PegSeries = nil
		s.PegCount = 0
		s.LastToPlay = -1
		s.PegCurrent = aiIdx
	} else {
		s.PegCurrent = aiIdx
	}

	s.advancePegging()
	writeJSON(w, s.view())
}

func handleNext(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	sessionMu.Lock()
	defer sessionMu.Unlock()

	if globalSession == nil {
		writeJSON(w, errView("no active game"))
		return
	}
	s := globalSession

	if s.Phase != PhaseScore {
		writeJSON(w, errView("not in score phase"))
		return
	}

	s.Dealer = 1 - s.Dealer
	s.dealRound()
	writeJSON(w, s.view())
}

func main() {
	start := time.Now()

	fmt.Print("Building hand score cache... ")
	cache := discard.NewSummaryCache()
	fmt.Printf("done (%s)\n", time.Since(start).Round(time.Millisecond))

	fmt.Print("Building crib score cache... ")
	cribCache := discard.NewTwoCribCache()
	fmt.Printf("done (%s)\n", time.Since(start).Round(time.Millisecond))

	d := discard.MaxAvgDiff{Cache: cache, CribCache: cribCache}
	p := peg.MaxSetup{}
	aiPlayer = strategy.NewStrategy("MaxAvgDiff/MaxSetup", d, p)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/new", handleNew)
	mux.HandleFunc("/api/discard", handleDiscard)
	mux.HandleFunc("/api/peg", handlePeg)
	mux.HandleFunc("/api/next", handleNext)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/play/", http.StatusFound)
			return
		}
		http.FileServer(http.Dir(".")).ServeHTTP(w, r)
	})

	addr := ":8080"
	fmt.Printf("Listening on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

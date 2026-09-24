package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"mattlove.dev/crib/engine"
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

	peg         *engine.Pegging // nil until the discard; Hands indexed by humanIdx/aiIdx
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
		PegSeries:   []int{},
		HumanPlayed: ids(s.HumanPlayed),
		AIPlayed:    ids(s.AIPlayed),
		Log:         s.Log,
		Cut:         -1,
	}
	if s.Log == nil {
		v.Log = []string{}
	}
	if s.peg != nil {
		v.PegCount = s.peg.Count
		v.PegSeries = ids(s.peg.Series)
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
		v.HumanHand = ids(s.peg.Hands[humanIdx])
		v.AIHandCount = len(s.peg.Hands[aiIdx])
		v.WhoseTurn = "human"
		if s.peg.Current == aiIdx {
			v.WhoseTurn = "ai"
		}
		v.CanPlay = s.peg.CanPlay(humanIdx)
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
	s.peg = nil
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

// playPeg plays card for the current player, logs it, and scores it.
// Returns true if that ended the game.
func (s *Session) playPeg(card game.Card) bool {
	player := s.peg.Current
	name := "You play"
	if player == aiIdx {
		name = "AI plays"
		s.AIPlayed = append(s.AIPlayed, card)
	} else {
		s.HumanPlayed = append(s.HumanPlayed, card)
	}

	count, pts := s.peg.Play(card)
	msg := fmt.Sprintf("%s %s (count: %d)", name, card.String(), count)
	if pts > 0 {
		msg += fmt.Sprintf(" +%d", pts)
	}
	s.Log = append(s.Log, msg)

	return pts > 0 && s.pegPoints(player, pts, scoreReason(count))
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

// advancePegging applies go/last-card points and runs AI turns until it's the
// human's turn, the game ends, or pegging is over (then scores the round).
func (s *Session) advancePegging() {
	for {
		if a, ok := s.peg.Resolve(); ok {
			if s.pegPoints(a.Player, a.Points, a.Kind.String()) {
				return
			}
			continue
		}
		if s.peg.Done() {
			break
		}
		if s.peg.Current == humanIdx {
			return
		}
		card := s.ai.Play(s.peg.Hands[aiIdx], s.peg.State())
		if s.playPeg(card) {
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

	// Setup pegging; the non-dealer leads
	s.peg = engine.NewPegging([2]game.Cards{humanIdx: s.HumanHand, aiIdx: s.AIHand}, s.Dealer)
	s.HumanPlayed = game.Cards{}
	s.AIPlayed = game.Cards{}

	s.Phase = PhasePeg

	// If AI leads, let it play
	if s.peg.Current == aiIdx {
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

	if s.peg.Current != humanIdx {
		writeJSON(w, errView("not your turn"))
		return
	}

	if req.Card == -1 {
		// Human says go
		if s.peg.CanPlay(humanIdx) {
			writeJSON(w, errView("you have a legal play"))
			return
		}
		s.advancePegging()
		writeJSON(w, s.view())
		return
	}

	// Find card in human's remaining hand
	var played game.Card
	found := false
	for _, c := range s.peg.Hands[humanIdx] {
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
	if s.peg.Count+played.Value > 31 {
		writeJSON(w, errView(fmt.Sprintf("would exceed 31 (count: %d)", s.peg.Count)))
		return
	}

	if !s.playPeg(played) {
		s.advancePegging()
	}
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

package engine

import (
	"mattlove.dev/crib/game"
	"mattlove.dev/crib/strategy"
)

// Pegging holds the state of one pegging phase and applies its rules. It is shared by
// the simulator (doPegging) and the interactive game server (play/), which drive it
// the same way: call Resolve until it reports no award, stop if Done, otherwise have
// the Current player choose a card and pass it to Play.
type Pegging struct {
	Hands      [2]game.Cards // cards each player has not yet played
	Series     game.Cards    // cards played since the last reset, in order
	Count      int           // running count of Series
	Current    int           // player whose turn it is
	LastToPlay int           // player who played the last card of Series, -1 if none
}

// AwardKind is why Resolve awarded a point.
type AwardKind int

const (
	Go       AwardKind = iota // neither player can play: 1 to the last to play
	LastCard                  // final card of pegging (not a 31): 1 to whoever played it
)

func (k AwardKind) String() string {
	if k == LastCard {
		return "last card"
	}
	return "go"
}

// Award is a point given by Resolve rather than for playing a card.
type Award struct {
	Player int
	Points int
	Kind   AwardKind
}

// NewPegging starts pegging with each player's kept hand; the non-dealer leads.
func NewPegging(hands [2]game.Cards, dealer int) *Pegging {
	return &Pegging{
		Hands:      [2]game.Cards{hands[0].Copy(), hands[1].Copy()},
		Current:    1 - dealer,
		LastToPlay: -1,
	}
}

// CanPlay reports whether player has a card that keeps the count at or below 31.
func (p *Pegging) CanPlay(player int) bool {
	for _, c := range p.Hands[player] {
		if p.Count+c.Value <= 31 {
			return true
		}
	}
	return false
}

// Done reports whether both players have played all their cards.
func (p *Pegging) Done() bool {
	return len(p.Hands[0]) == 0 && len(p.Hands[1]) == 0
}

// State is what the current player's Pegger sees.
func (p *Pegging) State() strategy.PeggingState {
	return strategy.PeggingState{Count: p.Count, Series: p.Series}
}

// Resolve applies everything that happens without a decision: passing the turn when
// the current player can't play, the go point and reset when neither can, and the
// last-card point once both hands are empty. It stops at each award so callers can
// check for a win, and returns ok=false once the Current player must play or pegging
// is Done.
func (p *Pegging) Resolve() (a Award, ok bool) {
	for {
		if p.Done() {
			// A 31 already scored 2 and reset LastToPlay, so this is only for other counts.
			if p.LastToPlay >= 0 {
				a = Award{Player: p.LastToPlay, Points: 1, Kind: LastCard}
				p.LastToPlay = -1
				return a, true
			}
			return Award{}, false
		}

		if p.CanPlay(p.Current) {
			return Award{}, false
		}
		if p.CanPlay(1 - p.Current) {
			p.Current = 1 - p.Current
			continue
		}

		// Both stuck: go point to the last player, then the other player leads a new series.
		last := p.LastToPlay
		p.reset()
		if last >= 0 {
			p.Current = 1 - last
			return Award{Player: last, Points: 1, Kind: Go}, true
		}
	}
}

// Play plays card for the Current player and passes the turn. It returns the count the
// card reached and the points it scored (15s, 31, pairs, runs). A 31 resets the count.
// The caller must ensure the card is in the Current player's hand and is legal.
func (p *Pegging) Play(card game.Card) (count, pts int) {
	player := p.Current
	p.Hands[player] = removeCard(p.Hands[player], card)
	p.Series = append(p.Series, card)
	p.Count += card.Value
	p.LastToPlay = player
	count = p.Count
	pts = game.ScorePeggingPlay(p.Count, p.Series)
	if p.Count == 31 {
		p.reset()
	}
	p.Current = 1 - player
	return count, pts
}

func (p *Pegging) reset() {
	p.Series = nil
	p.Count = 0
	p.LastToPlay = -1
}

// doPegging runs the pegging phase, updating result scores.
// Returns true if a player reached winScore during pegging.
func doPegging(players [2]strategy.Player, hands [2]game.Cards, dealer int, result *GameResult) bool {
	peg := NewPegging(hands, dealer)
	for {
		if a, ok := peg.Resolve(); ok {
			if addPegged(a.Player, a.Points, result) {
				return true
			}
			continue
		}
		if peg.Done() {
			return false
		}

		player := peg.Current
		card := players[player].Play(peg.Hands[player], peg.State())
		if _, pts := peg.Play(card); addPegged(player, pts, result) {
			return true
		}
	}
}

func addPegged(player, pts int, result *GameResult) bool {
	result.Scores[player] += pts
	result.PeggedPoints[player] += pts
	return result.Scores[player] >= winScore
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

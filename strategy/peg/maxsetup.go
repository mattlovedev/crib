package peg

import (
	"mattlove.dev/crib/game"
	"mattlove.dev/crib/strategy"
)

// MaxSetup plays the highest-scoring card if points are available (like MaxNext).
// If no card scores immediately, it plays the card that maximizes the best score
// achievable on its next turn — approximating pair setups (play a card you hold
// a duplicate of) and run setups (play adjacent cards you can extend later).
// The opponent's intervening card is not modeled.
type MaxSetup struct{}

func (MaxSetup) Play(hand game.Cards, state strategy.PeggingState) game.Card {
	// Phase 1: take points if available, same as MaxNext.
	bestPts := 0
	var best game.Card
	for _, c := range hand {
		if state.Count+c.Value > 31 {
			continue
		}
		series := append(state.Series[:len(state.Series):len(state.Series)], c)
		pts := game.ScorePeggingPlay(state.Count+c.Value, series)
		if pts > bestPts || (pts == bestPts && c.Value > best.Value) || (pts == bestPts && c.Value == best.Value && c.Id < best.Id) {
			bestPts = pts
			best = c
		}
	}
	if bestPts > 0 {
		return best
	}

	// Phase 2: no immediate points — pick the card with the best setup score.
	// Setup score = max points we could score on our very next play,
	// checking every remaining hand card against series+[c, r].
	bestSetup := -1
	for _, c := range hand {
		if state.Count+c.Value > 31 {
			continue
		}
		afterC := append(state.Series[:len(state.Series):len(state.Series)], c)
		countAfterC := state.Count + c.Value

		setup := 0
		for _, r := range hand {
			if r.Id == c.Id || countAfterC+r.Value > 31 {
				continue
			}
			afterR := append(afterC[:len(afterC):len(afterC)], r)
			if pts := game.ScorePeggingPlay(countAfterC+r.Value, afterR); pts > setup {
				setup = pts
			}
		}

		if setup > bestSetup || (setup == bestSetup && c.Value > best.Value) || (setup == bestSetup && c.Value == best.Value && c.Id < best.Id) {
			bestSetup = setup
			best = c
		}
	}
	return best
}

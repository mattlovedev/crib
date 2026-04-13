package peg

import (
	"mattlove.dev/crib/game"
	"mattlove.dev/crib/strategy"
)

// MinValue plays the lowest-value legal card during pegging.
type MinValue struct{}

func (MinValue) Play(hand game.Cards, state strategy.PeggingState) game.Card {
	var best game.Card
	bestSet := false
	for _, c := range hand {
		if state.Count+c.Value <= 31 {
			if !bestSet || c.Value < best.Value || (c.Value == best.Value && c.Id < best.Id) {
				best = c
				bestSet = true
			}
		}
	}
	return best
}

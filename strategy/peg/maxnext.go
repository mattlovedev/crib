package peg

import (
	"mattlove.dev/crib/game"
	"mattlove.dev/crib/strategy"
)

// MaxNext plays the legal card that scores the most points immediately.
// Ties are broken by highest value, then lowest Id.
type MaxNext struct{}

func (MaxNext) Play(hand game.Cards, state strategy.PeggingState) game.Card {
	var best game.Card
	bestPts := -1
	for _, c := range hand {
		if state.Count+c.Value > 31 {
			continue
		}
		series := append(state.Series[:len(state.Series):len(state.Series)], c) // full slice expr: never write into the caller's array
		pts := game.ScorePeggingPlay(state.Count+c.Value, series)
		if pts > bestPts || (pts == bestPts && c.Value > best.Value) || (pts == bestPts && c.Value == best.Value && c.Id < best.Id) {
			best = c
			bestPts = pts
		}
	}
	return best
}

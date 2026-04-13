package discard

import (
	"sort"

	"mattlove.dev/crib/game"
)

// MinValue keeps the 4 lowest-value cards, throwing the 2 highest to the crib.
type MinValue struct{}

func (MinValue) Discard(hand game.Cards, _ bool) (keep, crib game.Cards) {
	sorted := hand.Copy()
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Value != sorted[j].Value {
			return sorted[i].Value < sorted[j].Value
		}
		return sorted[i].Id < sorted[j].Id
	})
	return sorted[:4], sorted[4:]
}

package discard

import (
	"sort"

	"mattlove.dev/crib/game"
)

// MaxValue keeps the 4 highest-value cards, throwing the 2 lowest to the crib.
type MaxValue struct{}

func (MaxValue) Discard(hand game.Cards, _ bool) (keep, crib game.Cards) {
	sorted := hand.Copy()
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Value != sorted[j].Value {
			return sorted[i].Value > sorted[j].Value
		}
		return sorted[i].Id > sorted[j].Id
	})
	return sorted[:4], sorted[4:]
}

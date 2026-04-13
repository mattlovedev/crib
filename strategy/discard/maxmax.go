package discard

import "mattlove.dev/crib/game"

// MaxMax keeps the 4 cards with the highest maximum hand score across all cuts.
type MaxMax struct{ Cache SummaryCache }

func (m MaxMax) Discard(hand game.Cards, _ bool) (keep, crib game.Cards) {
	fours, twos := hand.ChooseFourWithRemaining()
	best := -1
	idx := 0
	for i, four := range fours {
		if v := m.Cache.lookup(four).Max; v > best {
			best = v
			idx = i
		}
	}
	return fours[idx], twos[idx]
}

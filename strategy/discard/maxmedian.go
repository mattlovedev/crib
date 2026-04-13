package discard

import "mattlove.dev/crib/game"

// MaxMedian keeps the 4 cards with the highest median hand score across all cuts.
type MaxMedian struct{ Cache SummaryCache }

func (m MaxMedian) Discard(hand game.Cards, _ bool) (keep, crib game.Cards) {
	fours, twos := hand.ChooseFourWithRemaining()
	best := -1
	idx := 0
	for i, four := range fours {
		if v := m.Cache.lookup(four).Median; v > best {
			best = v
			idx = i
		}
	}
	return fours[idx], twos[idx]
}

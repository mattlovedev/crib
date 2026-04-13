package discard

import "mattlove.dev/crib/game"

// MaxAvg keeps the 4 cards with the highest average hand score across all cuts.
type MaxAvg struct{ Cache SummaryCache }

func (m MaxAvg) Discard(hand game.Cards, _ bool) (keep, crib game.Cards) {
	fours, twos := hand.ChooseFourWithRemaining()
	best := -1.0
	idx := 0
	for i, four := range fours {
		if v := m.Cache.lookup(four).Avg; v > best {
			best = v
			idx = i
		}
	}
	return fours[idx], twos[idx]
}

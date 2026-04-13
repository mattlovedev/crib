package discard

import "mattlove.dev/crib/game"

// MaxMode keeps the 4 cards with the highest modal hand score across all cuts.
type MaxMode struct{ Cache SummaryCache }

func (m MaxMode) Discard(hand game.Cards, _ bool) (keep, crib game.Cards) {
	fours, twos := hand.ChooseFourWithRemaining()
	best := -1
	idx := 0
	for i, four := range fours {
		if v := m.Cache.lookup(four).Mode; v > best {
			best = v
			idx = i
		}
	}
	return fours[idx], twos[idx]
}

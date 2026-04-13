package discard

import "mattlove.dev/crib/game"

// MaxMin keeps the 4 cards with the highest minimum hand score across all cuts.
type MaxMin struct{ Cache SummaryCache }

func (m MaxMin) Discard(hand game.Cards, _ bool) (keep, crib game.Cards) {
	fours, twos := hand.ChooseFourWithRemaining()
	best := -1
	idx := 0
	for i, four := range fours {
		if v := m.Cache.lookup(four).Min; v > best {
			best = v
			idx = i
		}
	}
	return fours[idx], twos[idx]
}

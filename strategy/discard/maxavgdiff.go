package discard

import "mattlove.dev/crib/game"

// MaxAvgDiff keeps the 4 cards that maximize avg4 +/- avg2:
//   - avg4: average hand score across all cuts (from SummaryCache)
//   - avg2: average crib score of the 2 discarded cards (from TwoCribCache)
//   - if dealer crib:  score = avg4 + avg2
//   - if opponent crib: score = avg4 - avg2
type MaxAvgDiff struct {
	Cache     SummaryCache
	CribCache TwoCribCache
}

func (m MaxAvgDiff) Discard(hand game.Cards, dealerCrib bool) (keep, crib game.Cards) {
	fours, twos := hand.ChooseFourWithRemaining()

	best := -1e18
	idx := 0
	for i, four := range fours {
		avg4 := m.Cache.lookup(four).Avg
		avg2 := m.CribCache.lookup(twos[i])

		var score float64
		if dealerCrib {
			score = avg4 + avg2
		} else {
			score = avg4 - avg2
		}

		if score > best {
			best = score
			idx = i
		}
	}
	return fours[idx], twos[idx]
}

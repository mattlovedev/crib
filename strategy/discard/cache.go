package discard

import (
	"sync"

	"mattlove.dev/crib/game"
	"mattlove.dev/crib/game/counts"
)

// SummaryCache holds precomputed FourSummary for all 270,725 four-card hands.
// Build it once with NewSummaryCache and pass it to any discard strategy that needs it.
type SummaryCache map[string]counts.FourSummary

func NewSummaryCache() SummaryCache {
	hands := game.NewDeck().Cards.ChoseFour()
	cache := make(SummaryCache, len(hands))
	for _, hand := range hands {
		cache[hand.String()] = counts.MakeSummariesNoCounts(hand)
	}
	return cache
}

func (c SummaryCache) lookup(hand game.Cards) counts.FourSummary {
	hand.Sort()
	return c[hand.String()]
}

// TwoCribCache holds precomputed average crib scores for all 1,326 two-card combinations.
// Each entry is the expected crib value of those two cards averaged over all possible
// (cut, opp1, opp2) selections from the remaining 50-card deck.
// Build it once with NewTwoCribCache and pass it to MaxAvgDiff.
type TwoCribCache map[string]float64

func NewTwoCribCache() TwoCribCache {
	pairs := game.NewDeck().Cards.ChoseTwo()
	results := make([]float64, len(pairs))

	var wg sync.WaitGroup
	for i, pair := range pairs {
		wg.Add(1)
		i, pair := i, pair
		go func() {
			defer wg.Done()
			remaining := game.RemainingDeck(pair, nil).Cards
			n := len(remaining)
			total := 0
			count := 0
			for a := 0; a < n-2; a++ {
				for b := a + 1; b < n-1; b++ {
					for c := b + 1; c < n; c++ {
						r := [3]game.Card{remaining[a], remaining[b], remaining[c]}
						for rot := 0; rot < 3; rot++ {
							cut := r[rot]
							opp1 := r[(rot+1)%3]
							opp2 := r[(rot+2)%3]
							hole := game.Cards{pair[0], pair[1], opp1, opp2}
							total += game.CountCards(hole, &cut, true)
							count++
						}
					}
				}
			}
			if count > 0 {
				results[i] = float64(total) / float64(count)
			}
		}()
	}
	wg.Wait()

	cache := make(TwoCribCache, len(pairs))
	for i, pair := range pairs {
		pair.Sort()
		cache[pair.String()] = results[i]
	}
	return cache
}

func (c TwoCribCache) lookup(two game.Cards) float64 {
	two.Sort()
	return c[two.String()]
}

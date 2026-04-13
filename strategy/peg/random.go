package peg

import (
	"math/rand"
	"sync"

	"mattlove.dev/crib/game"
	"mattlove.dev/crib/strategy"
)

// Random plays a legal card chosen uniformly at random during pegging.
type Random struct {
	Rng *rand.Rand
	mu  sync.Mutex
}

func (r *Random) Play(hand game.Cards, state strategy.PeggingState) game.Card {
	legal := make(game.Cards, 0, len(hand))
	for _, c := range hand {
		if state.Count+c.Value <= 31 {
			legal = append(legal, c)
		}
	}
	r.mu.Lock()
	card := legal[r.Rng.Intn(len(legal))]
	r.mu.Unlock()
	return card
}

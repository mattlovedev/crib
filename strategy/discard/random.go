package discard

import (
	"math/rand"
	"sync"

	"mattlove.dev/crib/game"
)

// Random discards 2 cards chosen uniformly at random.
type Random struct {
	Rng *rand.Rand
	mu  sync.Mutex
}

func (r *Random) Discard(hand game.Cards, _ bool) (keep, crib game.Cards) {
	shuffled := hand.Copy()
	r.mu.Lock()
	r.Rng.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	r.mu.Unlock()
	return shuffled[:4], shuffled[4:]
}

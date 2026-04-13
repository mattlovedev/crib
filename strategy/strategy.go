package strategy

import "mattlove.dev/crib/game"

// PeggingState contains what a Pegger needs to decide which card to play.
type PeggingState struct {
	Count  int        // current running count (0–30)
	Series game.Cards // cards played since last reset, in order
}

// Discarder decides which 2 cards to throw to the crib.
// hand has 6 cards; keep has 4, crib has 2.
type Discarder interface {
	Discard(hand game.Cards, dealerCrib bool) (keep, crib game.Cards)
}

// Pegger decides which card to play during the pegging phase.
// hand contains only the cards not yet played.
// The engine guarantees at least one legal play exists when Play is called.
type Pegger interface {
	Play(hand game.Cards, state PeggingState) game.Card
}

// Player combines both decision interfaces.
type Player interface {
	Discarder
	Pegger
}

// Strategy composes a Discarder and a Pegger into a Player.
type Strategy struct {
	Name string
	Discarder
	Pegger
}

func NewStrategy(name string, d Discarder, p Pegger) Strategy {
	return Strategy{name, d, p}
}

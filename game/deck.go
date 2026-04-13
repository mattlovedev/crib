package game

import (
	"math/rand"
)

type Deck struct {
	Cards Cards
	Index int
}

func NewDeck() Deck {
	return Deck{cardsGlobal.Copy(), 0}
}

func NewShuffledDeck(seed int64) Deck {
	cards := cardsGlobal.Copy()
	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return Deck{cards, 0}
}

func RemainingDeck(hand Cards, crib Cards) Deck {
	cards := make(Cards, 0, NumCards-len(hand)-len(crib))

	removed := make(map[Card]struct{}, cap(cards))
	for _, card := range hand {
		removed[card] = struct{}{}
	}
	for _, card := range crib {
		removed[card] = struct{}{}
	}

	for _, card := range cardsGlobal {
		if _, found := removed[card]; !found {
			cards = append(cards, card)
		}
	}

	return Deck{
		Cards: cards,
		Index: 0,
	}
}

/*func (d *Deck) DealCard() Card {
	card := d.Cards[d.Index]
	d.Index++
	return card
}*/

func (d Deck) HasCards() bool {
	return d.Index < len(d.Cards)
}

func (d Deck) RemainingLen() int {
	return len(d.Cards) - d.Index
}

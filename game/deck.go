package game

type Deck struct {
	Cards Cards
}

func NewDeck() Deck {
	return Deck{cardsGlobal.Copy()}
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

	return Deck{Cards: cards}
}

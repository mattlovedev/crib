package game

type Deck struct {
	Cards Cards
	Index int
}

func NewDeck() Deck {
	return Deck{cardsGlobal.Copy(), 0}
}

func RemainingDeck(removed Cards) Deck {
	cards := make(Cards, NumCards-len(removed))
	index := 0

	for i := range cardsGlobal {
		if !removed.Contains(cardsGlobal[i]) {
			cards[index] = cardsGlobal[i]
			index++
		}
	}

	return Deck{
		Cards: cards,
		Index: 0,
	}
}

func (d *Deck) DealCard() Card {
	card := d.Cards[d.Index]
	d.Index++
	return card
}

func (d Deck) HasCards() bool {
	return d.Index < len(d.Cards)
}

func (d Deck) RemainingLen() int {
	return len(d.Cards) - d.Index
}

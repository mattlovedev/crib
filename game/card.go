package game

import (
	"fmt"
	"sort"
	"strings"

	"mattlove.dev/crib/game/math"
)

const (
	NumFaces = 13
	NumSuits = 4
	NumCards = NumFaces * NumSuits

	FaceIds = "a23456789tjqk"
	SuitIds = "cdhs"
)

const (
	Ace = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	Clubs = iota
	Diamonds
	Hearts
	Spades
)

const (
	AceOfClubs = iota
	AceOfDiamonds
	AceOfHearts
	AceOfSpades

	TwoOfClubs
	TwoOfDiamonds
	TwoOfHearts
	TwoOfSpades

	ThreeOfClubs
	ThreeOfDiamonds
	ThreeOfHearts
	ThreeOfSpades

	FourOfClubs
	FourOfDiamonds
	FourOfHearts
	FourOfSpades

	FiveOfClubs
	FiveOfDiamonds
	FiveOfHearts
	FiveOfSpades

	SixOfClubs
	SixOfDiamonds
	SixOfHearts
	SixOfSpades

	SevenOfClubs
	SevenOfDiamonds
	SevenOfHearts
	SevenOfSpades

	EightOfClubs
	EightOfDiamonds
	EightOfHearts
	EightOfSpades

	NineOfClubs
	NineOfDiamonds
	NineOfHearts
	NineOfSpades

	TenOfClubs
	TenOfDiamonds
	TenOfHearts
	TenOfSpades

	JackOfClubs
	JackOfDiamonds
	JackOfHearts
	JackOfSpades

	QueenOfClubs
	QueenOfDiamonds
	QueenOfHearts
	QueenOfSpades

	KingOfClubs
	KingOfDiamonds
	KingOfHearts
	KingOfSpades
)

var (
	faces       = []string{"ace", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "jack", "queen", "king"}
	suits       = []string{"hearts", "clubs", "diamonds", "spades"}
	cardsGlobal = make(Cards, NumCards)
)

func init() {
	newCard := func(id int) Card {
		face := id / NumSuits
		suit := id % NumSuits
		//face := id % NumFaces
		//suit := id / NumFaces
		value := face + 1
		if value > 10 {
			value = 10
		}
		return Card{id, face, suit, value}
	}

	for i := range cardsGlobal {
		cardsGlobal[i] = newCard(i)
	}
}

type Card struct {
	Id    int
	Face  int
	Suit  int
	Value int
}

func (c Card) String() string {
	return fmt.Sprintf("%c%c", FaceIds[c.Face], SuitIds[c.Suit])
}

func CardById(id int) Card {
	return cardsGlobal[id]
}

func CardByIdString(id string) Card {
	face := strings.IndexByte(FaceIds, id[0])
	suit := strings.IndexByte(SuitIds, id[1])
	return CardByFaceSuit(face, suit)
}

func CardByFaceSuit(face int, suit int) Card {
	id := face*NumSuits + suit
	//id := suit*NumFaces + face
	return cardsGlobal[id]
}

type Cards []Card

func SortedCards(cards ...Card) Cards {
	return Cards(cards).Sort()
}

func (c Cards) Contains(card Card) bool {
	for i := range c {
		if c[i].Id == card.Id {
			return true
		}
	}
	return false
}

func (c Cards) Copy() Cards {
	cards := make(Cards, len(c), cap(c))
	copy(cards, c)
	return cards
}

// TBD this might sort c as well
func (c Cards) Sort() Cards {
	sort.Slice(c, func(i int, j int) bool {
		return c[i].Id < c[j].Id
	})
	return c
}

func (c Cards) String() string {
	// changing sort order for FE
	sort.Slice(c, func(i int, j int) bool {
		if c[i].Suit != c[j].Suit {
			return c[i].Suit < c[j].Suit
		}
		return c[i].Face < c[j].Face
	})
	var sb strings.Builder
	for _, card := range c {
		sb.WriteString(card.String())
	}
	return sb.String()
}

func (c Cards) ChooseTwo() []Cards {
	pairs := make([]Cards, 0, math.NChooseR(len(c), 2))
	for i := 0; i < len(c)-1; i++ {
		for j := i + 1; j < len(c); j++ {
			pairs = append(pairs, []Card{c[i], c[j]})
		}
	}
	return pairs
}

func (c Cards) ChoseFour() []Cards {
	sets := make([]Cards, 0, math.NChooseR(len(c), 4))
	for i := 0; i < len(c)-3; i++ {
		for j := i + 1; j < len(c)-2; j++ {
			for k := j + 1; k < len(c)-1; k++ {
				for l := k + 1; l < len(c); l++ {
					sets = append(sets, Cards{c[i], c[j], c[k], c[l]})
				}
			}
		}
	}
	return sets
}

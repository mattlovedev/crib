package main

import (
	"fmt"
	"os"

	"mattlove.dev/crib/game"
	"mattlove.dev/crib/game/counts"
)

func countTwo(cards game.Cards) {
	/*d := game.RemainingDeck(cards)
	p := d.AllPossiblePairs()
	fmt.Println(p)*/
}

func countFour(cards game.Cards) {
	cards = cards.Sort()
	/*if c, err := counts.ReadAllFourHandCutCounts(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(c[cards.String()])
	}*/
	//fmt.Println(counts.MakeFourScores(cards))
	/*if s, err := counts.ReadAllFourHandSummaries(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(s[cards.String()])
	}*/
	fmt.Println(counts.MakeSummaries(cards))
}

func countFive(cards game.Cards) {
	//h := game.HandFromFiveCards(cards[0], cards[1], cards[2], cards[3], cards[4])
	//fmt.Println(h.Id, h.Count, h.CountWithCut, h.CountWithCrib)

	//h := game.HandFromFourCards(cards[0], cards[1], cards[2], cards[3]).WithCut(cards[4])
	//fmt.Println(h.Id, h.Count, h.CountWithCut, h.CountWithCrib)
}

func countSix(cards game.Cards) {
	/*sets := cards.ChoseFour()
	counts := make(map[string]int, len(sets))
	for _, set := range sets {
		counts[set.String()] = game.CountCards(set, nil, false)
	}
	fmt.Println(counts)*/
}

func main() {
	/*
		CardById(ACE_OF_SPADES)
		CardByFaceSuit(ACE, SPADES)
		CardByIdString("as")
	*/

	args := os.Args[1:]

	cards := make(game.Cards, len(args))
	for i := range args {
		cards[i] = game.CardByIdString(args[i])
	}

	switch len(args) {
	case 2:
		countTwo(cards)
	case 4:
		countFour(cards)
	case 5:
		countFive(cards)
	case 6:
		countSix(cards)
	default:
		fmt.Printf("Usage: %s <c1> <c2> <c3> <c4> [<c5> [<c6]]\n", os.Args[0])
	}

}

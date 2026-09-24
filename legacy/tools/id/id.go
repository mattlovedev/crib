package main

import (
	"fmt"
	"log"
	"mattlove.dev/crib/game"
	"mattlove.dev/crib/game/math"
	"os"
	"strconv"
)

func main1() {
	/*args := os.Args[1:]
	cards := make(game.Cards, len(args))
	for i := range args {
		cards[i] = game.CardByIdString(args[i])
	}*/

	/*gen := math.NewCombinationGenerator(52, 4)

	for i := 0; i < 2354; i++ {
		gen.Next()
		cards := game.CardsByIds(gen.Combination(nil))
		fmt.Printf("%3d - %s - %s\n", cards.Id(), cards.StringIds(), cards)
	}*/

	//*
	ids := [][]int{
		[]int{0, 1, 2, 3},   // 0
		[]int{0, 1, 2, 4},   // 1
		[]int{0, 1, 2, 51},  // 48
		[]int{0, 1, 3, 4},   // 49
		[]int{0, 1, 50, 51}, // 2352
		[]int{0, 2, 3, 4},   // 2353
		[]int{0, 2, 50, 51}, //
		[]int{0, 3, 4, 5},   //
		/*[]int{0, 49, 50, 51},  // 112944
		[]int{1, 2, 3, 4},     // 112945
		[]int{1, 2, 3, 51},    // 112992
		[]int{1, 2, 4, 5},     // 112993
		[]int{48, 49, 50, 51}, // 270724*/
	}

	for _, i := range ids {
		cards := game.CardsByIds(i)
		fmt.Printf("%3d - %s - %s\n", cards.Id(), cards.StringIds(), cards)
	}
	//*/

}

func function(input ...int) int {
	w, x, y, z := input[0], input[1], input[2], input[3]
	return w + x + y + z
}

func main2() {
	gen := math.NewCombinationGenerator(52, 4)
	i := 0
	for gen.Next() {
		input := gen.Combination(nil)
		output := function(input...)
		if output != i {
			log.Fatalf("expected: %d actual: %d input: %v", i, output, input)
		}
		i++
	}
}

func inputEquals(input []int, cmp []int) bool {
	for i := range input {
		if input[i] != cmp[i] {
			return false
		}
	}
	return true
}

func inputEqualsAny(input []int, cmp [][]int) bool {
	for i := range cmp {
		if inputEquals(input, cmp[i]) {
			return true
		}
	}
	return false
}

var test = [][]int{
	{0, 1, 2, 3},
	{0, 1, 2, 51},

	{0, 2, 3, 4},
	{0, 2, 3, 51},

	{0, 2, 49, 50},
	{0, 2, 49, 51},
	{0, 2, 50, 51},

	{0, 48, 49, 50},
	{0, 48, 49, 51},
	{0, 48, 50, 51},
	{0, 49, 50, 51},

	{1, 2, 3, 4},
}

func main3() {
	gen := math.NewCombinationGenerator(52, 4)
	i := 0
	for gen.Next() {
		input := gen.Combination(nil)
		if inputEqualsAny(input, test) {
			fmt.Printf("input: %v id: %d\n", input, i)
		}
		i++
	}
}

func main() {
	h, _ := strconv.Atoi(os.Args[1])
	c, _ := strconv.Atoi(os.Args[2])
	hI := math.IndexToCombination(nil, h, 52, 4)
	cI := math.IndexToCombination(nil, c, 52, 2)
	hand := game.CardsByIds(hI)
	crib := game.CardsByIds(cI)
	fmt.Printf("hand: %s, crib: %s\n", hand, crib)
}

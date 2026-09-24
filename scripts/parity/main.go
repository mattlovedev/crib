// Command parity prints Go's hand statistics so check.js can compare them against the
// JavaScript scoring in web/js/cards.js, which re-implements the same logic for the site.
//
//	go run ./scripts/parity | node scripts/parity/check.js
//
// Output is one line per summary:
//
//	4 <4 card ids> <stats>                       every 4-card hand
//	6 <6 card ids> <kept 4 ids> <crib 2 ids> <stats>   all 15 discards of a fixed sample of 6-card hands
//
// where <stats> is: avg min median max mode modeP belowAvg aboveAvg stdDev.
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"mattlove.dev/crib/game"
	"mattlove.dev/crib/game/counts"
)

const sixHandSamples = 2000

func main() {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for _, hand := range game.NewDeck().Cards.ChooseFour() {
		fmt.Fprintf(w, "4 %s %s\n", ids(hand), stats(counts.MakeSummariesNoCounts(hand)))
	}

	rng := rand.New(rand.NewSource(1))
	for i := 0; i < sixHandSamples; i++ {
		deck := game.NewDeck()
		rng.Shuffle(len(deck.Cards), func(i, j int) {
			deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
		})
		six := deck.Cards[:6].Copy().Sort()
		for _, h := range counts.MakeSixHands(six) {
			fmt.Fprintf(w, "6 %s %s %s %s\n", ids(six), ids(h.Hand), ids(h.Crib), stats(h.Summary))
		}
	}
}

func ids(cards game.Cards) string {
	s := make([]string, len(cards))
	for i, c := range cards {
		s[i] = fmt.Sprint(c.Id)
	}
	return strings.Join(s, " ")
}

func stats(s counts.FourSummary) string {
	return fmt.Sprintf("%.2f %d %d %d %d %.2f %d %d %.2f",
		s.Avg, s.Min, s.Median, s.Max, s.Mode, s.ModeP, s.BelowAvg, s.AboveAvg, s.StdDev)
}

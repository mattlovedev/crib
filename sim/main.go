package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"mattlove.dev/crib/engine"
	"mattlove.dev/crib/strategy"
	"mattlove.dev/crib/strategy/discard"
	"mattlove.dev/crib/strategy/peg"
)

type matchupResult struct {
	wins         [2]int
	peggedPoints [2]int64
	handPoints   [2]int64
	cribPoints   [2]int64
	diff         [2]int64
}

func runMatchup(p0, p1 strategy.Player, n int, rng *rand.Rand) matchupResult {
	var r matchupResult
	for i := 0; i < n; i++ {
		g := engine.RunGame(p0, p1, i%2, rng)
		r.wins[g.Winner]++
		for p := 0; p < 2; p++ {
			r.peggedPoints[p] += int64(g.PeggedPoints[p])
			r.handPoints[p] += int64(g.HandPoints[p])
			r.cribPoints[p] += int64(g.CribPoints[p])
			r.diff[p] += int64(g.Scores[p] - g.Scores[1-p])
		}
	}
	return r
}

func cellLines(r matchupResult, n, p int) []string {
	fn := float64(n)
	total := float64(r.peggedPoints[p]+r.handPoints[p]+r.cribPoints[p]) / fn
	return []string{
		fmt.Sprintf("W: %d (%.1f%%)", r.wins[p], float64(r.wins[p])/fn*100),
		fmt.Sprintf("Pg: %.2f", float64(r.peggedPoints[p])/fn),
		fmt.Sprintf("Hd: %.2f", float64(r.handPoints[p])/fn),
		fmt.Sprintf("Cr: %.2f", float64(r.cribPoints[p])/fn),
		fmt.Sprintf("T:  %.2f", total),
		fmt.Sprintf("D:  %+.2f", float64(r.diff[p])/fn),
	}
}

func allCombinations(rng *rand.Rand, cache discard.SummaryCache, cribCache discard.TwoCribCache) []strategy.Strategy {
	discards := []struct {
		name string
		d    strategy.Discarder
	}{
		// {"MinVal", discard.MinValue{}},
		// {"MaxVal", discard.MaxValue{}},
		// {"Rnd", &discard.Random{Rng: rng}},
		{"MaxAvg", discard.MaxAvg{Cache: cache}},
		{"MaxAvgDiff", discard.MaxAvgDiff{Cache: cache, CribCache: cribCache}},
		// {"MaxMin", discard.MaxMin{Cache: cache}},
		// {"MaxMed", discard.MaxMedian{Cache: cache}},
		// {"MaxMax", discard.MaxMax{Cache: cache}},
		// {"MaxMode", discard.MaxMode{Cache: cache}},
	}
	pegs := []struct {
		name string
		p    strategy.Pegger
	}{
		// {"MinVal", peg.MinValue{}},
		// {"MaxVal", peg.MaxValue{}},
		// {"Rnd", &peg.Random{Rng: rng}},
		{"MaxNext", peg.MaxNext{}},
		{"MaxSetup", peg.MaxSetup{}},
	}

	var all []strategy.Strategy
	for _, d := range discards {
		for _, p := range pegs {
			all = append(all, strategy.NewStrategy(d.name+"D/"+p.name+"P", d.d, p.p))
		}
	}
	return all
}

func main() {
	n := 1000
	if len(os.Args) > 1 {
		var err error
		n, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "usage: sim [num_games]")
			os.Exit(1)
		}
	}

	start := time.Now()
	rng := rand.New(rand.NewSource(start.UnixNano()))

	fmt.Print("Building hand score cache... ")
	cache := discard.NewSummaryCache()
	fmt.Printf("done (%s)\n", time.Since(start).Round(time.Millisecond))

	fmt.Print("Building crib score cache... ")
	cribCache := discard.NewTwoCribCache()
	fmt.Printf("done (%s)\n\n", time.Since(start).Round(time.Millisecond))

	players := allCombinations(rng, cache, cribCache)
	np := len(players)

	// Run upper triangle only; mirror cells reuse the swapped result.
	// Each matchup runs in its own goroutine with a dedicated RNG.
	results := make([][]matchupResult, np)
	for i := range players {
		results[i] = make([]matchupResult, np)
	}
	var wg sync.WaitGroup
	for i := range players {
		for j := i; j < np; j++ {
			wg.Add(1)
			i, j := i, j
			seed := rng.Int63()
			go func() {
				defer wg.Done()
				r := rand.New(rand.NewSource(seed))
				results[i][j] = runMatchup(players[i], players[j], n, r)
			}()
		}
	}
	wg.Wait()

	// Table: cols = P0, rows = P1.
	// Cell (row=r, col=c): P0=players[c], P1=players[r].
	//   c <= r → results[c][r], p=0
	//   c >  r → results[r][c], p=1
	rowLabel := 4
	for _, p := range players {
		if len(p.Name)+2 > rowLabel {
			rowLabel = len(p.Name) + 2
		}
	}
	// cellWidth must fit the widest cell line: the wins line at 100% with n games.
	minCellContent := len(fmt.Sprintf("W: %d (100.0%%)", n)) + 2
	cellWidth := rowLabel + 2
	if minCellContent > cellWidth {
		cellWidth = minCellContent
	}
	numLines := 6

	fmt.Printf("Cols = P0, Rows = P1 — %d games per matchup\n\n", n)

	fmt.Printf("%-*s", rowLabel, "P1 \\ P0")
	for _, col := range players {
		fmt.Printf("%-*s", cellWidth, col.Name)
	}
	fmt.Println()

	divider := ""
	for i := 0; i < rowLabel+cellWidth*np; i++ {
		divider += "-"
	}

	for r, row := range players {
		fmt.Println(divider)
		cells := make([][]string, np)
		for c := range players {
			if c <= r {
				cells[c] = cellLines(results[c][r], n, 0)
			} else {
				cells[c] = cellLines(results[r][c], n, 1)
			}
		}
		for line := 0; line < numLines; line++ {
			if line == 0 {
				fmt.Printf("%-*s", rowLabel, row.Name)
			} else {
				fmt.Printf("%-*s", rowLabel, "")
			}
			for c := range players {
				fmt.Printf("%-*s", cellWidth, cells[c][line])
			}
			fmt.Println()
		}
	}
	fmt.Println(divider)

	// Summary: aggregate each strategy's stats across all matchups, then rank.
	type summary struct {
		name   string
		wins   int
		pegged int64
		hand   int64
		crib   int64
		diff   int64
	}

	summaries := make([]summary, np)
	totalGames := np * n // each strategy plays np matchups of n games

	for k := range players {
		summaries[k].name = players[k].Name
		for j := 0; j < np; j++ {
			var p int
			var r matchupResult
			if k <= j {
				r, p = results[k][j], 0
			} else {
				r, p = results[j][k], 1
			}
			summaries[k].wins += r.wins[p]
			summaries[k].pegged += r.peggedPoints[p]
			summaries[k].hand += r.handPoints[p]
			summaries[k].crib += r.cribPoints[p]
			summaries[k].diff += r.diff[p]
		}
	}

	sort.Slice(summaries, func(i, j int) bool {
		if summaries[i].wins != summaries[j].wins {
			return summaries[i].wins > summaries[j].wins
		}
		return summaries[i].diff > summaries[j].diff
	})

	fg := float64(totalGames)
	fmt.Printf("\nRankings — %d total games per strategy\n\n", totalGames)
	fmt.Printf("%-4s %-*s %6s %8s %8s %8s %8s %8s\n",
		"Rank", rowLabel, "Strategy", "Win%", "Pegged", "Hand", "Crib", "Total", "Diff")
	for i, s := range summaries {
		total := float64(s.pegged+s.hand+s.crib) / fg
		fmt.Printf("%-4d %-*s %5.1f%% %8.2f %8.2f %8.2f %8.2f %+8.2f\n",
			i+1, rowLabel, s.name,
			float64(s.wins)/fg*100,
			float64(s.pegged)/fg,
			float64(s.hand)/fg,
			float64(s.crib)/fg,
			total,
			float64(s.diff)/fg,
		)
	}

	fmt.Printf("\nCompleted in %s\n", time.Since(start).Round(time.Millisecond))
}

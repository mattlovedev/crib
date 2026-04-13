package game

import "sort"

// ScorePeggingPlay returns the points scored by playing the last card in series.
func ScorePeggingPlay(count int, series Cards) int {
	pts := 0
	n := len(series)

	if count == 15 || count == 31 {
		pts += 2
	}

	// Pairs: consecutive matching faces from the end of series.
	if n >= 2 {
		face := series[n-1].Face
		pairLen := 1
		for i := n - 2; i >= 0 && series[i].Face == face; i-- {
			pairLen++
		}
		switch pairLen {
		case 2:
			pts += 2
		case 3:
			pts += 6
		case 4:
			pts += 12
		}
	}

	// Runs: longest run of 3+ ending at the last card played.
	for length := n; length >= 3; length-- {
		if isPeggingRun(series[n-length:]) {
			pts += length
			break
		}
	}

	return pts
}

func isPeggingRun(cards Cards) bool {
	faces := make([]int, len(cards))
	for i, c := range cards {
		faces[i] = c.Face
	}
	sort.Ints(faces)
	for i := 1; i < len(faces); i++ {
		if faces[i] != faces[i-1]+1 {
			return false
		}
	}
	return true
}

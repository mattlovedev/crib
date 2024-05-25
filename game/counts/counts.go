package counts

import (
	"encoding/json"
	"math"
	"sort"

	"mattlove.dev/crib/game"
	cribmath "mattlove.dev/crib/game/math"
	"mattlove.dev/crib/util"
)

// key is cut
type FourHandCutCounts map[string]int

func (c FourHandCutCounts) String() string {
	bytes, _ := json.MarshalIndent(c, "", "  ")
	return string(bytes)
}

type FourSummaryCounts map[int][]string

//func (c FourSummaryCounts) MarshalJSON() ([]byte, error) {}

type FourSummary struct {
	//Cuts   FourCounts
	Avg      float32
	Min      int
	Median   int
	Max      int
	Mode     int
	ModeP    float32
	BelowAvg float32
	AboveAvg float32
	Counts   FourSummaryCounts
	StdDev   float32
}

func (s FourSummary) String() string {
	//var sb strings.Builder
	//sb.WriteString(fmt.Sprintf("Avg: %f\n", s.Avg))
	//return sb.String()
	bytes, _ := json.MarshalIndent(s, "", "  ")
	return string(bytes)
}

type AllFourHandCutCounts map[string]FourHandCutCounts

func NewAllFourHandCutCounts() AllFourHandCutCounts {
	return make(AllFourHandCutCounts, cribmath.NCR52_4)
}

func ReadAllFourHandCutCounts() (AllFourHandCutCounts, error) {
	counts := NewAllFourHandCutCounts()
	err := util.ReadFile(util.FourCutsFile, &counts)
	return counts, err
}

type AllFourHandSummaries map[string]FourSummary

func NewAllFourHandSummaries(cap int) AllFourHandSummaries {
	return make(AllFourHandSummaries, cap)
}

func ReadAllFourHandSummaries() (AllFourHandSummaries, error) {
	s := NewAllFourHandSummaries(cribmath.NCR52_4)
	err := util.ReadFile(util.FourSummariesFile, &s)
	return s, err
}

func WriteAllFourHandSummaries(s AllFourHandSummaries) error {
	return util.WriteFile(util.FourSummariesFile, s)
}

func MakeFourScores(hand game.Cards) FourHandCutCounts {
	scores := make(FourHandCutCounts, game.NumCuts)
	deck := game.RemainingDeck(hand)
	for _, cut := range deck.Cards {
		count := game.CountCards(hand, &cut, false)
		scores[cut.String()] = count
	}
	return scores
}

func MakeSummaries(hand game.Cards) FourSummary {
	scores := MakeFourScores(hand)
	sum := 0
	vals := make([]int, 0, len(scores))
	countCuts := make(map[int][]string, 30)
	for cut, val := range scores {
		vals = append(vals, val)
		sum += val
		if _, found := countCuts[val]; !found {
			countCuts[val] = []string{cut}
		} else {
			countCuts[val] = append(countCuts[val], cut)
		}
	}
	avg := float32(sum) / float32(len(vals))
	sort.Ints(vals)
	min := vals[0]
	max := vals[len(vals)-1]
	median := vals[len(vals)/2]
	numBelow := 0
	numAbove := 0
	sumOfValsMinusAvgSquared := float32(0)
	for _, val := range vals {
		if float32(val) < avg {
			numBelow++
		} else if float32(val) > avg {
			numAbove++
		}
		sumOfValsMinusAvgSquared += (float32(val) - avg) * (float32(val) - avg)
	}
	below := float32(numBelow) / float32(len(vals))
	above := float32(numAbove) / float32(len(vals))

	mode := 0
	modeCount := 0
	for count, cuts := range countCuts {
		if len(cuts) > modeCount {
			mode = count
			modeCount = len(cuts)
		}
	}
	modeP := float32(modeCount) / float32(len(vals))

	stdDev := float32(math.Sqrt(float64(sumOfValsMinusAvgSquared / float32(len(vals)))))

	return FourSummary{
		Avg:      avg,
		Min:      min,
		Max:      max,
		Median:   median,
		BelowAvg: below,
		AboveAvg: above,
		Counts:   countCuts,
		Mode:     mode,
		ModeP:    modeP,
		StdDev:   stdDev,
	}
}

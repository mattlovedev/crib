package counts

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"mattlove.dev/crib/game"
	cribmath "mattlove.dev/crib/game/math"
	"mattlove.dev/crib/util"
	"sort"
)

type SixHand struct {
	Hand    game.Cards
	Crib    game.Cards
	Summary FourSummary
}

func (h SixHand) toData() sixHandChoice {
	sh := sixHandChoice{
		Hand:     uint32(h.Hand.Id()),
		Crib:     uint16(h.Crib.Id()),
		Avg:      float32(h.Summary.Avg),
		Min:      uint8(h.Summary.Min),
		Median:   uint8(h.Summary.Median),
		Max:      uint8(h.Summary.Max),
		Mode:     uint8(h.Summary.Mode),
		ModeP:    float32(h.Summary.ModeP),
		BelowAvg: uint8(h.Summary.BelowAvg),
		AboveAvg: uint8(h.Summary.AboveAvg),
		StdDev:   float32(h.Summary.StdDev),
	}
	fmt.Println(h)
	fmt.Println(sh)
	return sh
}

func (h SixHand) MarshalJSON() ([]byte, error) {
	//s := h.Summary.Interfaces()
	//i := append([]interface{}{}, h.Hand.Id(), h.Crib.Id())
	//i = append(i, s...)
	var buf bytes.Buffer
	buf.WriteString("{\"Hand\":\"")
	if _, err := buf.WriteString(h.Hand.String()); err != nil {
		return nil, err
	}
	buf.WriteString("\",\"Crib\":\"")
	if _, err := buf.WriteString(h.Crib.String()); err != nil {
		return nil, err
	}
	buf.WriteString("\",\"Summary\":")
	if _, err := buf.WriteString(h.Summary.String()); err != nil {
		return nil, err
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (h SixHand) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	err := binary.Write(&b, binary.LittleEndian, h.toData())
	return b.Bytes(), err
}

func (h *SixHand) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	s := sixHandChoice{}
	if err := binary.Read(r, binary.LittleEndian, &s); err != nil {
		panic(err)
		return err
	}
	h.Hand = game.CardsById(int(s.Hand), 4)
	h.Crib = game.CardsById(int(s.Crib), 2)
	h.Summary = FourSummary{
		Avg:      float64(s.Avg),
		Min:      int(s.Min),
		Median:   int(s.Median),
		Max:      int(s.Max),
		Mode:     int(s.Mode),
		ModeP:    float64(s.ModeP),
		BelowAvg: int(s.BelowAvg),
		AboveAvg: int(s.AboveAvg),
		StdDev:   float64(s.StdDev),
	}
	return nil
}

func (h SixHand) String() string {
	bytes, _ := json.MarshalIndent(h, "", "  ")
	return string(bytes)
}

type SixHands [15]SixHand

func (t *SixHands) Len() int           { return len(t) }
func (t *SixHands) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t *SixHands) Less(i, j int) bool { return t[i].Summary.Avg > t[j].Summary.Avg }

func (h SixHands) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	for i := range h {
		if b, err := h[i].MarshalBinary(); err != nil {
			return nil, err
		} else if _, err = buf.Write(b); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (h *SixHands) UnmarshalBinary(data []byte) error {
	size := len(data) / 15
	for i := range h {
		if err := h[i].UnmarshalBinary(data[i*size : (i+1)*size]); err != nil {
			panic(err)
			return err
		}
	}
	return nil
}

func (h SixHands) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("[")
	addComma := false
	for i := range h {
		if addComma {
			buf.WriteString(",")
		} else {
			addComma = true
		}
		b, err := h[i].MarshalJSON()
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}
	buf.WriteString("]")
	return buf.Bytes(), nil
}

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
	Avg    float64
	Min    int
	Median int
	Max    int
	Mode   int
	ModeP  float64
	//BelowAvg float32
	//AboveAvg float32
	BelowAvg int
	AboveAvg int
	Counts   FourSummaryCounts
	StdDev   float64
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func round32(num float32) int {
	return int(float64(num) + math.Copysign(0.5, float64(num)))
}

func fixed(f float64) float64 {
	output := math.Pow(10, float64(2))
	return float64(round(f*output)) / output
}

func fixed32(f float32) float32 {
	output := float32(math.Pow(10, float64(2)))
	return float32(round32(f*output)) / output
}

func (s FourSummary) Interfaces() []interface{} {
	//return []interface{}{s.Avg, s.Min, s.Median, s.Max, s.Mode, s.ModeP, s.BelowAvg, s.AboveAvg, s.StdDev}
	//return []interface{}{fmt.Sprintf("%.2f, %d, %d, %d, %.2f, %d, %d, %.2f", s.Avg, s.Min, s.Max, s.Mode, s.ModeP, s.BelowAvg, s.AboveAvg, s.StdDev)}
	i := append([]interface{}{}, fixed(s.Avg))
	i = append(i, s.Min, s.Median, s.Max, s.Mode)
	i = append(i, fixed(s.ModeP))
	i = append(i, s.BelowAvg, s.AboveAvg)
	return append(i, fixed(s.StdDev))
}

/*func (s FourSummary) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{s.Avg, s.Min, s.Median, s.Max, s.Mode, s.ModeP, s.BelowAvg, s.AboveAvg, s.StdDev})
}*/

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

func MakeFourScores(hand game.Cards, crib game.Cards) FourHandCutCounts {
	scores := make(FourHandCutCounts, game.NumCuts-len(crib))
	deck := game.RemainingDeck(hand, crib)
	for _, cut := range deck.Cards {
		count := game.CountCards(hand, &cut, false)
		scores[cut.String()] = count
	}
	return scores
}
func MakeSummaries(hand game.Cards) FourSummary {
	return makeSummaries(hand, nil, false)
}

func MakeSummariesNoCounts(hand game.Cards) FourSummary {
	return makeSummaries(hand, nil, true)
}

func MakeSummariesWithCrib(hand game.Cards, crib game.Cards) FourSummary {
	return makeSummaries(hand, crib, true)
}

func MakeSixHands(cards game.Cards) SixHands {
	sets, remaining := cards.ChooseFourWithRemaining()
	summaries := SixHands{}
	//summaries := make(map[string]counts.FourSummary, len(sets))
	for i := range sets {
		//summaries[sets[i].String()] = counts.MakeSummariesWithCrib(sets[i], remaining[i])
		summaries[i] = SixHand{
			Hand:    sets[i],
			Crib:    remaining[i],
			Summary: MakeSummariesWithCrib(sets[i], remaining[i]),
		}
	}
	sort.Sort(&summaries)
	return summaries
}

func makeSummaries(hand game.Cards, crib game.Cards, omitCounts bool) FourSummary {
	scores := MakeFourScores(hand, crib)
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
	for _, val := range countCuts {
		sort.Slice(val, func(i int, j int) bool {
			return game.CardByIdString(val[i]).Id < game.CardByIdString(val[j]).Id
		})
	}
	avg := float64(sum) / float64(len(vals))
	sort.Ints(vals)
	min := vals[0]
	max := vals[len(vals)-1]
	median := vals[len(vals)/2]
	numBelow := 0
	numAbove := 0
	sumOfValsMinusAvgSquared := float64(0)
	for _, val := range vals {
		if float64(val) < avg {
			numBelow++
		} else if float64(val) > avg {
			numAbove++
		}
		sumOfValsMinusAvgSquared += (float64(val) - avg) * (float64(val) - avg)
	}
	//below := float32(numBelow) / float32(len(vals))
	//above := float32(numAbove) / float32(len(vals))

	mode := 0
	modeCount := 0
	for count, cuts := range countCuts {
		if len(cuts) > modeCount {
			mode = count
			modeCount = len(cuts)
		}
	}
	modeP := float64(modeCount) / float64(len(vals))

	stdDev := float64(math.Sqrt(float64(sumOfValsMinusAvgSquared / float64(len(vals)))))

	fs := FourSummary{
		Avg:      fixed(avg),
		Min:      min,
		Max:      max,
		Median:   median,
		BelowAvg: numBelow,
		AboveAvg: numAbove,
		Mode:     mode,
		ModeP:    fixed(modeP),
		StdDev:   fixed(stdDev),
	}
	if !omitCounts {
		fs.Counts = countCuts
	}
	return fs
}

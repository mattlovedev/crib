package game

const (
	CribThrownCards = 2
	NumHoleCards    = 4
	NumPlayCards    = 5
	NumDealtCards   = 6

	NumCuts = NumCards - NumHoleCards

	//defaultMinCount = 30
	//defaultMaxCount = -1
)

const (
	FirstCardBit = 1 << iota
	SecondCardBit
	ThirdCardBit
	FourthCardBit
	FifthCardBit
	MaxMask
)

var (
	cardMasks = []int{FirstCardBit, SecondCardBit, ThirdCardBit, FourthCardBit, FifthCardBit, MaxMask}
)

/*type Hand struct {
	Id            string `json:"id"`
	Cards         Cards  `json:"-"`
	Cut           Card   `json:"-"`
	Count         int    `json:"count"` // sum of four hole cards
	CountWithCut  int    `json:"countCut"`
	CountWithCrib int    `json:"countCrib"`
	//WithCuts      map[string]Hand `json:"withCuts,omitempty"`
}*/

// goal is for us to be able to generate hand from either compute (while gen) or from json/file

type Hand interface {
	GetCount() int
}

type PreCutHand interface {
	Hand
	WithCut(cut Card) Hand
}

type FourCribHand struct{}

type FourPlayerHand struct{}

type FiveCribHand struct{}

type FivePlayerHand struct{}

// TODO move counts out of hand
/*type CountDetail struct {
	Value    int
	CutCards []string
}

func makeCountDetail(defVal int) CountDetail {
	return CountDetail{
		Value:    defVal,
		CutCards: make([]string, 0, NumFourHandCuts),
	}
}

type CountSummary struct {
	total   int
	Min     CountDetail
	Max     CountDetail
	Mode    CountDetail
	Median  CountDetail
	Average float32
}

func makeCountSummary() CountSummary {
	return CountSummary{
		total:   0,
		Min:     makeCountDetail(defaultMinCount),
		Max:     makeCountDetail(defaultMaxCount),
		Mode:    makeCountDetail(0),
		Median:  makeCountDetail(0),
		Average: 0,
	}
}

type FourHandSummary struct {
	Id   string
	Hand CountSummary
	//Crib CountSummary
}*/

func CountCards(hole Cards, cut *Card, asCrib bool) int {

	cards := hole.Copy()
	if cut != nil {
		cards = append(cards, *cut)
	}
	cards.Sort()

	countFifteens := func() int {
		count := 0
		for mask := 1; mask < cardMasks[len(cards)]; mask++ { // 1 bit for every card position 00001 to 11111
			sum := 0
			for card := range cards {
				if cardMasks[card]&mask > 0 {
					sum += cards[card].Value
				}
			}
			if sum == 15 {
				count += 2
			}
		}
		return count
	}

	countPairs := func() int {
		count := 0
		for i := 0; i < len(cards)-1; i++ {
			for j := i + 1; j < len(cards); j++ {
				if cards[i].Face == cards[j].Face {
					count += 2
				}
			}
		}
		return count
	}

	countRuns := func() int {
		uniques := cards.Copy()
		duplicates := make(Cards, 0, 3) // need at least 2 uniques between 5 cards

		// remove duplicates
		for i := 0; i < len(uniques)-1; i++ {
			for i < len(uniques)-1 && uniques[i].Face == uniques[i+1].Face {
				duplicates = append(duplicates, uniques[i+1])
				uniques = append(uniques[:i+1], uniques[i+2:]...)
			}
		}

		isStraight := func(start, length int) bool {
			for i := start; i < start+length-1; i++ {
				if uniques[i+1].Face-uniques[i].Face != 1 {
					return false
				}
			}
			return true
		}

		isInStraight := func(start, length int) int {
			count := 0
			dupes := make([]int, NumFaces)
			for i := 0; i < len(duplicates); i++ {
				for j := start; j < start+length; j++ {
					if duplicates[i].Face == uniques[j].Face {
						dupes[duplicates[i].Face]++
					}
				}
			}
			oneMatch := false
			for i := 0; i < len(dupes); i++ {
				if dupes[i] == 1 {
					if oneMatch {
						count += 2
					} else {
						count++
						oneMatch = true
					}
				} else if dupes[i] == 2 {
					count += 2
				}
			}
			return count
		}

		for length := len(uniques); length > 2; length-- {
			for i := 0; i <= len(uniques)-length; i++ {
				if isStraight(i, length) {
					return length * (1 + isInStraight(i, length))
				}
			}
		}

		return 0
	}

	// doesn't use cards in scope
	countFlush := func() int {
		for i := 1; i < len(hole); i++ {
			if hole[0].Suit != hole[i].Suit {
				return 0
			}
		}
		if asCrib && (cut == nil || cut.Suit != hole[0].Suit) {
			return 0
		}
		if hole[0].Suit == cut.Suit {
			return 5
		}
		return 4 // asCrib can't get down here
	}

	count := 0

	count += countFifteens()
	count += countPairs()
	count += countRuns()
	count += countFlush()
	if cut != nil && hole.Contains(CardByFaceSuit(Jack, cut.Suit)) {
		count += 1
	}

	return count
}

/*func makeHandCards(c1, c2, c3, c4 Card) Cards {
	cards := make(Cards, NumCardsInHand, MaxCardsSize)
	cards[0] = c1
	cards[1] = c2
	cards[2] = c3
	cards[3] = c4
	cards.Sort()
	return cards
}*/

/*func BuildFourHandSummary(c1, c2, c3, c4 Card) FourHandSummary {
	cards := makeHandCards(c1, c2, c3, c4)

	//type Counts struct {
	//	hand int
	//	crib int
	//}
	counts := make(map[string]int, NumFourHandCuts)

	deck := RemainingDeck(cards)

	for deck.HasCards() {
		cut := deck.DealCard()
		countWithCut := countHand(cards, cut, false)
		//countWithCrib := countHand(cards, cut, true)
		counts[cut.String()] = countWithCut
	}

	summary := FourHandSummary{
		Id:   fmt.Sprintf("%s%s%s%s", cards[0], cards[1], cards[2], cards[3]),
		Hand: makeCountSummary(),
		//Crib: makeCountSummary(),
	}

	for cut, count := range counts {
		summary.Hand.total += count

		if count < summary.Hand.Min.Value {
			summary.Hand.Min.Value = count
			summary.Hand.Min.CutCards = append(summary.Hand.Min.CutCards[:0], cut)
		} else if count == summary.Hand.Min.Value {
			summary.Hand.Min.CutCards = append(summary.Hand.Min.CutCards, cut)
		}

		if count > summary.Hand.Max.Value {
			summary.Hand.Max.Value = count
			summary.Hand.Max.CutCards = append(summary.Hand.Max.CutCards[:0], cut)
		} else if count == summary.Hand.Max.Value {
			summary.Hand.Max.CutCards = append(summary.Hand.Max.CutCards, cut)
		}
	}

	summary.Hand.Average = float32(summary.Hand.total) / NumFourHandCuts

	return summary
}*/

/*func HandFromFiveCards(c1, c2, c3, c4, cut Card) Hand {
	cards := makeHandCards(c1, c2, c3, c4)
	id := fmt.Sprintf("%s%s%s%s%s", cards[0].StringId(), cards[1].StringId(), cards[2].StringId(), cards[3].StringId(), cut.StringId())
	count := countHand(cards, nil, false)
	countWithCut := countHand(cards, &cut, false)
	countWithCrib := countHand(cards, &cut, true)
	return Hand{
		Id:            id,
		Cards:         cards,
		Cut:           cut,
		Count:         count,
		CountWithCut:  countWithCut,
		CountWithCrib: countWithCrib,
	}
}*/

/*func HandFromFourCards(c1, c2, c3, c4 Card) Hand {
	cards := makeHandCards(c1, c2, c3, c4)
	id := fmt.Sprintf("%s%s%s%s", cards[0], cards[1], cards[2], cards[3])
	withCuts := make(map[string]Hand, NumCards-4)
	return Hand{
		Id:       id,
		Cards:    cards,
		WithCuts: withCuts,
	}
}*/

/*func (h *Hand) WithCut(cut Card) {

	id := fmt.Sprintf("%s%s", h.Id, cut.StringId())
	countWithCut := countHand(h.Cards, &cut, false)
	countWithCrib := countHand(h.Cards, &cut, true)

	hand := Hand{
		Id:            id,
		Cards:         h.Cards,
		Cut:           cut,
		Count:         h.Count,
		CountWithCut:  countWithCut,
		CountWithCrib: countWithCrib,
	}

	h.WithCuts[cut.StringId()] = hand
}*/

/*func BuildAllHandsFromFourCards(c1, c2, c3, c4 Card) Hand {
	cards := makeHandCards(c1, c2, c3, c4)
	id := fmt.Sprintf("%s%s%s%s", cards[0].StringId(), cards[1].StringId(), cards[2].StringId(), cards[3].StringId())
	count := countHand(cards, nil, false)
	withCuts := make(map[string]Hand, NumCards-4)

	deck := RemainingDeck(cards)

	for deck.HasCards() {
		cut := deck.DealCard()
		subId := fmt.Sprintf("%s%s", id, cut.StringId())
		countWithCut := countHand(cards, &cut, false)
		countWithCrib := countHand(cards, &cut, true)

		hand := Hand{
			Id:            subId,
			Cut:           cut,
			Count:         count,
			CountWithCut:  countWithCut,
			CountWithCrib: countWithCrib,
		}
		withCuts[cut.StringId()] = hand
	}

	return Hand{
		Id:       id,
		Cards:    cards,
		Count:    count,
		WithCuts: withCuts,
	}
}*/

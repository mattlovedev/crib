package game

const (
	CribThrownCards = 2
	NumHoleCards    = 4
	NumPlayCards    = 5
	NumDealtCards   = 6

	NumCuts = NumCards - NumHoleCards
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
		if cut != nil && hole[0].Suit == cut.Suit {
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

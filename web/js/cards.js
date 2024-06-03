const NumFaces = 13
const NumSuits = 4
const NumCards = NumFaces * NumSuits

const Clubs = 0
const Diamonds = 1
const Hearts = 2
const Spades = 3

/*const AceOfClubsId = 0
const TwoOfClubsId = 1
const ThreeOfClubsId = 2
const FourOfClubsId = 3
const FiveOfClubsId = 4
const SixOfClubsId = 5
const SevenOfClubsId = 6
const EightOfClubsId = 7
const NineOfClubsId = 8
const TenOfClubsId = 9
const JackOfClubsId = 10
const QueenOfClubsId = 11
const KingOfClubsId = 12

const AceOfDiamondsId = 13
const TwoOfDiamondsId = 14
const ThreeOfDiamondsId = 15
const FourOfDiamondsId = 16
const FiveOfDiamondsId = 17
const SixOfDiamondsId = 18
const SevenOfDiamondsId = 19
const EightOfDiamondsId = 20
const NineOfDiamondsId = 21
const TenOfDiamondsId = 22
const JackOfDiamondsId = 23
const QueenOfDiamondsId = 24
const KingOfDiamondsId = 25

const AceOfHeartsId = 26
const TwoOfHeartsId = 27
const ThreeOfHeartsId = 28
const FourOfHeartsId = 29
const FiveOfHeartsId = 30
const SixOfHeartsId = 31
const SevenOfHeartsId = 32
const EightOfHeartsId = 33
const NineOfHeartsId = 34
const TenOfHeartsId = 35
const JackOfHeartsId = 36
const QueenOfHeartsId = 37
const KingOfHeartsId = 38

const AceOfSpadesId = 39
const TwoOfSpadesId = 40
const ThreeOfSpadesId = 41
const FourOfSpadesId = 42
const FiveOfSpadesId = 43
const SixOfSpadesId = 44
const SevenOfSpadesId = 45
const EightOfSpadesId = 46
const NineOfSpadesId = 47
const TenOfSpadesId = 48
const JackOfSpadesId = 49
const QueenOfSpadesId = 50
const KingOfSpadesId = 51
*/

/*const indexToString = [
    "ac", "2c", "3c", "4c", "5c", "6c", "7c", "8c", "9c", "tc", "jc", "qc", "kc",
    "ad", "2d", "3d", "4d", "5d", "6d", "7d", "8d", "9d", "td", "jd", "qd", "kd",
    "ah", "2h", "3h", "4h", "5h", "6h", "7h", "8h", "9h", "th", "jh", "qh", "kh",
    "as", "2s", "3s", "4s", "5s", "6s", "7s", "8s", "9s", "ts", "js", "qs", "ks"
]*/

const indexToString = [
    "ac", "ad", "ah", "as",
    "2c", "2d", "2h", "2s",
    "3c", "3d", "3h", "3s",
    "4c", "4d", "4h", "4s",
    "5c", "5d", "5h", "5s",
    "6c", "6d", "6h", "6s",
    "7c", "7d", "7h", "7s",
    "8c", "8d", "8h", "8s",
    "9c", "9d", "9h", "9s",
    "tc", "td", "th", "ts",
    "jc", "jd", "jh", "js",
    "qc", "qd", "qh", "qs",
    "kc", "kd", "kh", "ks"
]

const stringToIndex = indexToString.reduce((obj, str, i) => {
    obj[str] = i
    return obj
}, {})

const Card = id => {
    const face = Math.floor(id / NumSuits)
    const suit = id % NumSuits
    const value = (face + 1 > 10) ? 10 : face + 1
    return {
        id, face, suit, value
    }
}

const Deck = cards => {
    const removed = {}
    cards.forEach(card => {
        removed[card.id] = card
    })
    const deck = []
    for (let i = 0; i < NumCards; i++) {
        if (removed[i] == undefined) {
            deck.push(Card(i))
        }
    }
    return deck
}

const countCards = (hole, cut, asCrib) => {
    const cards = [...hole]
    cards.push(cut)
    cards.sort((a, b) => a.id - b.id)

    const countFifteens = () => {
        const FirstCardBit = 1
        const SecondCardBit = 2
        const ThirdCardBit = 4
        const FourthCardBit = 8
        const FifthCardBit = 16
        const MaxMask = 32
        const CardMasks = [ FirstCardBit, SecondCardBit, ThirdCardBit, FourthCardBit, FifthCardBit, MaxMask ]


        let count = 0
        for (let mask = 1; mask < CardMasks[cards.length]; mask++) {
            let sum = 0
            for (let card = 0; card < cards.length; card++) {
                if ((CardMasks[card] & mask) > 0) {
                    sum += cards[card].value
                }
            }
            if (sum == 15) {
                count += 2
            }
        }
        return count
    }

    const countPairs = () => {
        let count = 0
        for (let i = 0; i < cards.length-1; i++) {
            for (let j = i + 1; j < cards.length; j++) {
                if (cards[i].face == cards[j].face) {
                    count += 2
                }
            }
        }
        return count
    }

    const countRuns = () => {
        const uniques = [...cards]
        const duplicates = []

        for (let i = 0; i < uniques.length; i++) {
            while (i < uniques.length-1 && uniques[i].face == uniques[i+1].face) {
                duplicates.push(uniques[i+1])
                uniques.splice(i+1, 1)
            }
        }

        const isStraight = (start, len) => {
            for (let i = start; i < start+len-1; i++) {
                if (uniques[i+1].face - uniques[i].face != 1) {
                    return false
                }
            }
            return true
        }

        const isInStraight = (start, len) =>  {
            let count = 0
            const dupes = new Array(NumFaces).fill(0)
            for (let i = 0; i < duplicates.length; i++) {
                for (let j = start; j < start+len; j++) {
                    if (duplicates[i].face == uniques[j].face) {
                        dupes[duplicates[i].face]++
                    }
                }
            }
            let oneMatch = false
            for (let i = 0; i < dupes.length; i++) {
                if (dupes[i] == 1) {
                    if (oneMatch) {
                        count += 2
                    } else {
                        count++
                        oneMatch = true
                    }
                } else if (dupes[i] == 2) {
                    count += 2
                }
            }
            return count
        }

        for (var len = uniques.length; len > 2; len--) {
            for (var i = 0; i <= uniques.length-len; i++) {
                if (isStraight(i, len)) {
                    return len * (1 + isInStraight(i, len))
                }
            }
        }

        return 0
    }

    const countFlush = () => {
        for (let i = 1; i < hole.length; i++) {
            if (hole[0].suit != hole[i].suit) {
                return 0
            }
        }
        if (asCrib && cut.suit != hole[0].suit) {
            return 0
        }
        if (hole[0].suit == cut.suit) {
            return 5
        }
        return 4
    }

    let count = 0
    count += countFifteens()
    count += countPairs()
    count += countRuns()
    count += countFlush()
    hole.forEach(holeCard => {
        if (holeCard.face == 10 && holeCard.suit == cut.suit) { // 10 is Jack
            count += 1
        }
    })
    return count
}

const makeFourScores = hand => {
    const scores = {}
    const deck = Deck(hand)
    deck.forEach(cut => {
        const count = countCards(hand, cut, false)
        scores[cut.id] = count
    })
    return scores
}


const makeSummaries = hand => {
    const scores = makeFourScores(hand)
	let sum = 0
	const vals = []
	const countCuts = {}
    for (const [cut, val] of Object.entries(scores)) {
        vals.push(val)
        sum += val
        if (countCuts[val] == undefined) {
            countCuts[val] = []
        }
        countCuts[val].push(cut)
    }
    const avg = sum / vals.length
    vals.sort()
    const min = vals[0]
    const max = vals[vals.length-1]
    const median = vals[vals.length/2]
    let numBelow = 0
    let numAbove = 0
    let sumOfValsMinusAvgSquared = 0
    for (const val of vals) {
        if (val < avg) {
            numBelow++
        } else if (val > avg) {
            numAbove++
        }
        sumOfValsMinusAvgSquared += (val - avg) * (val - avg)
    }
    let mode = 0
    let modeCount = 0
    for (const [count, cuts] of Object.entries(countCuts)) {
        if (cuts.length > modeCount) {
            mode = count
            modeCount = cuts.length
        }
    }
    const modeP = modeCount / vals.length
    const stdDev = Math.sqrt(sumOfValsMinusAvgSquared / vals.length)
    return {
        Avg: avg.toFixed(2),
        Min: min,
        Max: max,
        Median: median,
        BelowAvg: numBelow,
        AboveAvg: numAbove,
        Mode: mode,
        ModeP: modeP.toFixed(2),
        StdDev: stdDev.toFixed(2),
        Counts: countCuts
    }
}
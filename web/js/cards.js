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
package engine

import (
	"math/rand"

	"mattlove.dev/crib/game"
	"mattlove.dev/crib/strategy"
)

const winScore = 121

// GameResult holds the outcome of a single game.
type GameResult struct {
	Winner       int
	Scores       [2]int
	PeggedPoints [2]int
	HandPoints   [2]int // includes nibs (cut Jack = 2 pts for dealer)
	CribPoints   [2]int
}

// RunGame plays a complete game between two players and returns the result.
// dealer is the index (0 or 1) of the first dealer. rng is used for shuffling.
func RunGame(p0, p1 strategy.Player, dealer int, rng *rand.Rand) GameResult {
	players := [2]strategy.Player{p0, p1}
	result := GameResult{}

	for result.Scores[0] < winScore && result.Scores[1] < winScore {
		if playRound(players, dealer, rng, &result) {
			break
		}
		dealer = 1 - dealer
	}

	if result.Scores[0] >= winScore {
		result.Winner = 0
	} else {
		result.Winner = 1
	}
	return result
}

// playRound runs one complete round: deal, discard, cut, peg, score.
// Returns true if a player reached winScore during this round.
func playRound(players [2]strategy.Player, dealer int, rng *rand.Rand, result *GameResult) bool {
	nonDealer := 1 - dealer

	// Shuffle and deal 6 cards to each player; card [12] is the cut.
	deck := game.NewDeck()
	rng.Shuffle(len(deck.Cards), func(i, j int) {
		deck.Cards[i], deck.Cards[j] = deck.Cards[j], deck.Cards[i]
	})
	hands := [2]game.Cards{
		deck.Cards[0:6].Copy(),
		deck.Cards[6:12].Copy(),
	}
	cut := deck.Cards[12]

	// Discard phase: each player throws 2 cards to the crib.
	var kept [2]game.Cards
	crib := make(game.Cards, 0, 4)
	for p := 0; p < 2; p++ {
		var discard game.Cards
		kept[p], discard = players[p].Discard(hands[p], p == dealer)
		crib = append(crib, discard...)
	}

	// Nibs: cut card is a Jack — dealer scores 2 immediately.
	if cut.Face == game.Jack {
		result.Scores[dealer] += 2
		result.HandPoints[dealer] += 2
		if result.Scores[dealer] >= winScore {
			return true
		}
	}

	// Pegging phase.
	if doPegging(players, kept, dealer, result) {
		return true
	}

	// Scoring: non-dealer hand, then dealer hand, then crib.
	if addPoints(nonDealer, game.CountCards(kept[nonDealer], &cut, false), result) {
		return true
	}
	if addPoints(dealer, game.CountCards(kept[dealer], &cut, false), result) {
		return true
	}
	dealerCribPts := game.CountCards(crib, &cut, true)
	result.Scores[dealer] += dealerCribPts
	result.CribPoints[dealer] += dealerCribPts
	return result.Scores[dealer] >= winScore
}

func addPoints(player, pts int, result *GameResult) bool {
	result.Scores[player] += pts
	result.HandPoints[player] += pts
	return result.Scores[player] >= winScore
}

package engine

import (
	"mattlove.dev/crib/game"
	"mattlove.dev/crib/strategy"
)

// doPegging runs the pegging phase, updating result scores.
// Returns true if a player reached winScore during pegging.
func doPegging(players [2]strategy.Player, hands [2]game.Cards, dealer int, result *GameResult) bool {
	inHand := [2]game.Cards{hands[0].Copy(), hands[1].Copy()}
	var series game.Cards
	count := 0
	current := 1 - dealer // non-dealer leads
	lastToPlay := -1

	for len(inHand[0]) > 0 || len(inHand[1]) > 0 {
		canPlay := [2]bool{
			canPlayAny(inHand[0], count),
			canPlayAny(inHand[1], count),
		}

		if !canPlay[0] && !canPlay[1] {
			// Both stuck: award 1 go point to last player (count==31 already scored 2).
			if lastToPlay >= 0 && count != 31 {
				result.Scores[lastToPlay]++
				result.PeggedPoints[lastToPlay]++
				if result.Scores[lastToPlay] >= winScore {
					return true
				}
			}
			// Reset for next series; player who didn't play last leads.
			series = nil
			count = 0
			if lastToPlay >= 0 {
				current = 1 - lastToPlay
			}
			lastToPlay = -1
			continue
		}

		if !canPlay[current] {
			current = 1 - current
			continue
		}

		card := players[current].Play(inHand[current], strategy.PeggingState{Count: count, Series: series})
		inHand[current] = removeCard(inHand[current], card)
		series = append(series, card)
		count += card.Value
		lastToPlay = current

		pts := game.ScorePeggingPlay(count, series)
		result.Scores[current] += pts
		result.PeggedPoints[current] += pts
		if result.Scores[current] >= winScore {
			return true
		}

		if count == 31 {
			series = nil
			count = 0
			lastToPlay = -1
			current = 1 - current
		} else {
			current = 1 - current
		}
	}

	return false
}

func canPlayAny(hand game.Cards, count int) bool {
	for _, c := range hand {
		if count+c.Value <= 31 {
			return true
		}
	}
	return false
}

func removeCard(hand game.Cards, card game.Card) game.Cards {
	out := make(game.Cards, 0, len(hand)-1)
	removed := false
	for _, c := range hand {
		if !removed && c.Id == card.Id {
			removed = true
			continue
		}
		out = append(out, c)
	}
	return out
}


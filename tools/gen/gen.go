package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"log"
	"sync"
	"time"

	"mattlove.dev/crib/game"
	"mattlove.dev/crib/game/counts"
	"mattlove.dev/crib/game/math"
	"mattlove.dev/crib/util"
)

func generateAllHands() error {
	hands := game.NewDeck().Cards.ChoseFour()

	allScores := counts.NewAllFourHandCutCounts()

	for _, hand := range hands {
		allScores[hand.String()] = counts.MakeFourScores(hand, nil)
	}

	return util.WriteFile(util.FourCutsFile, allScores)
}

func generateSplitAllHands() error {
	hands := game.NewDeck().Cards.ChoseFour()
	scoresMaps := make(map[string]counts.AllFourHandCutCounts)

	for _, hand := range hands {
		handString := hand.String()
		scoreMap, found := scoresMaps[handString[0:1]]
		if !found {
			scoresMaps[handString[0:1]] = counts.NewAllFourHandCutCounts()
			scoreMap = scoresMaps[handString[0:1]]
		}
		scoreMap[handString] = counts.MakeFourScores(hand, nil)
	}

	for prefix, scoreMap := range scoresMaps {
		if err := util.WriteFile(fmt.Sprintf("scores/four_cuts_%s.json.gz", prefix), scoreMap); err != nil {
			return err
		}
	}
	return nil
}

func generateAllSummaries() error {
	hands := game.NewDeck().Cards.ChoseFour()

	allSummaries := counts.NewAllFourHandSummaries(math.NCR52_4)

	for _, hand := range hands {
		allSummaries[hand.String()] = counts.MakeSummaries(hand)
	}

	return counts.WriteAllFourHandSummaries(allSummaries)
}

const FourPrime = 17
const SixPrime = 47

func generateSplitAllSummaries() error {
	hands := game.NewDeck().Cards.ChoseFour()
	summariesMaps := make(map[int]counts.AllFourHandSummaries, FourPrime)

	for i := 0; i < FourPrime; i++ {
		summariesMaps[i] = counts.NewAllFourHandSummaries(math.NCR52_4 / FourPrime)
	}

	for _, hand := range hands {
		handString := hand.String()
		hash := sha1.Sum([]byte(handString))
		buf := bytes.NewBuffer(hash[:])
		var prefix int64
		if err := binary.Read(buf, binary.LittleEndian, &prefix); err != nil {
			log.Fatalf(err.Error())
		}
		prefixI := (int(prefix%FourPrime) + FourPrime) % FourPrime
		summariesMaps[prefixI][handString] = counts.MakeSummaries(hand)
	}

	for prefix, summaryMap := range summariesMaps {
		//if err := util.WriteFile(fmt.Sprintf("%s/four_summaries_%s.json.gz", dumpDir, prefix), summaryMap); err != nil {
		if err := util.WriteFile(fmt.Sprintf("scores/four/four_summaries_%d.json", prefix), summaryMap); err != nil {
			return err
		}
	}
	return nil
}

func generateSixes() error {
	sixes := []game.Cards{
		game.Cards{
			game.CardById(game.FourOfClubs),
			game.CardById(game.FourOfDiamonds),
			game.CardById(game.FiveOfClubs),
			game.CardById(game.FiveOfDiamonds),
			game.CardById(game.SixOfClubs),
			game.CardById(game.SixOfDiamonds),
		},
	}
	maps := make(map[string]counts.SixHands)
	for _, hand := range sixes {
		maps[hand.String()] = counts.MakeSixHands(hand)
	}

	for h, s := range maps {
		if err := util.WriteFile(fmt.Sprintf("scores/six/%s.json", h), s); err != nil {
			log.Fatal(err)
		}
	}

	return util.WriteFile("scores/six/six_summaries_0.json", maps)
}

func generateSixSetup() error {
	begin := time.Now()
	gen := math.NewCombinationGenerator(52, 6)
	for gen.Next() {
		/*cards*/ _ = game.CardsByIds(gen.Combination(nil))
		//fmt.Println(cards)
	}
	fmt.Println(time.Now().Sub(begin))

	return nil
}

func hashedCards(cards game.Cards, prime int) int {
	hash := sha1.Sum([]byte(cards.String()))
	buf := bytes.NewBuffer(hash[:])
	var prefix int64
	if err := binary.Read(buf, binary.LittleEndian, &prefix); err != nil {
		log.Fatalf(err.Error())
	}
	prefixI := (int(prefix%int64(prime)) + prime) % prime
	return prefixI
}

func generateFoursConcurrent() {
	chs := make([]chan game.Cards, FourPrime)
	var wg sync.WaitGroup
	gen := math.NewCombinationGenerator(52, 4)
	for i := range chs {
		wg.Add(1)
		chs[i] = make(chan game.Cards, 100)
		go func() {
			defer wg.Done()
			generateFour(i, chs[i])
		}()
	}
	for gen.Next() {
		cards := game.CardsByIds(gen.Combination(nil))
		prefixI := hashedCards(cards, FourPrime)
		chs[prefixI] <- cards
	}
	for i := range chs {
		close(chs[i])
	}
	wg.Wait()
}

func generateSixsConcurrent() {
	chs := make([]chan game.Cards, SixPrime)
	var wg sync.WaitGroup
	gen := math.NewCombinationGenerator(52, 6)
	for i := range chs {
		wg.Add(1)
		chs[i] = make(chan game.Cards, 100)
		chNum := i
		ch := chs[i]
		go func() {
			defer wg.Done()
			generateSixesJob(chNum, ch)
		}()
	}
	for gen.Next() {
		cards := game.CardsByIds(gen.Combination(nil))
		prefixI := hashedCards(cards, SixPrime)
		chs[prefixI] <- cards
	}
	for i := range chs {
		close(chs[i])
	}
	wg.Wait()
}

func generateFoursIndex() {

}

func generateOne() {
	/*cards := game.CardsByIds([]int{4, 5, 6, 7})
	four := counts.MakeSummaries(cards)
	util.WriteFileNoErr("scores/four.json", four)*/

	util.WriteFileNoErr("scores/six.json",
		counts.MakeSixHands(
			game.CardsByIds([]int{0, 1, 2, 3, 4, 5})))

	util.BinaryWriteFileNoErr("scores/six.dat",
		counts.MakeSixHands(
			game.CardsByIds([]int{0, 1, 2, 3, 4, 5})))
}

func main() {
	//fmt.Println(generateAllHands())
	//fmt.Println(generateAllSummaries())
	//fmt.Println(generateSplitAllHands())
	//fmt.Println(generateSplitAllSummaries())
	//fmt.Println(generateSixes())
	//fmt.Println(generateSixSetup())
	//generateFoursConcurrent()
	//generateSixsConcurrent()
	//generateFoursIndex()
	generateOne()
}

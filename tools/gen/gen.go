package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"mattlove.dev/crib/game"
	"mattlove.dev/crib/game/counts"
	"mattlove.dev/crib/game/math"
	"mattlove.dev/crib/util"
)

func generateAllHands() error {
	hands := game.NewDeck().Cards.ChoseFour()

	allScores := counts.NewAllFourHandCutCounts()

	for _, hand := range hands {
		allScores[hand.String()] = counts.MakeFourScores(hand)
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
		scoreMap[handString] = counts.MakeFourScores(hand)
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

func generateSplitAllSummaries(dumpDir string) error {
	hands := game.NewDeck().Cards.ChoseFour()
	summariesMaps := make(map[int]counts.AllFourHandSummaries, 49)

	for i := 0; i < 49; i++ {
		summariesMaps[i] = counts.NewAllFourHandSummaries(math.NCR52_4 / 49)
	}

	for _, hand := range hands {
		handString := hand.String()
		hash := sha1.Sum([]byte(handString))
		buf := bytes.NewBuffer(hash[:])
		var prefix int64
		if err := binary.Read(buf, binary.LittleEndian, &prefix); err != nil {
			log.Fatalf(err.Error())
		}
		prefixI := (int(prefix%49) + 49) % 49
		summariesMaps[prefixI][handString] = counts.MakeSummaries(hand)
	}

	for prefix, summaryMap := range summariesMaps {
		//if err := util.WriteFile(fmt.Sprintf("%s/four_summaries_%s.json.gz", dumpDir, prefix), summaryMap); err != nil {
		if err := util.WriteFile(fmt.Sprintf("%s/four_summaries_%d.json", dumpDir, prefix), summaryMap); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	//fmt.Println(generateAllHands())
	//fmt.Println(generateAllSummaries())
	//fmt.Println(generateSplitAllHands())
	fmt.Println(generateSplitAllSummaries(os.Args[1]))
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"mattlove.dev/crib/game"
	"mattlove.dev/crib/game/counts"
)

func generateFour(chNum int, ch chan game.Cards) {
	file, err := os.Create(fmt.Sprintf("scores/four/four_summaries_%d.json", chNum))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = file.Close()
	}()
	if _, err = file.WriteString("{"); err != nil {
		log.Fatal(err)
	}
	insertComma := false
	for cards := range ch {
		if insertComma {
			if _, err = file.WriteString(","); err != nil {
				log.Fatal(err)
			}
		} else {
			insertComma = true
		}
		if _, err = file.WriteString(fmt.Sprintf("\"%s\":", cards.String())); err != nil {
			log.Fatal(err)
		}
		summary := counts.MakeSummaries(cards)
		bytes, err := json.Marshal(summary)
		if err != nil {
			log.Fatal(err)
		}
		if _, err = file.Write(bytes); err != nil {
			log.Fatal(err)
		}
	}
	if _, err = file.WriteString("}"); err != nil {
		log.Fatal(err)
	}
}

func generateSixesJob(chNum int, ch chan game.Cards) {
	file, err := os.Create(fmt.Sprintf("scores/six/six_summaries_%d.json", chNum))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = file.Close()
	}()
	if _, err = file.WriteString("{"); err != nil {
		log.Fatal(err)
	}
	insertComma := false
	for cards := range ch {
		if insertComma {
			if _, err = file.WriteString(","); err != nil {
				log.Fatal(err)
			}
		} else {
			insertComma = true
		}
		if _, err = file.WriteString(fmt.Sprintf("\"%s\":", cards.String())); err != nil {
			log.Fatal(err)
		}
		summary := counts.MakeSixHands(cards)
		bytes, err := json.Marshal(summary)
		if err != nil {
			log.Fatal(err)
		}
		if _, err = file.Write(bytes); err != nil {
			log.Fatal(err)
		}
	}
	if _, err = file.WriteString("}"); err != nil {
		log.Fatal(err)
	}
}

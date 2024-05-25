package main

import (
	"fmt"
	"os"

	"mattlove.dev/crib/game/counts"
)

func dumpFourCuts() {
	if c, err := counts.ReadAllFourHandCutCounts(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(c)
	}
}

func dumpFourSummaries() {
	if s, err := counts.ReadAllFourHandSummaries(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(s)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: dump <count>")
		return
	}

	switch os.Args[1] {
	case "4c":
		dumpFourCuts()
	case "4s":
		dumpFourSummaries()
	default:
		fmt.Printf("Invalid arg: %s\n", os.Args[1])
	}
}

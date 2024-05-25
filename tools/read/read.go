package main

import (
	"log"
	"mattlove.dev/crib/util"
	"os"
)

func main() {
	var obj map[string]interface{}
	if err := util.ReadFile(os.Args[1], &obj); err != nil {
		log.Fatalf(err.Error())
	}
	log.Println(obj[os.Args[2]])
}

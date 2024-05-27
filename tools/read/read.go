package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mattlove.dev/crib/util"
	"os"
)

func _main() {
	var obj map[string]interface{}
	if err := util.ReadFile(os.Args[1], &obj); err != nil {
		log.Fatal(err)
	}
	log.Println(obj[os.Args[2]])
}

func main() {
	infos, err := os.ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	for _, info := range infos {
		var obj map[string]interface{}
		if bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", os.Args[1], info.Name())); err != nil {
			log.Fatal(err)
		} else if err = json.Unmarshal(bytes, &obj); err != nil {
			log.Fatal(err)
		}
		fmt.Println(info.Name(), len(obj))
	}
}

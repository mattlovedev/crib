package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mattlove.dev/crib/game/counts"
	"mattlove.dev/crib/util"
	"os"
)

func main1() {
	var obj map[string]interface{}
	if err := util.ReadFile(os.Args[1], &obj); err != nil {
		log.Fatal(err)
	}
	log.Println(obj[os.Args[2]])
}

func main2() {
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

func main3() {
	h := counts.SixHands{}
	util.BinaryReadFileNoErr("scores/six.dat", &h)
	b, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func main4() {
	h := counts.SixHands{}
	if err := util.ReadFile("scores/six.json", &h); err != nil {
		panic(err)
	}
	b, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func main() {
	main3()
	//main4()
}

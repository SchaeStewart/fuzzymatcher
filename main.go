package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tidwall/gjson"
	// . "github.com/kkdai/trigram"
)

func readJSON(path string) gjson.Result {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	return gjson.ParseBytes(byteValue)
}

func main() {
	readJSON("durham_pool_permits.json")
}

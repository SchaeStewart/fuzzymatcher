package main

import (
	"flag"
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
	var file1 string
	var file1Selector string
	var file2 string
	var file2Selector string
	flag.StringVar(&file1, "file1", "", "the first file")
	flag.StringVar(&file1Selector, "file1-selector", "", "the first file")
	flag.StringVar(&file2, "file2", "", "the second file")
	flag.StringVar(&file2Selector, "file2-selector", "", "TODO: better help messages")
	flag.Parse()

	file1Res := readJSON(file1)
	file2Res := readJSON(file2)

	fmt.Println(file1Res.Get(file1Selector))
	fmt.Println(file2Res.Get(file2Selector))
}

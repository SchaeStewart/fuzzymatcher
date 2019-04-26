package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/kkdai/trigram"
	"github.com/tidwall/gjson"
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

	f1Result := readJSON(file1).Get(file1Selector).Array()
	f2Result := readJSON(file2).Get(file2Selector).Array()

	// TODO: get many from selectors. Maybe only except JSON arrays as valid input
	// O(n^2) compare each item in list1 to each item in list2
	// `#` is like * but for lists

	for _, f1Item := range f1Result {
		ti := NewTrigramIndex()
		ti.Add(f1Item.String())
		for _, f2Item := range f2Result {
			ret := ti.Query(f2Item.String())
			if len(ret) > 0 {
				fmt.Println(f1Item)
				fmt.Println(f2Item)
			}
		}
	}

}

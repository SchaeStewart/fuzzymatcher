package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kkdai/trigram"
	"github.com/tidwall/gjson"
)

/* TODO:
Ensure users know that both files must be arrays
Flag help messages
Better search algorith:
- Create a list of all TrigramIndexes for r1
- For each item in r2, compare against each item in list of TrigramIndexes
Multiple selector support
*/

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
	flag.StringVar(&file2Selector, "file2-selector", "", "Test")
	flag.Parse()

	f1Result := readJSON(file1)
	f2Result := readJSON(file2)

	for _, r1 := range f1Result.Array() {
		ti := trigram.NewTrigramIndex()
		ti.Add(r1.Get(file1Selector).String())
		// TODO: improve search time
		// Store results in array and print array after complete
		// Test against pull permits
		// Pick up here

		for _, r2 := range f2Result.Array() {
			ret := ti.Query(r2.Get(file2Selector).String())
			if len(ret) > 0 {
				fmt.Println(r1)
			}
		}
	}

}

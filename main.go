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
- Store results in array and print array after complete
- Test against pull permits
*/

type results [][2]gjson.Result

func (r *results) print() {
	fmt.Println("[")
	for i, result := range *r {
		fmt.Println("[")
		fmt.Println(result[0], ",")
		fmt.Print(result[1])
		if i < len(*r)-1 {
			fmt.Println("],")
		} else {
			fmt.Println("]")
		}
	}
	fmt.Println("]")
}

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

	f1Result := readJSON(file1).Array()
	f2Result := readJSON(file2).Array()

	var res results
	for _, r1 := range f1Result {
		ti := trigram.NewTrigramIndex()
		ti.Add(r1.Get(file1Selector).String())
		for _, r2 := range f2Result {
			ret := ti.Query(r2.Get(file2Selector).String())
			if len(ret) > 0 {
				res = append(res, [2]gjson.Result{r1.Get("@pretty"), r2.Get("@pretty")})
			}
		}
	}

	res.print()

}

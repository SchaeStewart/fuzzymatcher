package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/kkdai/trigram"

	"github.com/tidwall/gjson"
)

/* TODO:
Ensure users know that both files must be arrays
Flag help messages
Multiple selector support
- Test against pull permits
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

type resultMap map[int][]gjson.Result

func createSearchTerms(f1Result []gjson.Result, f2Result []gjson.Result) (resultMap, resultMap, [][]*trigram.TrigramIndex) {
	f1SearchTerms := make(resultMap, len(f1Result))
	tIdxs := make([][]*trigram.TrigramIndex, len(f1Result))
	f2SearchTerms := make(resultMap, len(f2Result))

	for f1Idx, f1R := range f1Result {
		terms := gjson.GetMany(f1R.String(), file1Selectors...)
		f1SearchTerms[f1Idx] = terms

		tis := make([]*trigram.TrigramIndex, len(terms))
		for i, term := range terms {
			ti := trigram.NewTrigramIndex()
			ti.Add(term.String())
			tis[i] = ti
		}
		tIdxs[f1Idx] = tis
	}

	for f2Idx, f2R := range f2Result {
		f2SearchTerms[f2Idx] = gjson.GetMany(f2R.String(), file2Selectors...)
	}

	return f1SearchTerms, f2SearchTerms, tIdxs
}

func (r *results) Match(f1Result []gjson.Result, f2Result []gjson.Result, file1Selectors []string, file2Selectors []string) {
	if len(file1Selectors) != len(file2Selectors) {
		// TODO: return error instead of failing
		log.Fatal("You must have the same number of selectors in file1Selectors and file2Selectors")
	}

	_, f2SearchTerms, tIdxs := createSearchTerms(f1Result, f2Result)

	var matchIds []struct{ f1Idx, f2Idx int }
	for f2Idx, f2Term := range f2SearchTerms {
		for f1Idx, tis := range tIdxs {
			for termIdx, term := range f2Term {
				ret := tis[termIdx].Query(term.String())
				if len(ret) == 0 {
					break
				}
				matchIds = append(matchIds, struct{ f1Idx, f2Idx int }{f1Idx, f2Idx})
			}
		}
	}

	for _, ids := range matchIds {
		*r = append(*r, [2]gjson.Result{f1Result[ids.f1Idx].Get("@pretty"), f2Result[ids.f2Idx].Get("@pretty")})
	}
}

var file1 string
var file1Selectors []string
var file2 string
var file2Selectors []string

func init() {
	var file1Selector string
	var file2Selector string
	flag.StringVar(&file1, "file1", "", "the first file")
	flag.StringVar(&file1Selector, "file1-selector", "", "the first file")
	flag.StringVar(&file2, "file2", "", "the second file")
	flag.StringVar(&file2Selector, "file2-selector", "", "Test")
	flag.Parse()
	file1Selectors = strings.Split(file1Selector, ",")
	file2Selectors = strings.Split(file2Selector, ",")
}

func main() {
	f1Result := readJSON(file1).Array()
	f2Result := readJSON(file2).Array()

	var r results
	r.Match(f1Result, f2Result, file1Selectors, file2Selectors)
	r.print()
}

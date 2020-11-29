package main

import (
	dervaze "dervaze/lang"
	"log"
	"regexp"
	"strings"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func generateListOfStrings() *[]string {
	searchMap := dervaze.GetTurkishLatinIndex()
	searchList := make([]string, 0)

	for _, d := range *searchMap {
		searchList = append(searchList, d...)
	}
	return &searchList
}

func generateSingleString() *string {
	searchMap := dervaze.GetTurkishLatinIndex()

	var sb strings.Builder
	for _, d := range *searchMap {
		sb.WriteString(strings.Join(d, "\n"))
	}

	out := sb.String()
	return &out
}

func searchListOfStrings(data *[]string, str string) *[]string {
	defer timeTrack(time.Now(), "List of Strings")
	runes := []rune(str)

	sb := strings.Builder{}
	sb.WriteString(".*")
	for _, r := range runes {
		sb.WriteRune(r)
		sb.WriteString(".*")
	}
	searchRegex := regexp.MustCompile(sb.String())
	out := make([]string, 0)

	for _, l := range *data {
		if searchRegex.MatchString(l) {
			out = append(out, l)
		}
	}
	log.Println("List of Strings Search Results")
	log.Println(out[:20])
	return &out
}

func searchSingleString(data *string, str string) *[]string {
	defer timeTrack(time.Now(), "Single String")

	runes := []rune(str)
	var sb strings.Builder

	sb.WriteString("(.*")

	for _, r := range runes {
		sb.WriteRune(r)
		sb.WriteString(".*")
	}

	sb.WriteString("\n)")
	searchRegex := regexp.MustCompile(sb.String())

	results := searchRegex.FindAllString(*data, -1)
	log.Println("Single Search Results")
	log.Println(results[:20])

	return &results
}

func main() {
	dervaze.InitSearch("../assets/dervaze-rootset.protobuf")
	data1 := generateListOfStrings()
	data2 := generateSingleString()

	searchListOfStrings(data1, "emre")
	searchSingleString(data2, "emre")

}

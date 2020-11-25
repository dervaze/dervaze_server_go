package main

import (
	dervaze "dervaze/lang"
	"encoding/csv"
	"fmt"
	"io"
	"ioutil"
	"log"
	"path/filepath"
)

func readCSVFile(filename string, pos dervaze.PartOfSpeech, roots *[]Root, rooti *int) []Root {
	data, err = ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}

	csvr := csv.NewReader(string(data))
	records, err := csvr.ReadAll()
	if err != nil {
		log.Println(err)
	}

	for i, record := range records {
		if len(record) == 2 {
			latin := record[0]
			visenc := record[1]
			*roots[*rooti] = dervaze.MakeRoot(latin, visenc, pos)
			*rooti++
		} else {
			log.Println("Record error in %s line %d - %s", vfn, i, record)

		}
	}
}

func loadWordFiles() []Root {

	roots := make([]Root, 0, 250000)
	rooti := 0

	verbFiles, err := filepath.Glob("../../assets/rootdata/v/*.csv")
	if err != nil {
		verbFiles = make([]string, 0)
	}

	nounFiles, err := filepath.Glob("../../assets/rootdata/n/*.csv")
	if err != nil {
		nounFiles = make([]string, 0)
	}
	properFiles, err := filepath.Glob("../../assets/rootdata/p/*.csv")
	if err != nil {
		properFiles = make([]string, 0)
	}

	fmt.Println(verbFiles)
	fmt.Println(nounFiles)
	fmt.Println(properFiles)

	for _, fn := range verbFiles {
		readCSVFile(fn, dervaze.PartOfSpeech_VERB, &roots, &rooti)
	}
	for _, fn := range nounFiles {
		readCSVFile(fn, dervaze.PartOfSpeech_NOUN, &roots, &rooti)
	}

	for _, fn := range properFiles {
		readCSVFile(fn, dervaze.PartOfSpeech_PROPER_NOUN, &roots, &rooti)
	}

	Log.Println("Read %d records", rooti)

}

func storeProtobuf(roots []Root) {
	dervaze.
}

func main() {

	// dervaze.MakeRoot("emre", "emrh")
	// ow := dervaze.OttomanWord{
	//
	// 	Visenc:           "aa",
	// 	Unicode:          "bb",
	// 	Abjad:            123,
	// 	VisencLetters:    []string{"a", "a"},
	// 	SearchKey:        "aa",
	// 	DotlessSearchKey: "dd",
	// }
	// fmt.Println(ow.Abjad)

	roots := loadWordFiles()

}

package main

import (
	dervaze "dervaze/lang"
	"encoding/csv"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	// "time"
	"unicode/utf8"
)

func readCSVFile(filename string, pos dervaze.PartOfSpeech) []*dervaze.Root {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}

	csvr := csv.NewReader(strings.NewReader(string(data)))
	records, err := csvr.ReadAll()
	if err != nil {
		log.Println(err)
	}

	roots := make([]*dervaze.Root, 0)

	for i, record := range records {
		if len(record) == 2 {
			latin := record[0]
			visenc := record[1]
			roots = append(roots, &dervaze.MakeRoot(latin, visenc, pos))
		} else {
			log.Println("Record error in %s line %d - %s", filename, i, record)
		}
	}
	return roots[:]
}

func loadWordFiles() []dervaze.Root {

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

	rootset = new(dervaze.RootSet)

	for _, fn := range verbFiles {
		fmt.Println("%s\n", fn)
		froots := readCSVFile(fn, dervaze.PartOfSpeech_VERB)
		rootset.Roots = append(rootset.Roots, froots...)
	}
	for _, fn := range nounFiles {
		fmt.Println("%s\n", fn)
		froots := readCSVFile(fn, dervaze.PartOfSpeech_NOUN)
		rootset.Roots = append(rootset.Roots, froots...)
	}

	for _, fn := range properFiles {
		fmt.Println("%s\n", fn)
		froots := readCSVFile(fn, dervaze.PartOfSpeech_PROPER_NOUN)
		rootset.Roots = append(rootset.Roots, froots...)
	}

	log.Println("Read %d records", len(rootset.Roots))

	return rootset

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

	rootset := loadWordFiles()
	SaveRootSetProtobuf(rootset)

}

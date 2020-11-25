package main

import (
	dervaze "dervaze/lang"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	// "time"
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
			root := dervaze.MakeRoot(latin, visenc, pos)
			roots = append(roots, &root)
		} else {
			log.Println("Record error in %s line %d - %s", filename, i, record)
		}
	}
	return roots[:]
}

func loadWordFiles(rootdatadir string) *dervaze.RootSet {

	verbglob := fmt.Sprintf("%s/v/*.csv", rootdatadir)
	nounglob := fmt.Sprintf("%s/n/*.csv", rootdatadir)
	propernounglob := fmt.Sprintf("%s/p/*.csv", rootdatadir)

	verbFiles, err := filepath.Glob(verbglob)
	if err != nil {
		verbFiles = make([]string, 0)
	}

	nounFiles, err := filepath.Glob(nounglob)
	if err != nil {
		nounFiles = make([]string, 0)
	}
	properFiles, err := filepath.Glob(propernounglob)
	if err != nil {
		properFiles = make([]string, 0)
	}

	rootset := new(dervaze.RootSet)

	for _, fn := range verbFiles {
		fmt.Printf("%s\n", fn)
		froots := readCSVFile(fn, dervaze.PartOfSpeech_VERB)
		rootset.Roots = append(rootset.Roots, froots...)
	}
	for _, fn := range nounFiles {
		fmt.Printf("%s\n", fn)
		froots := readCSVFile(fn, dervaze.PartOfSpeech_NOUN)
		rootset.Roots = append(rootset.Roots, froots...)
	}

	for _, fn := range properFiles {
		fmt.Printf("%s\n", fn)
		froots := readCSVFile(fn, dervaze.PartOfSpeech_PROPER_NOUN)
		rootset.Roots = append(rootset.Roots, froots...)
	}

	log.Printf("Read %d records", len(rootset.Roots))

	return rootset

}

func main() {

	var inputdir string
	var outputfile string
	flag.StringVar(&inputdir, "i", "assets/rootdata/", "Input dir where n/ v/ p/ directories reside")
	flag.StringVar(&outputfile, "o", "assets/dervaze-rootset.protobuf", "Output file to store the protobuf file")

	flag.Parse()

	rootset := loadWordFiles(inputdir)
	dervaze.SaveRootSetProtobuf(outputfile, rootset)

}

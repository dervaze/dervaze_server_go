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

func csvLine(root *dervaze.Root) string {
	line := root.TurkishLatin + "," + root.Ottoman.Visenc + "," + root.Ottoman.Unicode + "\n"
	return line
}

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

func loadRootSet(inputfile string, outputdir string) *dervaze.RootSet {
	dervaze.InitSearch(inputfile)
	rootSet := dervaze.GetRootSet()

	verbFiles := make(map[string]string)
	nounFiles := make(map[string]string)
	properFiles := make(map[string]string)
	otherFiles := make(map[string]string)

	for i, r := range rootSet {
		line := csvLine(r)
		key := string([]rune{line}[0])
		switch r.PartOfSpeech {
		case dervaze.PartOfSpeech_NOUN:
			nounFiles[key] += line
		case dervaze.PartOfSpeech_VERB:
			verbFiles[key] += line
		case dervaze.PartOfSpeech_PROPER_NOUN:
			properFiles[key] += line
		default:
			otherFiles[key] += line
		}
	}

}

func main() {

	var outputdir string
	var inputfile string
	flag.StringVar(&inputfile, "i", "../../assets/dervaze-rootset.protobuf", "Input file that contains the objects")
	flag.StringVar(&outputdir, "o", "/tmp/rootdata/", "Output dir to write files")

	flag.Parse()

	rootset := loadWordFiles(inputdir)
	dervaze.SaveRootSetProtobuf(outputfile, rootset)

}

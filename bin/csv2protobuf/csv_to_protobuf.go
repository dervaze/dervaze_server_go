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
	"time"
)

// readRecord reads a single line of CSV and returns a Root type
func readRecord(record []string, pos dervaze.PartOfSpeech) *dervaze.Root {
	if record[0][0] == '#' {
		return nil
	}

	if len(record) >= 2 {
		latin := record[0]
		visenc := record[1]
		root := dervaze.NewRoot(latin, visenc, pos)
		return root
	}
	return nil
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

	roots := make([]*dervaze.Root, 0, len(records))

	for i, record := range records {
		if r := readRecord(record, pos); r != nil {
			roots = append(roots, r)
		} else {
			log.Printf("Record error in %s line %d - %s", filename, i, record)
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

func generateSuffixData(rootset *dervaze.RootSet) (*dervaze.RootSet, *dervaze.SuffixSet) {

	// generate map of roots

	rootMap := make(map[string]([]*dervaze.Root), len(rootset.Roots))
	suffixList := make([]int, 0, len(rootset.Roots))
	suffixes := make([]*dervaze.Suffix, 0)

	for i, r := range rootset.Roots {
		if !strings.Contains(r.TurkishLatin, "_") {
			_, exists := rootMap[r.TurkishLatin]
			if !exists {
				rootMap[r.TurkishLatin] = make([]*dervaze.Root, 0, 5)
			}
			rootMap[r.TurkishLatin] = append(rootMap[r.TurkishLatin], r)
		} else {
			suffixList = append(suffixList, i)
		}
	}

	// for each element in the suffix list, generate a suffix record by splitting root

	for _, si := range suffixList {
		s := rootset.Roots[si]
		tr := s.TurkishLatin
		ve := s.Ottoman.Visenc

		trEls := strings.Split(tr, "_")
		if len(trEls) != 2 {
			log.Printf("Skipping suffixes for %s. Multiple or no suffix parts", tr)
			continue
		}

		trRoot := trEls[0]
		trSuffix := trEls[1]

		trRootOrig, exists := rootMap[trRoot]
		if !exists {
			log.Printf("No orig root for %s found", tr)
		}

		// for all visenc candidates, check if any of them has the same prefix with our current visenc word

		for _, tro := range trRootOrig {
			if strings.HasPrefix(ve, tro.Ottoman.Visenc) {
				log.Printf("Found orig root for %s - %s => %s - %s", tr, ve, tro.TurkishLatin, tro.Ottoman.Visenc)
				veSuffix := strings.TrimPrefix(ve, tro.Ottoman.Visenc)
				log.Printf("Suffix Correspondence: %s - %s", trSuffix, veSuffix)

				suffixPOS := s.PartOfSpeech
				suffixRLV := tro.LastVowel
				// we don't have Req_MAYBE here because we need more than one example for that
				var suffixREWV dervaze.Req
				if tro.EndsWithVowel {
					suffixREWV = dervaze.Req_ALWAYS
				} else {
					suffixREWV = dervaze.Req_NEVER
				}

				var suffixRHSV dervaze.Req

				if tro.HasSingleVowel {
					suffixRHSV = dervaze.Req_ALWAYS
				} else {
					suffixRHSV = dervaze.Req_NEVER
				}

				var suffixLCH dervaze.Req
				if tro.LastVowelHard {
					suffixLCH = dervaze.Req_ALWAYS
				} else {
					suffixLCH = dervaze.Req_NEVER
				}

				suffixSLV := s.LastVowel

				suffixCPOS := tro.PartOfSpeech
				suffixEWV := s.EndsWithVowel

				ot, err := dervaze.MakeOttomanWord(veSuffix, "")

				if err == nil {
					suffix := dervaze.Suffix{
						TurkishLatin:              trSuffix,
						Ottoman:                   ot,
						MorphologicalClass:        "auto",
						RequiredLastVowel:         suffixRLV,
						RequiresPOS:               suffixPOS,
						RequiresEndsWithVowel:     suffixREWV,
						RequiresHasSingleVowel:    suffixRHSV,
						RequiresLastConsonantHard: suffixLCH,
						SetsLastVowelTo:           suffixSLV,
						ConvertsPOSto:             suffixCPOS,
						EndsWithVowel:             suffixEWV}
					log.Printf("Adding suffix: %s", suffix.String())
					suffixes = append(suffixes, &suffix)
				} else {
					log.Printf("Error in Ottoman Word: %s", err.Error())

				}

			}

		}

	}
	// delete elements in suffix list from rootset by generating a new rootset

	onlyRoots := make([]*dervaze.Root, 0, len(rootset.Roots))

	for _, rootList := range rootMap {
		onlyRoots = append(onlyRoots, rootList...)
	}

	log.Printf("Generated %d suffixes", len(suffixes))
	log.Printf("Roots now has %d elements", len(onlyRoots))

	suffixSet := dervaze.SuffixSet{Suffixes: suffixes}
	newRootSet := dervaze.RootSet{Roots: onlyRoots}

	return &newRootSet, &suffixSet
}

func main() {

	var inputdir string
	var rootsetfile string
	var suffixsetfile string
	var format string
	t := time.Now().Format("2006-01-02-03-04-05")
	flag.StringVar(&inputdir, "i", "../../assets/rootdata/", "Input dir where n/ v/ p/ directories reside")
	flag.StringVar(&rootsetfile, "r", fmt.Sprintf("../../assets/dervaze-rootset-%s.protobuf", t), "Output file to store the rootset file")
	flag.StringVar(&suffixsetfile, "s", fmt.Sprintf("../../assets/dervaze-suffixset-%s.protobuf", t), "Output file to store the suffixset file")
	flag.StringVar(&format, "f", "protobuf", "Output file to store the suffixset file")

	flag.Parse()

	rootset := loadWordFiles(inputdir)
	newrootset, suffixset := generateSuffixData(rootset)
	if format == "protobuf" {
		dervaze.SaveRootSetProtobuf(rootsetfile, newrootset)
		dervaze.SaveSuffixSetProtobuf(suffixsetfile, suffixset)

	} else if format == "json" {

		dervaze.SaveRootSetJSON(rootsetfile, newrootset)
		dervaze.SaveSuffixSetJSON(suffixsetfile, suffixset)

	} else {
		println("format should be either protobuf or json")
	}

}

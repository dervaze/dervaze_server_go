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
func readRecord(record []string, pos dervaze.PartOfSpeech) *Root {
	if records[0][0] == '#' {
		return nil
	}

	if len(record) >= 2 {
		latin := record[0]
		visenc := record[1]
		root := dervaze.MakeRoot(latin, visenc, pos)
		return &root
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

func generateSuffixData(rootset *dervaze.RootSet) (*dervaze.RootSet, *dervaze.SuffixSet) {

	// generate map of roots

	rootMap := make(map[string][]*Root, 0, len(roots))
	suffixList := make([]int, 0, len(roots))
	suffixes := make([]*Suffix, 0)

	for i, r := range rootset.Roots {
		if !strings.Contains(r.TurkishLatin, "_") {
			ll, exists := rootMap[r.TurkishLatin]
			if !exists {
				rootMap[r.TurkishLatin] = make([]*Root, 0, 5)
			}
			rootMap[r.TurkishLatin] = append(rootMap[r.TurkishLatin], r)
		} else {
			suffixList = append(suffixList, i)
		}

		// for each element in the suffix list, generate a suffix record from corresponding root

		for _, si := range suffixList {
			s := rootset[si]
			tr := s.TurkishLatin
			ve := s.Ottoman.Visenc

			tr_els := strings(tr, "_")
			if len(tr_els) != 2 {
				log.Println("Skipping suffixes for %s. Multiple or no suffix parts", tr)
				continue
			}

			tr_root := tr_els[0]
			tr_suffix := tr_els[1]

			tr_root_orig, exists := rootMap[tr_root]
			if !exists {
				log.Println("No orig root for %s found", tr)
			}

			for j, tro := range tr_root_orig {
				if strings.HasPrefix(ve, tro.Ottoman.Visenc) {
					log.Println("Found orig root for %s - %s => %s - %s", tr, ve, tro.TurkishLatin, tro.Ottoman.Visenc)
					ve_suffix := strings.TrimPrefix(ve, tro.Ottoman.Visenc)
					log.Println("Suffix Correspondence: %s - %s", tr_suffix, ve_suffix)

					suffixPOS := s.PartOfSpeech
					suffixRLV := tro.LastVowel
					// we don't have Req_MAYBE here because we need more than one example for that
					var suffixREWV Req
					if tro.EndsWithVowel {
						suffixREWV = dervaze.Req_ALWAYS
					} else {
						suffixREWV = dervaze.Req_NEVER
					}

					var suffixRHSV Req

					if tro.HasSingleVowel {
						suffixRHSV = dervaze.Req_ALWAYS
					} else {
						suffixRHSV = dervaze.Req_NEVER
					}

					var suffixLCH Req
					if tro.LastVowelHard {
						suffixLCH = dervaze.Req_ALWAYS
					} else {
						suffixLCH = dervaze.Req_NEVER
					}

					suffixSLV := s.LastVowel

					suffixCPOS := tro.PartOfSpeech
					suffixEWV := s.EndsWithVowel

					suffix := Suffix{
						TurkishLatin:              tr_suffix,
						Ottoman:                   dervaze.MakeOttomanWord(ve_suffix, ""),
						MorphologicalClass:        "auto",
						RequiredLastVowel:         suffixRLV,
						RequiresPOS:               suffixPOS,
						RequiresEndsWithVowel:     suffixREWV,
						RequiresHasSingleVowel:    suffixRHSV,
						RequiresLastConsonantHard: suffixLCH,
						SetsLastVowelTo:           suffixSLV,
						ConvertsPOSto:             suffixCPOS,
						EndsWithVowel:             suffixEWV}
				}

				log.Println("Adding suffix: %s", string(suffix))
				suffixes = append(suffixes, &suffix)
			}

		}

		// delete elements in suffix list from rootset by generating a new rootset

		only_roots := make([]*Root, 0, len(rootset.Roots))

		for _, root_list := range rootMap {
			only_roots = append(only_roots, root_list...)
		}

		log.Println("Generated %d suffixes", len(suffixes))
		log.Println("Roots now has %d elements", len(only_roots))

		suffixSet := SuffixSet{Suffixes: suffixes}
		newRootSet := RootSet{Roots: only_roots}

		return &newRootSet, &suffixSet

	}

}

func main() {

	var inputdir string
	var rootsetfile string
	var suffixsetfile string
	flag.StringVar(&inputdir, "i", "../../assets/rootdata/", "Input dir where n/ v/ p/ directories reside")
	flag.StringVar(&rootsetfile, "r", fmt.Sprintf("../../assets/dervaze-rootset-%s.protobuf", ""), "Output file to store the rootset file")
	flag.StringVar(&suffixsetfile, "s", fmt.Sprintf("../../assets/dervaze-suffixset-%s.protobuf", ""), "Output file to store the suffixset file")

	flag.Parse()

	rootset := loadWordFiles(inputdir)
	suffixset, newrootset := generateSuffixData(rootset)
	dervaze.SaveRootSetProtobuf(rootsetfile, newrootset)
	dervaze.SaveSuffixSetProtobuf(suffixsetfile, suffixset)

}

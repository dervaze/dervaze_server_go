package main

import (
	dervaze "dervaze/lang"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	// "time"
)

func csvLine(root *dervaze.Root) string {
	line := root.TurkishLatin + "," + root.Ottoman.Visenc + "," + root.Ottoman.Unicode + "\n"
	return line
}

func WriteRootSetToCSV(inputfile string, outputdir string) {
	dervaze.InitSearch(inputfile)
	rootSet := dervaze.GetRootSet()

	verbFiles := make(map[string]string)
	nounFiles := make(map[string]string)
	properFiles := make(map[string]string)
	otherFiles := make(map[string]string)

	for _, r := range rootSet.Roots {
		line := csvLine(r)
		lineRune := []rune(line)
		firstRune := lineRune[0]
		key := string(firstRune)
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

	verbDir := filepath.Join(outputdir, "v")
	nounDir := filepath.Join(outputdir, "n")
	properDir := filepath.Join(outputdir, "p")
	otherDir := filepath.Join(outputdir, "o")

	log.Println(verbDir)
	log.Println(nounDir)
	log.Println(properDir)
	log.Println(otherDir)

	err := os.MkdirAll(verbDir, 0777)
	if err != nil {
		log.Println(err)
	}
	err = os.MkdirAll(nounDir, 0777)
	if err != nil {
		log.Println(err)
	}
	err = os.MkdirAll(properDir, 0777)
	if err != nil {
		log.Println(err)
	}
	err = os.MkdirAll(otherDir, 0777)
	if err != nil {
		log.Println(err)
	}

	writeFunc := func(m map[string]string, outdir string, fileprefix string) {
		for name, content := range m {
			filename := outdir + "/" + fileprefix + name + ".csv"
			log.Println(filename)

			filelines := strings.Split(content, "\n")
			sort.Strings(filelines)
			filecontent := strings.Join(filelines, "\n")
			err := ioutil.WriteFile(filename, []byte(filecontent), 0666)
			if err != nil {
				log.Println(err)
			}
		}
	}

	writeFunc(nounFiles, nounDir, "noun-")
	writeFunc(verbFiles, verbDir, "verb-")
	writeFunc(properFiles, properDir, "proper-")
	writeFunc(otherFiles, properDir, "other-")

}

func main() {

	var outputdir string
	var inputfile string
	flag.StringVar(&inputfile, "i", "../../assets/dervaze-rootset.protobuf", "Input file that contains the objects")
	flag.StringVar(&outputdir, "o", "/tmp/rootdata/", "Output dir to write files")

	flag.Parse()

	WriteRootSetToCSV(inputfile, outputdir)

}

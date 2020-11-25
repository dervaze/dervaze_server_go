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
	"time"
)

func readCSVFile(filename string, pos dervaze.PartOfSpeech, roots []dervaze.Root, rooti *int) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
	}

	csvr := csv.NewReader(strings.NewReader(string(data)))
	records, err := csvr.ReadAll()
	if err != nil {
		log.Println(err)
	}

	for i, record := range records {
		if len(record) == 2 {
			latin := record[0]
			visenc := record[1]
			roots[*rooti] = dervaze.MakeRoot(latin, visenc, pos)
			(*rooti)++
		} else {
			log.Println("Record error in %s line %d - %s", filename, i, record)

		}
	}
}

func loadWordFiles() []dervaze.Root {

	roots := make([]dervaze.Root, 0, 250000)
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
		readCSVFile(fn, dervaze.PartOfSpeech_VERB, roots[:], &rooti)
	}
	for _, fn := range nounFiles {
		readCSVFile(fn, dervaze.PartOfSpeech_NOUN, roots[:], &rooti)
	}

	for _, fn := range properFiles {
		readCSVFile(fn, dervaze.PartOfSpeech_PROPER_NOUN, roots[:], &rooti)
	}

	log.Println("Read %d records", rooti)

	return roots[:]

}

func storeProtobuf(roots []dervaze.Root) {
	t := time.Now().Format("2006-01-02-03-04-05")
	filename := fmt.Sprintf("dervaze-roots-%s.bin", t)
	file, err := os.OpenFile(
		filename,
		os.O_APPEND|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write bytes to file
	totalBytes := 0
	for i, r := range roots {
		byteSlice, err := proto.Marshal(r)
		if err != nil {
			log.Fatal(err)
		}
		bytesWritten, err := file.Write(byteSlice)
		if err != nil {
			log.Fatal(err)
		}
		totalBytes += bytesWritten
	}
	log.Printf("%s: Wrote %d bytes.\n", filename, totalBytes)
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
	storeProtobuf(roots)

}

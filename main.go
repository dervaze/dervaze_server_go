package main

import (
	"dervaze"
	"encoding/csv"
	"fmt"
	"path/filepath"
)

func loadWordFiles() {
	verbFiles, err := filepath.Glob("assets/csv/v/*.csv")
	if err != nil {
		verbFiles = make([]string, 0)
	}
	nounFiles, err := filepath.Glob("assets/csv/n/*.csv")
	if err != nil {
		nounFiles = make([]string, 0)
	}
	properFiles, err := filepath.Glob("assets/csv/p/*.csv")
	if err != nil {
		properFiles = make([]string, 0)
	}

	for i, vfn := range verbFiles {

	}

	io.FileReader()
	reader := csv.NewReader(file)
}

func main() {

	dervaze.MakeRoot("emre", "emrh")
	ow := dervaze.OttomanWord{

		Visenc:           "aa",
		Unicode:          "bb",
		Abjad:            123,
		VisencLetters:    []string{"a", "a"},
		SearchKey:        "aa",
		DotlessSearchKey: "dd",
	}
	fmt.Println(ow.Abjad)
}

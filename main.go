package main

import (
	dervazepb "deneme/dervaze.com/dervazepb"
	"encoding/csv"
	"fmt"
	"path/filepath"
)

func makeRoot(latin string, visenc string, pos dervazepb.PartOfSpeech) dervazepb.Root
{

	ow := OttomanWord{
		Visenc: visenc,
		Unicode: VisencToUnicode(visenc),
		Abjad: VisencToAbjad(visenc),
		VisencLetters: VisencLetters(visenc),
		SearchKey: VisencToSearchKey(visenc),
		DotlessSearchKey: VisencToDotlessSearchKey(visenc)
	}

	r := dervazepb.Root{
		TurkishLatin: latin,
		Ottoman: ow,
		LastVowel: LatinLastVowel(latin),
		LastConsonant: LatinLastConsonant(latin),
		EffectiveLastVowel: LatinEffectiveLastVowel(latin),
		EffectiveTurkishLatin: LatinEffectiveTurkishLatin(latin),
		EffectiveVisenc: LatinEffectiveVisenc,
		Abjad: ow.Abjad,
		PartOfSpeech: pos,
		EndsWithVowel: LatinEndsWithVowel(latin),
		HasSingleVowel: LatinHasSingleVowel(latin),
		LastVowelHard: LatinLastVowelHard(latin),
		LastConsonantHard: LatinLastConsonantHard(latin),
		HasConsonantSoftening: LatinHasConsonantSoftening(latin)
	}
}

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

	ow := dervazepb.OttomanWord{

		Visenc:           "aa",
		Unicode:          "bb",
		Abjad:            123,
		VisencLetters:    []string{"a", "a"},
		SearchKey:        "aa",
		DotlessSearchKey: "dd",
	}
	fmt.Println(ow.Abjad)
}

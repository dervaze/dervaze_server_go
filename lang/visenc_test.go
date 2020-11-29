package lang

import (
	"fmt"
	"testing"
)

/*

Characteristics of a Golang test function:

    The first and only parameter needs to be t *testing.T
    It begins with the word Test followed by a word or phrase starting with a capital letter.
    (usually the method under test i.e. TestValidateClient)
    Calls t.Error or t.Fail to indicate a failure (I called t.Errorf to provide more details)
    t.Log can be used to provide non-failing debug information
    Must be saved in a file named something_test.go such as: addition_test.go
*/

// func SearchKey(s string) string {
func TestSearchKey(t *testing.T) {
	testDict := map[string]string{
		"elfo1":        "elfo1",
		"emrh":         "emrh",
		"eo5mrh":       "emrh",
		"eo6mro1h":     "emro1h",
		"eo5so3ro8bu2": "eso3rbu2"}

	for i, o := range testDict {
		if SearchKey(i) != o {
			t.Log(fmt.Sprintf("%s, %s fails for SearchKey", i, o))
			t.Fail()
		}
	}
}

// func DotlessSearchKey(s string) string {
func TestDotlessSearchKey(t *testing.T) {
	testDict := map[string]string{
		"elfo1":        "elf",
		"emrh":         "emrh",
		"eo5mrh":       "emrh",
		"eo6mro1h":     "emrh",
		"eo5so3ro8bu2": "esrb"}

	for i, o := range testDict {
		if DotlessSearchKey(i) != o {
			t.Log(fmt.Sprintf("%s, %s fails for DotlessSearchKey", i, o))
			t.Fail()
		}
	}
}

// func EndsWithVowel(s string) bool {
func TestEndsWithVowel(t *testing.T) {
	testDict := map[string]bool{
		"emre":   true,
		"esra":   true,
		"araba":  true,
		"meydan": false,
		"sarâ":   true,
		"dev":    false,
	}

	for i, o := range testDict {
		if EndsWithVowel(i) != o {
			t.Log(fmt.Sprintf("%s, %t fails for EndsWithVowel", i, o))
			t.Fail()
		}
	}
}

// func SplitVisenc(s string, addInvalidLetters bool) []string {
func TestSplitVisenc(t *testing.T) {

	testDict := map[string][]string{
		"emrh":     {"e", "m", "r", "h"},
		"efo2bu2x": {"e", "fo2", "bu2", "x"},
		"aübo1eo5": {"a", "ü", "bo1", "eo5"},
		"fo3d":     {"f", "o", "3", "d"},
		"brdh":     {"b", "r", "d", "h"},
	}
	for i, o := range testDict {
		if CompareStringSlices(SplitVisenc(i, true), o) == false {
			t.Log(fmt.Sprintf("SplitVisenc(%s, true) returns %s instead of %s", i, SplitVisenc(i, true), o))
			t.Fail()
		}
	}
}

// func MakeOttomanWord(visenc string, unicode string) (*OttomanWord, error) {
func TestMakeOttomanWord(t *testing.T) {

	ow, err := MakeOttomanWord("emre", "emrh")

	if err != nil {
		t.Fail()
		t.Log(err)
	}

	t.Log(ow)
}

// func VisencToAbjad(s string) int32 {
func TestVisencToAbjad(t *testing.T) {
	testDict := map[string]int32{
		"emrh":      246,
		"mlk":       90,
		"ewao1wro1": 1020}

	for i, o := range testDict {
		res := VisencToAbjad(i)
		if res != o {
			t.Log(fmt.Sprintf("VisencToAbjad(%s) gives %d. It should be %d", i, res, o))
			t.Fail()
		}
	}
}

// func MakeRoot(latin string, visenc string, pos PartOfSpeech) Root {
// func TestMakeRoot(t *testing.T) { t.Fail() }

// func VisencToUnicode(s string) string {
func TestVisencToUnicode(t *testing.T) {

	testDict := map[string]string{
		"eo6bu1dsbo2sro1":             "آبدستسز",
		"axu1mbu2ebo1":                "عجميان",
		"eo6xo1erlemfo2":              "آخارلامق",
		"so3o4ro0fo2u4 so3u4mo4elbu2": "شَرْقِ شِمَالي",
		"bo1o4to1bu2fo1o0":            "نَظيفْ",
		"eu4ko0ro4emu4 eu4lo4hu4bu2":  "اِکْرَامِ اِلَهِي",
		"klmh||lr":                    "کلمه‌لر",
	}

	for v, u := range testDict {
		if r := VisencToUnicode(v); r != u {
			t.Log(fmt.Sprintf("Fails for %s -> %s - returned %s", v, u, r))
			t.Fail()
		}
	}

}

// func UnicodeToVisenc(s string) string {
func TestUnicodeToVisenc(t *testing.T) {

	testDict := map[string]string{
		"eo6bu1dsbo2sro1":             "آبدستسز",
		"axu1mbu2ebo1":                "عجميان",
		"eo6xo1erlemfo2":              "آخارلامق",
		"so3o4ro0fo2u4 so3u4mo4elbu2": "شَرْقِ شِمَالي",
		"bo1o4to1bu2fo1o0":            "نَظيفْ",
		"eu4ko0ro4emu4 eu4lo4hu4bu2":  "اِکْرَامِ اِلَهِي",
		"klmh||lr":                    "کلمه‌لر",
	}

	for v, u := range testDict {
		if r := UnicodeToVisenc(u); r != v {
			t.Log(fmt.Sprintf("Fails for %s -> %s - returned %s", u, v, r))
			t.Fail()
		}
	}

}

// func UnicodeToAbjad(s string) int32 {
// func TestUnicodeToAbjad(t *testing.T) { t.Fail() }

// func ContainsArabicChars(s string) bool {
func TestContainsArabicChars(t *testing.T) {

	testDict := map[string]bool{
		"eo6bu1dsbo2sro1":             false,
		"آبدستسز":                     true,
		"axu1mbu2ebo1":                false,
		"عجميان":                      true,
		"eo6xo1erlemfo2":              false,
		"آخارلامق":                    true,
		"so3o4ro0fo2u4 so3u4mo4elbu2": false,
		"شَرْقِ شِمَالي":              true,
		"bo1o4to1bu2fo1o0نَظيفْ":      true,
		"eu4ko0ro4emu4 4bu2اِکْرَامِ اِلَهِي": true,
		"klmh||lr  کلمه‌لر":                   true,
	}

	for s, b := range testDict {
		if r := ContainsArabicChars(s); r != b {
			t.Log(fmt.Sprintf("Fails for %s -> %t", s, b))
			t.Fail()
		}
	}

}

// func ContainsDigits(s string) bool {
func TestContainsDigits(t *testing.T) {
	testDict := map[string]bool{
		"eo6bu1dsbo2sro1":             true,
		"آبدستسز":                     false,
		"axu1mbu2ebo1":                true,
		"عجميان":                      false,
		"eo6xo1erlemfo2":              true,
		"آخارلامق":                    false,
		"so3o4ro0fo2u4 so3u4mo4elbu2": true,
		"شَرْقِ شِمَالي":              false,
		"bo1o4to1bu2fo1o0نَظيفْ":      true,
		"eu4ko0ro4emu4 4bu2اِکْرَامِ اِلَهِي": true,
		"klmh||lr  کلمه‌لر":                   false,
	}

	for s, b := range testDict {
		if r := ContainsDigits(s); r != b {
			t.Log(fmt.Sprintf("Fails for %s -> %t", s, b))
			t.Fail()
		}
	}

}

// func HasSingleVowel(s string) bool {
func TestHasSingleVowel(t *testing.T) {
	testDict := map[string]bool{
		"cevat":   false,
		"merâ":    false,
		"semer":   false,
		"ezeli":   false,
		"cemaat":  false,
		"umur":    false,
		"said":    false,
		"çalış":   false,
		"kaynak":  false,
		"gol":     true,
		"saygı":   false,
		"fıstık":  false,
		"trüf":    true,
		"bilinç":  false,
		"bilinc":  false,
		"birim":   false,
		"aşk":     true,
		"dert":    true,
		"kelebek": false,
		"baal":    false,
	}

	for w, sv := range testDict {
		if r := HasSingleVowel(w); r != sv {
			t.Log(fmt.Sprintf("Fails for %s -> %t", w, sv))
			t.Fail()
		}
	}
}

// func LastConsonantHard(s string) bool {
func TestLastConsonantHard(t *testing.T) {

	testDict := map[string]bool{
		"cevat":  true,
		"merâ":   false,
		"semer":  false,
		"ezeli":  false,
		"cemaat": true,
		"umur":   false,
		"said":   false,
		"çalış":  true,
		"kaynak": true,
		"gol":    false,
		"saygı":  false,
		"fıstık": true,
		"trüf":   true,
		"bilinç": true,
		"bilinc": false,
		"birim":  false,
	}

	for w, lch := range testDict {
		if r := LastConsonantHard(w); r != lch {
			t.Log(fmt.Sprintf("Fails for %s -> %t", w, lch))
			t.Fail()
		}
	}

}

// func LastVowelHard(s string) bool {
func TestLastVowelHard(t *testing.T) {

	testDict := map[string]bool{
		"cevat":  true,
		"merâ":   true,
		"semer":  false,
		"ezeli":  false,
		"cemaat": true,
		"umur":   true,
		"said":   false,
		"çalış":  true,
		"kaynak": true,
		"gol":    true,
		"saygı":  true}

	for w, lv := range testDict {
		if r := LastVowelHard(w); r != lv {
			t.Log(fmt.Sprintf("Fails for %s -> %t", w, lv))
			t.Fail()
		}
	}
}

// func EffectiveLastVowel(s string) string {
func TestEffectiveLastVowel(t *testing.T) { t.Fail() }

// func LastVowel(s string) string {
func TestLastVowel(t *testing.T) {

	testDict := map[string]string{
		"cevat":  "a",
		"merâ":   "â",
		"semer":  "e",
		"ezeli":  "i",
		"cemaat": "a"}

	for w, lv := range testDict {
		if r := LastVowel(w); r != lv {
			t.Log(fmt.Sprintf("Fails for %s -> %s -- returns %s", w, lv, r))
			t.Fail()
		}
	}

}

// func LastConsonant(s string) string {
func TestLastConsonant(t *testing.T) {

	testDict := map[string]string{
		"cevat":  "t",
		"mera":   "r",
		"semer":  "r",
		"ezeli":  "l",
		"cemaat": "t"}

	for w, lc := range testDict {
		if r := LastConsonant(w); r != lc {
			t.Log(fmt.Sprintf("Fails for %s -> %s -- returns %s", w, lc, r))
			t.Fail()
		}
	}
}

// func UpdateEffectiveSoftening(r *Root) {
func TestUpdateEffectiveSoftening(t *testing.T) { t.Fail() }

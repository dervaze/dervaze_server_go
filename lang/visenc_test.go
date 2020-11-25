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

func CompareSlicesString(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i, v := range slice1 {
		if slice2[i] != v {
			return false
		}
	}
	return true
}

func TestSplitVisenc(t *testing.T) {

	testDict := map[string][]string{
		"emrh":     []string{"e", "m", "r", "h"},
		"efo2bu2x": []string{"e", "fo2", "bu2", "x"},
		"aübo1eo5": []string{"a", "ü", "bo1", "eo5"},
		"fo3d":     []string{"f", "o", "3", "d"},
		"brdh":     []string{"b", "r", "d", "h"},
	}
	for i, o := range testDict {
		if CompareSlicesString(SplitVisenc(i, true), o) == false {
			t.Log(fmt.Sprintf("SplitVisenc(%s, true) returns %s instead of %s", i, SplitVisenc(i, true), o))
			t.Fail()
		}
	}
}

func TestMakeOttomanWord(t *testing.T) {

	ow, err := MakeOttomanWord("emre", "emrh")

	if err != nil {
		t.Fail()
		t.Log(err)
	}

	t.Log(ow)
}

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

package dervaze

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
			t.Log(fmt.Sprintf("%s, %s fails for EndsWithVowel", i, o))
			t.Fail()
		}
	}
}

func TestSplitVisenc(t *testing.T) {

	testDict := map[string][]string{
		"emre":     []string{"e", "m", "r", "h"},
		"efo2bu2x": []string{"e", "fo2", "bu2", "x"},
		"ağbo1eo5": []string{"a", "ğ", "bo1", "eo5"},
		"fo3d":     []string{"f", "o", "3", "d"},
		"brdh":     []string{"b", "r", "d", "h"},
	}
	for i, o := range testDict {
		if SplitVisenc(i) != o {
			t.Log(fmt.Sprintf("%s, %s fails for SplitVisenc", i, o))
			t.Fail()
		}
	}
}

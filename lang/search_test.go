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

const PROTOBUFFILE = "../assets/dervaze-rootset.protobuf"

// func InitSearch(protobuffile string) {
func TestInitSearch(t *testing.T) {
	fmt.Println("InitSearch")
	InitSearch(PROTOBUFFILE)
}

// func GetRootSet() *RootSet {
func TestGetRootSet(t *testing.T) {
	rootSet := GetRootSet()

	if len(rootSet.Roots) == 0 {
		t.Log("len(rootset) is 0")
		t.Fail()
	}
}

// func PrefixSearchTurkishLatin(turkishLatin string) []*Root {

func TestPrefixSearchTurkishLatin(t *testing.T) {

}

// func PrefixSearchTurkishLatinExact(turkishLatin string) []*Root {
// func PrefixSearchVisenc(visenc string) []*Root {
// func PrefixSearchVisencExact(visenc string) []*Root {
// func PrefixSearchUnicode(unicode string) []*Root {
// func PrefixSearchUnicodeExact(unicode string) []*Root {
// func PrefixSearchAbjad(abjad int32) []*Root {
// func PrefixSearchAll(term string) []*Root {
// func splitRootKeyFromIndex(k string) (int, error) {
// func indexRegexSearch(keylist []string, r regex.Regex) []*Root {
// func RegexSearchTurkishLatin(word string) []*Root {
// func RegexSearchUnicode(word string) []*Root {
// func RegexSearchVisenc(word string) []*Root {
// func IndexSearchAbjad(abjad int32) []*Root {
// func PrintRoots(roots []*Root) string {

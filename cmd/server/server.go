package main

import (
	dervaze "dervaze/lang"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
)

// ## `/v1/json/prefix/tr/{word}
//
// Sends a list of Turkish words starting with `word` sorted by length
//
// ```
// [ "word", "worda", "wordb", "wordabc"]
// ```
//

func JsonPrefixTr(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("JsonPrefixTr Vars: %s", vars)
	roots := dervaze.SearchTurkishLatin(vars["word"])
	log.Printf("roots: %s", roots)
	jsonBytes, err := proto.MarshalMessageSetJSON(roots)
	if err == nil {
		fmt.Fprintln(w, "", jsonBytes)
	} else {
		log.Fatal(err)
	}
}

// ## `/v1/json/prefix/ot/{word}
//
// Sends a list of Ottoman words starting with `word`
//
// ```
// [ "word", "worda", "wordb", "wordabc"]
// ```
//
// ## `/v1/json/calc/abjad/{word}
//
// Calculates abjad for the `word` given in unicode
//
// ```
// { "ottoman_unicode": "word",
// "abjad": 1234 }
// ```
//
// ## `/v1/json/exact/tr/{word}
//
// Returns records with Turkish Latin == `word`
//
// ```
// [ { "ottoman_unicode": "abcd",
//     "latin": "word",
//     "abjad": 1234,
//     "meanings": [
//         {"source": "dervaze",
//          "meaning": "meaning 1 of word"}]
//     },
//     { "ottoman_unicode": "abcddd",
//     "latin": "word",
//     "abjad": 1255,
//     "meanings": [
//         {"source": "kanar",
//          "meaning": "meaning 1 of word"}]
//     }]
//     ```
//
// ## `/v1/json/exact/ot/{word}
//
// Returns records with Ottoman == `word`
//
// ```
// [ { "ottoman_unicode": "word",
//     "latin": "abdc",
//     "abjad": 1298,
//     "meanings": [
//         {"source": "dervaze",
//          "meaning": "meaning 1 of word"}]
//     },
//     { "ottoman_unicode": "word",
//     "latin": "anbdn",
//     "abjad": 2191,
//     "meanings": [
//         {"source": "kanar",
//          "meaning": "meaning 1 of word"}]
//     }]
//     ```
//
// ## `/v1/json/exact/abjad/{number}
//
// Returns records with abjad == `number`
//
// ```
// [ { "ottoman_unicode": "word",
//     "latin": "abdc",
//     "abjad": number,
//     "meanings": [
//         {"source": "dervaze",
//          "meaning": "meaning 1 of word"}]
//     },
//     { "ottoman_unicode": "drwo",
//     "latin": "anbdn",
//     "abjad": number,
//     "meanings": [
//         {"source": "kanar",
//          "meaning": "meaning 1 of word"}]
//     }]
//     ```
//
//
// ## `/v1/json/v2u/{word}
//
// Converts `word` from visenc to unicode
//
// ```
// {"ottoman_visenc": <word>,
// "ottoman_unicode": <unicode>}
// ```
//
// ## `/v1/json/u2v/{word}
//
// Converts `word` from unicode to visenc
//
// ```
// {"ottoman_visenc": <visenc>,
// "ottoman_unicode": <word>}
// ```
//

func server() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/v1/json/prefix/tr/{word}", JsonPrefixTr)
	router.HandleFunc("/v1/json/prefix/ot/{word}", JsonPrefixOt)
	router.HandleFunc("/v1/json/calc/abjad/{word}", JsonCalcAbjad)
	router.HandleFunc("/v1/json/exact/tr/{word}", JsonExactTr)
	router.HandleFunc("/v1/json/exact/ot/{word}", JsonExactOt)
	router.HandleFunc("/v1/json/exact/abjad/{number}", JsonExactAbjad)
	router.HandleFunc("/v1/json/v2u/{word}", JsonV2U)
	router.HandleFunc("/v1/json/u2v/{word}", JsonU2V)

}

func main() {

	var inputfile string
	flag.StringVar(&inputfile, "i", "assets/dervaze-rootset.protobuf", "protobuffer file to load roots")

	flag.Parse()
	dervaze.InitSearch(inputfile)
	server()
}

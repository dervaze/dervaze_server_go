package main

import (
	dervaze "dervaze/lang"
	// "encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
	// "github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	// "google.golang.org/protobuf/encoding/protojson"
	// "google.golang.org/protobuf/proto"
)

func transformRoots(roots []*dervaze.Root, transformer func(*dervaze.Root) *dervaze.Root) *dervaze.RootSet {
	out := make([]*dervaze.Root, len(roots))

	for i, r := range roots {
		out[i] = transformer(r)
	}

	r := dervaze.RootSet{
		Roots: out,
	}
	return &r
}

// ## `/v1/json/prefix/tr/{word}
//
// Sends a list of Turkish words starting with `word` sorted by length
//
// ```
// [ "word", "worda", "wordb", "wordabc"]
// ```
//

func JsonPrefixTr(w http.ResponseWriter, r *http.Request) {

	transformer := func(root *dervaze.Root) *dervaze.Root {
		r := dervaze.Root{
			TurkishLatin: root.TurkishLatin,
		}
		return &r
	}
	vars := mux.Vars(r)
	log.Printf("JsonPrefixTr Vars: %s", vars)
	roots := dervaze.SearchTurkishLatin(vars["word"])
	log.Printf("roots: %s", roots)

	outputRootSet := transformRoots(roots, transformer)
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
		Indent:       "  ",
	}
	jsonStr, err := marshaler.MarshalToString(outputRootSet)

	if err == nil {
		fmt.Fprintln(w, "", jsonStr)
	} else {
		log.Printf("Marshal Error: %s", err)
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

func JsonPrefixOt(w http.ResponseWriter, r *http.Request) {
	transformer := func(root *dervaze.Root) *dervaze.Root {
		r := dervaze.Root{
			Ottoman: OttomanWord{
				Unicode: root.Ottoman.Unicode 
			}
		}
		return &r
	}
	vars := mux.Vars(r)
	log.Printf("JsonPrefixTr Vars: %s", vars)
	roots := dervaze.SearchUnicode(vars["word"])
	log.Printf("roots: %s", roots)

	outputRootSet := transformRoots(roots, transformer)
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
		Indent:       "  ",
	}
	jsonStr, err := marshaler.MarshalToString(outputRootSet)

	if err == nil {
		fmt.Fprintln(w, "", jsonStr)
	} else {
		log.Printf("Marshal Error: %s", err)
	}
}

// ## `/v1/json/calc/abjad/{word}
//
// Calculates abjad for the `word` given in unicode
//
// ```
// { "ottoman_unicode": "word",
// "abjad": 1234 }
// ```
//
func JsonCalcAbjad(w http.ResponseWriter, r *http.Request) {

	
	vars := mux.Vars(r)
	log.Printf("JsonPrefixTr Vars: %s", vars)
	abjad := dervaze.UnicodeToAbjad(vars["word"])
	str := ftm.Sprintf("{ \"ottoman_unicode\": \"%s\", \"abjad\": %d }", vars["word"], abjad)
	fmt.Fprintln(w, "", str)

}

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
func JsonExactTr(w http.ResponseWriter, r *http.Request) {
}

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
func JsonExactOt(w http.ResponseWriter, r *http.Request) {}

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
func JsonExactAbjad(w http.ResponseWriter, r *http.Request) {}

// ## `/v1/json/v2u/{word}
//
// Converts `word` from visenc to unicode
//
// ```
// {"ottoman_visenc": <word>,
// "ottoman_unicode": <unicode>}
// ```
//
func JsonV2U(w http.ResponseWriter, r *http.Request) {}

// ## `/v1/json/u2v/{word}
//
// Converts `word` from unicode to visenc
//
// ```
// {"ottoman_visenc": <visenc>,
// "ottoman_unicode": <word>}
// ```
//
func JsonU2V(w http.ResponseWriter, r *http.Request) {}

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

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:9876",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

func main() {

	var inputfile string
	flag.StringVar(&inputfile, "i", "../../assets/dervaze-rootset.protobuf", "protobuffer file to load roots")

	flag.Parse()
	dervaze.InitSearch(inputfile)
	server()
}

package main

import (
	dervaze "dervaze/lang"
	"os/exec"
	"strconv"
	"strings"

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

func marshalRoots(outputRootSet *dervaze.RootSet) (string, error) {
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
		Indent:       "  ",
	}
	jsonStr, err := marshaler.MarshalToString(outputRootSet)

	if err == nil {
		return jsonStr, nil
	}
	log.Printf("Marshal Error: %s", err)
	return "", err
}

// JSONPrefixTr makes a prefix search with the word
// ## `/v1/json/prefix/tr/{word}
//
// Sends a list of Turkish words starting with `word` sorted by length
//
func JSONPrefixTr(w http.ResponseWriter, r *http.Request) {

	transformer := func(root *dervaze.Root) *dervaze.Root {
		r := dervaze.Root{
			TurkishLatin: root.TurkishLatin,
		}
		return &r
	}
	vars := mux.Vars(r)
	log.Printf("JsonPrefixTr Vars: %s", vars)
	roots := dervaze.RegexSearchTurkishLatin(vars["word"])
	log.Printf("roots: %s", roots)

	outputRootSet := transformRoots(roots, transformer)

	if m, err := marshalRoots(outputRootSet); err == nil {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, "", m)
	}
}

// JSONPrefixOt responds to a prefix search
// ## `/v1/json/prefix/ot/{word}
//
// Sends a list of Ottoman words starting with `word`
//
// ```
// [ "word", "worda", "wordb", "wordabc"]
// ```
//
func JSONPrefixOt(w http.ResponseWriter, r *http.Request) {
	transformer := func(root *dervaze.Root) *dervaze.Root {
		r := dervaze.Root{
			Ottoman: &dervaze.OttomanWord{
				Unicode: root.Ottoman.Unicode,
			},
		}
		return &r
	}
	vars := mux.Vars(r)
	log.Printf("JsonPrefixTr Vars: %s", vars)
	roots := dervaze.RegexSearchUnicode(vars["word"])
	log.Printf("roots: %s", roots)

	outputRootSet := transformRoots(roots, transformer)

	if m, err := marshalRoots(outputRootSet); err == nil {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, "", m)
	}
}

// JSONExactTr searches `word` exactly without prefix or regex
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
func JSONExactTr(w http.ResponseWriter, r *http.Request) {
	transformer := func(root *dervaze.Root) *dervaze.Root {
		r := dervaze.Root{
			TurkishLatin: root.TurkishLatin,
			Abjad:        root.Abjad,
			Ottoman: &dervaze.OttomanWord{
				Unicode: root.Ottoman.Unicode,
			},
		}
		return &r
	}
	vars := mux.Vars(r)
	log.Printf("JsonExactTr Vars: %s", vars)
	roots := dervaze.PrefixSearchTurkishLatinExact(vars["word"])
	log.Printf("roots: %s", roots)

	outputRootSet := transformRoots(roots, transformer)
	if m, err := marshalRoots(outputRootSet); err == nil {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, "", m)
	}
}

// JSONExactOt searches `word` exactly as written without prefix or regex
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
func JSONExactOt(w http.ResponseWriter, r *http.Request) {
	transformer := func(root *dervaze.Root) *dervaze.Root {
		r := dervaze.Root{
			TurkishLatin: root.TurkishLatin,
			Abjad:        root.Abjad,
			Ottoman: &dervaze.OttomanWord{
				Unicode: root.Ottoman.Unicode,
			},
		}
		return &r
	}
	vars := mux.Vars(r)
	log.Printf("JsonExactTr Vars: %s", vars)
	roots := dervaze.PrefixSearchUnicodeExact(vars["word"])
	log.Printf("roots: %s", roots)

	outputRootSet := transformRoots(roots, transformer)
	if m, err := marshalRoots(outputRootSet); err == nil {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, "", m)
	}
}

// JSONSearchTr makes a regex search by interleaving .? between runes of `word`
// `/v1/json/search/tr/{word}`
func JSONSearchTr(w http.ResponseWriter, r *http.Request) {
	transformer := func(root *dervaze.Root) *dervaze.Root {
		r := dervaze.Root{
			TurkishLatin: root.TurkishLatin,
			Abjad:        root.Abjad,
			Ottoman: &dervaze.OttomanWord{
				Unicode: root.Ottoman.Unicode,
			},
		}
		return &r
	}
	vars := mux.Vars(r)
	log.Printf("JsonExactTr Vars: %s", vars)
	roots := dervaze.RegexSearchTurkishLatin(vars["word"])
	log.Printf("roots: %s", roots)

	outputRootSet := transformRoots(roots, transformer)

	if m, err := marshalRoots(outputRootSet); err == nil {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, "", m)
	}
}

// JSONSearchOt makes a regex search by interleaving .? between runes of `word`
// `/v1/json/search/ot/{word}`
func JSONSearchOt(w http.ResponseWriter, r *http.Request) {
	transformer := func(root *dervaze.Root) *dervaze.Root {
		r := dervaze.Root{
			TurkishLatin: root.TurkishLatin,
			Abjad:        root.Abjad,
			Ottoman: &dervaze.OttomanWord{
				Unicode: root.Ottoman.Unicode,
			},
		}
		return &r
	}
	vars := mux.Vars(r)
	log.Printf("JsonExactTr Vars: %s", vars)
	roots := dervaze.RegexSearchUnicode(vars["word"])
	log.Printf("roots: %s", roots)

	outputRootSet := transformRoots(roots, transformer)
	if m, err := marshalRoots(outputRootSet); err == nil {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, "", m)
	}
}

// JSONSearchAuto makes a regex search by interleaving .? between runes of `word`
// if `word` contains Arabic letters, it makes Arabic search
// if `word` contains digits only it makes an abjad search
// if `word` contains characters and digits mixed, it makes a visenc search
// otherwise it searches as a Turkish latin word
// `/v1/json/search/ot/{word}`
func JSONSearchAuto(w http.ResponseWriter, r *http.Request) {
	transformer := func(root *dervaze.Root) *dervaze.Root {
		r := dervaze.Root{
			TurkishLatin: root.TurkishLatin,
			Abjad:        root.Abjad,
			Ottoman: &dervaze.OttomanWord{
				Unicode: root.Ottoman.Unicode,
			},
		}
		return &r
	}
	vars := mux.Vars(r)
	log.Printf("JsonExactTr Vars: %s", vars)
	roots := dervaze.RegexSearchAuto(vars["word"])
	log.Printf("roots: %s", roots)

	outputRootSet := transformRoots(roots, transformer)
	if m, err := marshalRoots(outputRootSet); err == nil {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, "", m)
	}
}

// JSONExactAbjad searches words with `number` as their abjad counterpart
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
func JSONExactAbjad(w http.ResponseWriter, r *http.Request) {

	transformer := func(root *dervaze.Root) *dervaze.Root {
		r := dervaze.Root{
			TurkishLatin: root.TurkishLatin,
			Abjad:        root.Abjad,
			Ottoman: &dervaze.OttomanWord{
				Unicode: root.Ottoman.Unicode,
			},
		}
		return &r
	}
	vars := mux.Vars(r)
	log.Printf("JsonExactTr Vars: %s", vars)
	val, err := strconv.Atoi(vars["word"])
	var roots []*dervaze.Root
	if err == nil {
		roots = dervaze.IndexSearchAbjad(int32(val))
	} else {
		log.Printf("Error in SearchAbjad parameter: %s", vars["word"])
		roots = make([]*dervaze.Root, 0)
	}
	log.Printf("roots: %s", roots)

	outputRootSet := transformRoots(roots, transformer)
	if m, err := marshalRoots(outputRootSet); err == nil {

		w.Header().Add("Access-Control-Allow-Origin", "*")

		fmt.Fprintln(w, "", m)
	}
}

// JSONCalcAbjad calculates the abjad of a word
// ## `/v1/json/calc/abjad/{word}
//
// Calculates abjad for the `word` given in unicode
//
// ```
// { "ottoman_unicode": "word",
// "abjad": 1234 }
// ```
//
func JSONCalcAbjad(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	log.Printf("JsonPrefixTr Vars: %s", vars)
	abjad := dervaze.UnicodeToAbjad(vars["word"])
	str := fmt.Sprintf("{ \"ottoman_unicode\": \"%s\", \"abjad\": %d }", vars["word"], abjad)
	fmt.Fprintln(w, "", str)

}

// JSONV2U converts a visenc string to unicode
// ## `/v1/json/v2u/{word}
//
// Converts `word` from visenc to unicode
//
// ```
// {"ottoman_visenc": <word>,
// "ottoman_unicode": <unicode>}
// ```
//
func JSONV2U(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	log.Printf("JsonV2U Vars: %s", vars)
	unicode := dervaze.VisencToUnicode(vars["word"])
	str := fmt.Sprintf("{ \"ottoman_unicode\": \"%s\", \"ottoman_unicode\": \"%s\" }", unicode, vars["word"])
	fmt.Fprintln(w, "", str)
}

// JSONU2V converts a unicode string to visenc
// ## `/v1/json/u2v/{word}
//
// Converts `word` from unicode to visenc
//
// ```
// {"ottoman_visenc": <visenc>,
// "ottoman_unicode": <word>}
// ```
//
func JSONU2V(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	log.Printf("JsonU2V Vars: %s", vars)
	visenc := dervaze.UnicodeToVisenc(vars["word"])
	str := fmt.Sprintf("{ \"ottoman_unicode\": \"%s\", \"ottoman_unicode\": \"%s\" }", vars["word"], visenc)
	fmt.Fprintln(w, "", str)
}

// JSONVersion sends git version information
func JSONVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	out, err := exec.Command("git", "log", "--summary", "--max-count", "1").Output()
	if err == nil {
		head := strings.Split(strings(out), "\n")
		head = head[:3]
		fmt.Fprintln(w, "", strings.Join(head, "\n"))
	} else {
		fmt.Fprintln(w, "", string(err.Error()))
	}
}
func server(host, port string) {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/v1/json/prefix/tr/{word}", JSONPrefixTr)
	router.HandleFunc("/v1/json/prefix/ot/{word}", JSONPrefixOt)
	router.HandleFunc("/v1/json/exact/tr/{word}", JSONExactTr)
	router.HandleFunc("/v1/json/exact/ot/{word}", JSONExactOt)
	router.HandleFunc("/v1/json/search/any/{word}", JSONSearchAuto)
	router.HandleFunc("/v1/json/search/ot/{word}", JSONSearchOt)
	router.HandleFunc("/v1/json/search/tr/{word}", JSONSearchTr)
	router.HandleFunc("/v1/json/exact/abjad/{number}", JSONExactAbjad)
	router.HandleFunc("/v1/json/calc/abjad/{word}", JSONCalcAbjad)
	router.HandleFunc("/v1/json/v2u/{word}", JSONV2U)
	router.HandleFunc("/v1/json/u2v/{word}", JSONU2V)
	router.HandleFunc("/v1/version/", JSONVersion)
	srv := &http.Server{
		Handler:      router,
		Addr:         host + ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

func main() {

	var inputfile string
	var port string
	var host string

	flag.StringVar(&inputfile, "i", "../../assets/dervaze-rootset.protobuf", "protobuffer file to load roots")
	flag.StringVar(&host, "h", "127.0.0.1", "IP address or hostname to listen to")
	flag.StringVar(&port, "p", "9876", "port to listen to")

	flag.Parse()

	fmt.Println("Starting Server")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("%s=%s [%s] \n", f.Name, f.Value.String(), f.Usage)
	})

	dervaze.InitSearch(inputfile)
	server(host, port)
}

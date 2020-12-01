package lang

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/tchap/go-patricia/patricia"
)

var rootSet *RootSet
var turkishLatinTrie *patricia.Trie
var visencTrie *patricia.Trie
var unicodeTrie *patricia.Trie

var turkishLatinIndex *map[rune][]string
var visencIndex *map[rune][]string
var unicodeIndex *map[rune][]string

var abjadIndex *map[int32][]int

func buildTrie(roots []*Root, keyfunc func(*Root, int) string) *patricia.Trie {
	trie := patricia.NewTrie()

	for i, r := range roots {
		trie.Insert(patricia.Prefix(keyfunc(r, i)), i)
	}

	return trie

}

func buildIndex(roots []*Root, keyfunc func(*Root, int) string) *map[rune][]string {
	m := make(map[rune][]string)

	for i, r := range roots {
		s := keyfunc(r, i)
		r := []rune(s)
		var mapkey rune

		if len(r) > 0 {
			mapkey = r[0]
		} else {
			continue
		}

		_, exists := m[mapkey]
		if exists == false {
			m[mapkey] = make([]string, 0)
		}
		m[mapkey] = append(m[mapkey], s)
	}
	return &m
}

func buildAbjadIndex(roots []*Root) *map[int32][]int {

	m := make(map[int32][]int)

	for i, r := range roots {
		abj := r.Abjad
		_, exists := m[abj]
		if !exists {
			m[abj] = make([]int, 0)
		}
		m[abj] = append(m[abj], i)
	}

	return &m
}

// filterResults filters the identical elements by checking TurkishLatin and Unicode
func filterResults(roots []*Root) []*Root {
	theSet := map[string]*Root{}
	outList := make([]*Root, 0, len(roots))

	for _, r := range roots {
		k := r.TurkishLatin + r.Ottoman.Unicode
		_, exists := theSet[k]
		if !exists {
			theSet[k] = r
			outList = append(outList, r)
		}
	}
	return outList
}

// sorts roots by length of TurkishLatin
func sortByLength(roots []*Root) []*Root {
	sort.Slice(roots, func(i, j int) bool {
		return len(roots[i].TurkishLatin) < len(roots[j].TurkishLatin)
	})
	return roots
}

// InitSearch loads protobuf file and builds Trie and []string indices for turkishLatin, visenc and unicode
func InitSearch(protobuffile string) {
	rootSet = LoadRootSetProtobuf(protobuffile)

	turkishLatinTrie = buildTrie(rootSet.Roots, func(r *Root, i int) string { return fmt.Sprintf("%s#%d", r.TurkishLatin, i) })
	visencTrie = buildTrie(rootSet.Roots, func(r *Root, i int) string { return fmt.Sprintf("%s#%d", r.Ottoman.Visenc, i) })
	unicodeTrie = buildTrie(rootSet.Roots, func(r *Root, i int) string { return fmt.Sprintf("%s#%d", r.Ottoman.Unicode, i) })

	turkishLatinIndex = buildIndex(rootSet.Roots, func(r *Root, i int) string { return fmt.Sprintf("%s#%d", r.TurkishLatin, i) })
	visencIndex = buildIndex(rootSet.Roots, func(r *Root, i int) string { return fmt.Sprintf("%s#%d", r.Ottoman.Visenc, i) })
	unicodeIndex = buildIndex(rootSet.Roots, func(r *Root, i int) string { return fmt.Sprintf("%s#%d", r.Ottoman.Unicode, i) })

	abjadIndex = buildAbjadIndex(rootSet.Roots)
}

func GetRootSet() *RootSet {
	return rootSet
}

func GetTurkishLatinTrie() *patricia.Trie {
	return turkishLatinTrie
}

func GetVisencTrie() *patricia.Trie {
	return visencTrie
}

func GetUnicodeTrie() *patricia.Trie {
	return unicodeTrie
}

func GetTurkishLatinIndex() *map[rune][]string {
	return turkishLatinIndex
}

func GetVisencIndex() *map[rune][]string {
	return visencIndex
}

func GetUnicodeIndex() *map[rune][]string {
	return unicodeIndex
}

func GetAbjadIndex() *map[int32][]int {
	return abjadIndex
}

func PrefixSearchTurkishLatin(turkishLatin string, maxLen int) []*Root {
	results := make([]*Root, 0)
	visitFunc := func(_ patricia.Prefix, item patricia.Item) error {
		i, ok := item.(int)
		if ok {
			results = append(results, rootSet.Roots[i])
		} else {
			log.Printf("Error for %s in SearchTurkishLatin", item)
			return errors.New("item error")
		}
		return nil

	}
	turkishLatinTrie.VisitSubtree(patricia.Prefix(turkishLatin), visitFunc)

	results = filterResults(results)
	results = sortByLength(results)
	if maxLen > len(results) {
		results = results[:maxLen]
	}

	return results
}

// PrefixSearchTurkishLatinExact returns a single Root where Root.TurkishLatin == turkishLatin
func PrefixSearchTurkishLatinExact(turkishLatin string) []*Root {
	return PrefixSearchTurkishLatin(turkishLatin+"#", 1)
}

func PrefixSearchVisenc(visenc string, maxLen int) []*Root {
	results := make([]*Root, 0)
	visitFunc := func(_ patricia.Prefix, item patricia.Item) error {
		i, ok := item.(int)
		if ok {
			results = append(results, rootSet.Roots[i])
		} else {
			log.Printf("Error for %s in SearchVisenc", item)
			return errors.New("item error")
		}
		return nil
	}
	visencTrie.VisitSubtree(patricia.Prefix(visenc), visitFunc)

	results = filterResults(results)
	results = sortByLength(results)
	if maxLen > len(results) {
		results = results[:maxLen]
	}

	return results
}

func PrefixSearchVisencExact(visenc string) []*Root {
	return PrefixSearchVisenc(visenc+"#", 1)
}

func PrefixSearchUnicode(unicode string, maxLen int) []*Root {
	results := make([]*Root, 0)
	visitFunc := func(_ patricia.Prefix, item patricia.Item) error {
		i, ok := item.(int)
		if ok {
			results = append(results, rootSet.Roots[i])
		} else {
			log.Printf("Error for %s in SearchUnicode", item)
			return errors.New("item error")
		}
		return nil
	}
	unicodeTrie.VisitSubtree(patricia.Prefix(unicode), visitFunc)

	results = filterResults(results)
	results = sortByLength(results)
	if maxLen > len(results) {
		results = results[:maxLen]
	}

	return results
}

func PrefixSearchUnicodeExact(unicode string) []*Root {
	return PrefixSearchUnicode(unicode+"#", 1)
}

func PrefixSearchAll(term string, maxLen int) []*Root {
	results := make([]*Root, 0)
	val, err := strconv.Atoi(term)
	if err == nil {
		results = append(results, IndexSearchAbjad(int32(val), maxLen)...)
	}

	results = append(results, PrefixSearchTurkishLatin(term, maxLen)...)
	results = append(results, PrefixSearchUnicode(term, maxLen)...)
	results = append(results, PrefixSearchVisenc(term, maxLen)...)

	results = filterResults(results)
	results = sortByLength(results)
	if maxLen > len(results) {
		results = results[:maxLen]
	}

	return results
}

func splitRootKeyFromIndex(k string) (int, error) {
	elements := strings.Split(k, "#")
	if len(elements) != 2 {
		return -1, errors.New(fmt.Sprintf("Malformed Index String %s", k))
	}

	i, err := strconv.Atoi(elements[1])
	if err != nil {
		return -1, err
	}

	return i, nil
}

func indexRegexSearch(keylist []string, r *regexp.Regexp) []*Root {

	roots := make([]*Root, 0)

	for _, k := range keylist {
		if r.MatchString(k) {
			ri, err := splitRootKeyFromIndex(k)
			if err == nil {
				roots = append(roots, rootSet.Roots[ri])
			} else {
				log.Println(err)
			}
		}
	}
	return roots
}

func RegexSearchTurkishLatin(word string, maxLen int) []*Root {

	runes := []rune(word)
	var sb strings.Builder
	sb.WriteString(".?")
	for _, r := range runes {
		sb.WriteRune(r)
		sb.WriteString(".?")
	}

	searchRegex := regexp.MustCompile(sb.String())

	results := make([]*Root, 0)

	// TODO convert the search to goroutines

	for _, indexList := range *turkishLatinIndex {
		docResults := indexRegexSearch(indexList, searchRegex)
		results = append(results, docResults...)
	}

	results = filterResults(results)
	results = sortByLength(results)
	if maxLen > len(results) {
		results = results[:maxLen]
	}

	return results
}

func RegexSearchUnicode(word string, maxLen int) []*Root {

	runes := []rune(word)
	var sb strings.Builder
	sb.WriteString(".?")
	for _, r := range runes {
		sb.WriteRune(r)
		sb.WriteString(".?")
	}

	searchRegex := regexp.MustCompile(sb.String())

	results := make([]*Root, 0)

	// TODO convert the search to goroutines

	for _, indexList := range *unicodeIndex {
		docResults := indexRegexSearch(indexList, searchRegex)
		results = append(results, docResults...)
	}

	results = filterResults(results)
	results = sortByLength(results)
	if maxLen > len(results) {
		results = results[:maxLen]
	}

	return results
}

func RegexSearchVisenc(word string, maxLen int) []*Root {

	visencLetters := SplitVisenc(word, false)

	var sb strings.Builder
	sb.WriteString(".?")
	for _, v := range visencLetters {
		sb.WriteString(v)
		sb.WriteString(".?")
	}

	searchRegex := regexp.MustCompile(sb.String())

	results := make([]*Root, 0)

	// TODO convert the search to goroutines

	for _, indexList := range *visencIndex {
		docResults := indexRegexSearch(indexList, searchRegex)
		results = append(results, docResults...)
	}

	results = filterResults(results)
	results = sortByLength(results)
	if maxLen > len(results) {
		results = results[:maxLen]
	}

	return results
}

func RegexSearchAuto(word string, maxLen int) []*Root {

	if ContainsArabicChars(word) {
		return RegexSearchUnicode(word, maxLen)
	} else if ContainsDigits(word) {
		if val, err := strconv.Atoi(word); err == nil {
			return IndexSearchAbjad(int32(val), maxLen)
		} else {
			return RegexSearchVisenc(word, maxLen)
		}
	} else {
		return RegexSearchTurkishLatin(word, maxLen)
	}
}

func IndexSearchAbjad(abjad int32, maxLen int) []*Root {

	indices, exists := (*abjadIndex)[abjad]

	if !exists {
		indices = make([]int, 0, 0)
	}

	roots := make([]*Root, len(indices))

	for i, v := range indices {
		roots[i] = rootSet.Roots[v]
	}

	roots = filterResults(roots)
	roots = sortByLength(roots)
	if maxLen > len(roots) {
		roots = roots[:maxLen]
	}

	return roots
}

func PrintRoots(roots []*Root) string {
	out := ""
	for i, r := range roots {
		out += fmt.Sprintf("%d - %s | %s | %s | %d\n", i, r.TurkishLatin, r.Ottoman.Unicode, r.Ottoman.Visenc, r.Ottoman.Abjad)
	}
	return out
}

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

// GetRootSet returns the rootSet whole package uses
func GetRootSet() *RootSet {
	return rootSet
}

// GetTurkishLatinTrie returns a trie keeping turkishLatin roots
func GetTurkishLatinTrie() *patricia.Trie {
	return turkishLatinTrie
}

// GetVisencTrie returns a trie keeping visenc of roots
func GetVisencTrie() *patricia.Trie {
	return visencTrie
}

// GetUnicodeTrie returns a trie for unicode roots
func GetUnicodeTrie() *patricia.Trie {
	return unicodeTrie
}

// GetTurkishLatinIndex returns turkishLatinIndex
func GetTurkishLatinIndex() *map[rune][]string {
	return turkishLatinIndex
}

// GetVisencIndex returns visencIndex
func GetVisencIndex() *map[rune][]string {
	return visencIndex
}

// GetUnicodeIndex returns unicode index
func GetUnicodeIndex() *map[rune][]string {
	return unicodeIndex
}

// GetAbjadIndex returns index of all roots sharing common abjad value
func GetAbjadIndex() *map[int32][]int {
	return abjadIndex
}

// PrefixSearchTurkishLatin returns list of roots whose TurkishLatin begins with `turkishLatin`
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
	if maxLen < len(results) {
		results = results[:maxLen]
	}

	return results
}

// PrefixSearchTurkishLatinExact returns a single Root where Root.TurkishLatin == turkishLatin
func PrefixSearchTurkishLatinExact(turkishLatin string) []*Root {
	return PrefixSearchTurkishLatin(turkishLatin+"#", 1)
}

// PrefixSearchVisenc returns list of roots whose Visenc starts with `visenc`
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
	if maxLen < len(results) {
		results = results[:maxLen]
	}

	return results
}

// PrefixSearchVisencExact returns a maximum of 10 Root having Visenc = `visenc`
func PrefixSearchVisencExact(visenc string) []*Root {
	return PrefixSearchVisenc(visenc+"#", 10)
}

// PrefixSearchUnicode searches roots by unicode string
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
	if maxLen < len(results) {
		results = results[:maxLen]
	}

	return results
}

//PrefixSearchUnicodeExact returns maximum 10 roots with having a prefix unicode
func PrefixSearchUnicodeExact(unicode string) []*Root {
	return PrefixSearchUnicode(unicode+"#", 10)
}

// PrefixSearchAll runs PrefixSearchTurkishLatin, PrefixSearchUnicode, PrefixSearchVisenc, IndexSearchAbjad and combines results.
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
	if maxLen < len(results) {
		results = results[:maxLen]
	}

	return results
}

func splitRootKeyFromIndex(k string) (int, error) {
	elements := strings.Split(k, "#")
	if len(elements) != 2 {

		return -1, fmt.Errorf("Malformed Index String %s", k)
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

// FuzzySearchTurkishLatin searches word in the string index via regexes.
// `word` is searched as `.?w.?o.?r.?d.?`
func FuzzySearchTurkishLatin(word string, maxLen int) []*Root {

	runes := []rune(word)
	var sb strings.Builder
	sb.WriteString(".*")
	for _, r := range runes {
		sb.WriteRune(r)
		sb.WriteString(".*")
	}

	searchRegex := regexp.MustCompile(sb.String())
	return RegexSearchTurkishLatin(searchRegex, maxLen)
}

// RegexSearchTurkishLatin searches turkishLatinIndex with the supplied regex
func RegexSearchTurkishLatin(regex *regexp.Regexp, maxLen int) []*Root {

	results := make([]*Root, 0)

	// TODO convert the search to goroutines

	for _, indexList := range *turkishLatinIndex {
		docResults := indexRegexSearch(indexList, regex)
		results = append(results, docResults...)
	}

	results = filterResults(results)
	results = sortByLength(results)
	if maxLen < len(results) {
		results = results[:maxLen]
	}

	return results
}

// FuzzySearchUnicode searches `word` in unicode indices
func FuzzySearchUnicode(word string, maxLen int) []*Root {

	runes := []rune(word)
	var sb strings.Builder
	sb.WriteString(".*")
	for _, r := range runes {
		sb.WriteRune(r)
		sb.WriteString(".*")
	}

	searchRegex := regexp.MustCompile(sb.String())
	return RegexSearchUnicode(searchRegex, maxLen)
}

// RegexSearchUnicode searches unicodeIndex with the supplied regex and returns at most maxLen results
func RegexSearchUnicode(regex *regexp.Regexp, maxLen int) []*Root {

	results := make([]*Root, 0)

	// TODO convert the search to goroutines

	for _, indexList := range *unicodeIndex {
		docResults := indexRegexSearch(indexList, regex)
		results = append(results, docResults...)
	}

	results = filterResults(results)
	results = sortByLength(results)
	if maxLen < len(results) {
		results = results[:maxLen]
	}

	return results
}

// FuzzySearchVisenc searches word in visencIndices using fuzzy matching
func FuzzySearchVisenc(word string, maxLen int) []*Root {

	visencLetters := SplitVisenc(word, false)

	var sb strings.Builder
	sb.WriteString(".*")
	for _, v := range visencLetters {
		sb.WriteString(v)
		sb.WriteString(".*")
	}

	searchRegex := regexp.MustCompile(sb.String())
	return RegexSearchVisenc(searchRegex, maxLen)

}

// RegexSearchVisenc makes a search in visenc field with the supplied regexp
func RegexSearchVisenc(regex *regexp.Regexp, maxLen int) []*Root {

	results := make([]*Root, 0)

	// TODO convert the search to goroutines

	for _, indexList := range *visencIndex {
		docResults := indexRegexSearch(indexList, regex)
		results = append(results, docResults...)
	}

	results = filterResults(results)
	results = sortByLength(results)
	if maxLen < len(results) {
		results = results[:maxLen]
	}

	return results
}

// FuzzySearchAuto searches word in either of FuzzySearchUnicode, FuzzySearchTurkishLatin, FuzzySearchVisenc and IndexSearchAbjad
func FuzzySearchAuto(word string, maxLen int) []*Root {

	if ContainsArabicChars(word) {
		return FuzzySearchUnicode(word, maxLen)
	} else if ContainsDigits(word) {
		if val, err := strconv.Atoi(word); err == nil {
			return IndexSearchAbjad(int32(val), maxLen)
		}
		return FuzzySearchVisenc(word, maxLen)
	}
	return FuzzySearchTurkishLatin(word, maxLen)
}

// RegexSearchAuto searches word in either of RegexSearchUnicode, RegexSearchTurkishLatin, RegexSearchVisenc and IndexSearchAbjad
func RegexSearchAuto(regexp *regexp.Regexp, maxLen int) []*Root {

	word := regexp.String()

	if ContainsArabicChars(word) {
		return RegexSearchUnicode(regexp, maxLen)
	} else if ContainsDigits(word) {
		if val, err := strconv.Atoi(word); err == nil {
			return IndexSearchAbjad(int32(val), maxLen)
		}
		return RegexSearchVisenc(regexp, maxLen)
	}
	return RegexSearchTurkishLatin(regexp, maxLen)
}

// IndexSearchAbjad searches returns list roots containing `abjad` as value
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
	if maxLen < len(roots) {
		roots = roots[:maxLen]
	}

	return roots
}

// PrintRoots returns roots' TurkishLatin, Unicode, Visenc and Abjad as a single string
func PrintRoots(roots []*Root) string {
	out := ""
	for i, r := range roots {
		out += fmt.Sprintf("%d - %s | %s | %s | %d\n", i, r.TurkishLatin, r.Ottoman.Unicode, r.Ottoman.Visenc, r.Ottoman.Abjad)
	}
	return out
}

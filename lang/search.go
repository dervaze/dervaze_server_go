package lang

import (
	"errors"
	"fmt"
	"github.com/tchap/go-patricia/patricia"
	"log"
	"strconv"
)

var rootSet *RootSet
var turkishLatinTrie *patricia.Trie
var visencTrie *patricia.Trie
var unicodeTrie *patricia.Trie

func buildTrie(roots []*Root, keyfunc func(*Root, int) string) *patricia.Trie {
	trie := patricia.NewTrie()

	for i, r := range roots {
		trie.Insert(patricia.Prefix(keyfunc(r, i)), i)
	}

	return trie

}

func buildTries(rootset *RootSet) (*patricia.Trie, *patricia.Trie, *patricia.Trie) {
	turkishLatinTrie := buildTrie(rootset.Roots, func(r *Root, i int) string { return r.TurkishLatin + "#" + string(i) })
	visencTrie := buildTrie(rootset.Roots, func(r *Root, i int) string { return r.Ottoman.Visenc + "#" + string(i) })
	unicodeTrie := buildTrie(rootset.Roots, func(r *Root, i int) string { return r.Ottoman.Unicode + "#" + string(i) })
	return turkishLatinTrie, visencTrie, unicodeTrie
}

func InitSearch(protobuffile string) {
	RootSet = LoadRootSetProtobuf(protobuffile)
	turkishLatinTrie, visencTrie, unicodeTrie = buildTries(rootSet)
}

func GetRootSet() *RootSet {
	return rootSet
}

func SearchTurkishLatin(turkishLatin string) []*Root {
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
	return results
}

func SearchTurkishLatinExact(turkishLatin string) []*Root {
	return SearchTurkishLatin(turkishLatin + "#")
}

func SearchVisenc(visenc string) []*Root {
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

	return results
}

func SearchVisencExact(visenc string) []*Root {
	return SearchVisenc(visenc + "#")
}

func SearchUnicode(unicode string) []*Root {
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

	return results
}

func SearchUnicodeExact(unicode string) []*Root {
	return SearchUnicode(unicode + "#")
}

func SearchAbjad(abjad int32) []*Root {
	results := make([]*Root, 0)

	for _, r := range rootSet.Roots {
		if r.Ottoman.Abjad == abjad {
			results = append(results, r)
		}
	}
	return results
}

func SearchAll(term string) []*Root {
	results := make([]*Root, 0)
	val, err := strconv.Atoi(term)
	if err == nil {
		results = append(results, SearchAbjad(int32(val))...)
	}

	results = append(results, SearchTurkishLatin(term)...)
	results = append(results, SearchUnicode(term)...)
	results = append(results, SearchVisenc(term)...)

	return results
}

func PrintRoots(roots []*Root) string {
	out := ""
	for i, r := range roots {
		out += fmt.Sprintf("%d - %s | %s | %s | %d\n", i, r.TurkishLatin, r.Ottoman.Unicode, r.Ottoman.Visenc, r.Ottoman.Abjad)
	}
	return out
}

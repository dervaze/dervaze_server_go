package lang

import (
	"github.com/tchap/go-patricia/patricia"
)

var rootSet *RootSet
var turkishLatinTrie *patricia.Trie
var visencTrie *patricia.Trie
var unicodeTrie *patricia.Trie

func buildTrie(roots []*Root, keyfunc func(*Root) string) *patricia.Trie {
	trie := patricia.NewTrie()

	for i, r := range roots {
		trie.Insert(patricia.Prefix(keyfunc(r)), i)
	}

	return trie

}

func buildTries(rootset *RootSet) (*patricia.Trie, *patricia.Trie, *patricia.Trie) {
	turkishLatinTrie := buildTrie(rootset.Roots, func(r *Root) string { return r.TurkishLatin })
	visencTrie := buildTrie(rootset.Roots, func(r *Root) string { return r.Ottoman.Visenc })
	unicodeTrie := buildTrie(rootset.Roots, func(r *Root) string { return r.Ottoman.Unicode })
	return turkishLatinTrie, visencTrie, unicodeTrie
}

func InitSearch(protobuffile string) {
	rootSet = LoadRootSetProtobuf(protobuffile)
	turkishLatinTrie, visencTrie, unicodeTrie = buildTries(rootSet)
}

func SearchTurkishLatin(turkishLatin string) []*Root {
	results := make([]*Root, 0)
	visitFunc := func(_ patricia.Prefix, item patricia.Item) error {
		results = append(results, rootSet.Roots[item])
	}
	turkishLatinTrie.VisitPrefixes(turkishLatin, visitFunc)

	return results

}

func SearchVisenc(visenc string) []*Root {
	results := make([]*Root, 0)
	visitFunc := func(_ patricia.Prefix, item patricia.Item) error {
		results = append(results, rootSet.Roots[item])
	}
	visencTrie.VisitPrefixes(visenc, visitFunc)

	return results
}

func SearchUnicode(unicode string) []*Root {
	results := make([]*Root, 0)
	visitFunc := func(_ patricia.Prefix, item patricia.Item) error {
		results = append(results, rootSet.Roots[item])
	}
	unicodeTrie.VisitPrefixes(unicode, visitFunc)

	return results
}

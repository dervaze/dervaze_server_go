package lang

import (
	"github.com/tchap/go-patricia/patricia"
)

func BuildTrie(roots []*Root, keyfunc func(*Root) string) *patricia.Trie {
	trie := patricia.NewTrie()

	for i, r := range roots {
		trie.Insert(patricia.Prefix(keyfunc(r)), i)
	}

	return trie

}

func BuildTries(rootset RootSet) (*patricia.Trie, *patricia.Trie, *patricia.Trie) {
	turkishLatinTrie := BuildTrie(rootset.Roots, func(r *Root) string { return r.TurkishLatin })
	visencTrie := BuildTrie(rootset.Roots, func(r *Root) string { return r.Ottoman.Visenc })
	unicodeTrie := BuildTrie(rootset.Roots, func(r *Root) string { return r.Ottoman.Unicode })
	return turkishLatinTrie, visencTrie, unicodeTrie
}

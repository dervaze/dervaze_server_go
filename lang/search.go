package lang

import (
	"github.com/tchap/go-patricia/patricia"
	"io/ioutil"
)

func BuildTrie(roots []Root, keyfunc func(Root) string) Trie {
	trie := NewTrie()

	for i, r := range roots {
		trie.Insert(Prefix(keyfunc(r)), i)
	}

	return trie

}

func BuildTries(roots []Root) (Trie, Trie, Trie) {
	turkishLatinTrie := BuildTrie(roots, GetTurkishLatin)
	visencTrie := BuildTrie(roots, func(r Root) { return r.Ottoman.Visenc })
	unicodeTrie := BuildTrie(roots, func(r Root) { return r.Ottoman.Unicode })
	return turkishLatinTrie, visencTrie, unicodeTrie
}

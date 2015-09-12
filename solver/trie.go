package solver

import (
	"strings"
)

type Trie struct {
	isWord   bool
	children [26]*Trie
}

func NewTrie() *Trie {
	return new(Trie)
}

// Insert stores a word in the Trie.
func (t *Trie) Insert(word string) {
	var index int
	var next *Trie
	for _, c := range strings.ToLower(word) {
		index = int(c - 'a')
		next = t.children[index]
		if next == nil {
			next = new(Trie)
			t.children[index] = next
		}
		t = next
	}
	t.isWord = true
}

// Contains checks to see if the Trie contains the given word.
func (t *Trie) ContainsWord(word string) bool {
	for _, c := range strings.ToLower(word) {
		t = t.children[int(c-'a')]
		if t == nil {
			return false
		}
	}
	return t.isWord
}

// ContainsChar checks to see if the Trie contains the given letter.
func (t *Trie) ContainsLetter(c rune) bool {
	return t.Child(c) != nil
}

// GetChild returns the sub-Trie for the given character.
func (t *Trie) Child(c rune) *Trie {
	return t.children[int(c-'a')]
}

// IsWord returns true if the Trie completes a word.
func (t *Trie) IsWord() bool {
	return t.isWord
}

package solver

import (
	"strings"
)

type Trie struct {
	isWord   bool
	children map[rune]*Trie
}

func NewTrie() *Trie {
	return new(Trie)
}

// Insert stores a word in the Trie.
func (t *Trie) Insert(word string) {
	var next *Trie
	for _, c := range strings.ToLower(word) {
		next = t.children[c]
		if next == nil {
			next = new(Trie)
			if t.children == nil {
				t.children = map[rune]*Trie{}
			}
			t.children[c] = next
		}
		t = next
	}
	t.isWord = true
}

// Contains checks to see if the Trie contains the given word.
func (t *Trie) ContainsWord(word string) bool {
	for _, c := range strings.ToLower(word) {
		t = t.children[c]
		if t == nil {
			return false
		}
	}
	return t.isWord
}

// ContainsChar checks to see if the Trie node contains the given letter.
func (t *Trie) ContainsLetter(c rune) bool {
	return t.Child(c) != nil
}

// Child returns the sub-Trie for the given character.
func (t *Trie) Child(c rune) *Trie {
	return t.children[c]
}

// IsWord returns true if the given Trie node completes a word.
func (t *Trie) IsWord() bool {
	return t.isWord
}

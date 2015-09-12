package solver

import (
	"testing"
)

func TestInsert(t *testing.T) {
	root := NewTrie()

	root.Insert("hello")
	root.Insert("world")
	root.Insert("worldwide")

	if !root.ContainsWord("hello") {
		t.Error("missing word")
	}

	if !root.ContainsWord("world") {
		t.Error("missing word")
	}

	if !root.ContainsWord("worldwide") {
		t.Error("missing word")
	}

	if root.ContainsWord("foo") {
		t.Error("invalid result from Contains()")
	}
}

func TestPath(t *testing.T) {
	root := NewTrie()

	root.Insert("hi")

	if root.ContainsLetter('a') {
		t.Error("Trie has invalid letter 'a'")
	}
	if !root.ContainsLetter('h') {
		t.Error("Trie missing letter 'h'")
	}

	if root.Child('a') != nil {
		t.Error("invalid child")
	}

	child := root.Child('h')
	if child == nil {
		t.Error("missing child")
	}
	if child.IsWord() {
		t.Error("invalid word end")
	}

	child = child.Child('i')
	if child == nil {
		t.Error("missing child")
	}
	if !child.IsWord() {
		t.Error("expected word missing")
	}
}

package data

import (
	"testing"
)

func TestInsertSimple(t *testing.T) {
	trie := NewTrie()
	trie.Insert("one", "foo")

	v, ok := trie.Get("foo")
	if ok != true || v != "one" {
		t.Errorf("Incorrect value")
	}
}

func TestInsertMiss(t *testing.T) {
	trie := NewTrie()
	trie.Insert("one", "foo")

	_, ok := trie.Get("bar")
	if ok != false {
		t.Errorf("Incorrect value")
	}
}

func TestInsertPath(t *testing.T) {
	trie := NewTrie()
	trie.Insert("one", "foo", "bar", "baz")

	v, ok := trie.Get("foo", "bar", "baz")
	if ok != true || v != "one" {
		t.Errorf("Incorrect value")
	}
}

func TestInsertPathMiss(t *testing.T) {
	trie := NewTrie()
	trie.Insert("one", "foo", "bar", "baz")

	_, ok := trie.Get("foo", "bar", "barry")
	if ok != false {
		t.Errorf("Incorrect value")
	}
}

func TestMergeSimple(t *testing.T) {
	trie := NewTrie()
	trie.Insert("one", "foo")
	trie.Insert("two", "bar")

	v, ok := trie.Get("foo")
	if ok != true || v != "one" {
		t.Errorf("Incorrect value")
	}

	v, ok = trie.Get("bar")
	if ok != true || v != "two" {
		t.Errorf("Incorrect value")
	}
}

func TestMergePaths(t *testing.T) {
	trie := NewTrie()
	trie.Insert("one", "foo", "bar", "baz")
	trie.Insert("two", "foo", "bar", "barry")
	trie.Insert("three", "foo", "lol", "rofl")
	trie.Insert("zero", "foo")

	v, ok := trie.Get("foo")
	if ok != true || v != "zero" {
		t.Errorf("Incorrect value")
	}

	v, ok = trie.Get("foo", "bar", "baz")
	if ok != true || v != "one" {
		t.Errorf("Incorrect value")
	}

	v, ok = trie.Get("foo", "bar", "barry")
	if ok != true || v != "two" {
		t.Errorf("Incorrect value")
	}

	v, ok = trie.Get("foo", "lol", "rofl")
	if ok != true || v != "three" {
		t.Errorf("Incorrect value")
	}
}

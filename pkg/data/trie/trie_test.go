package trie

import (
	"testing"

	"github.com/mt-inside/go-usvc"
	"github.com/stretchr/testify/assert"
)

func TestInsertSimple(t *testing.T) {
	trie := NewTrie(usvc.GetLogger(true).WithName("trie"))
	trie.Insert(Some("one"), "foo")

	v, ok := trie.Get("foo")
	assert.Equal(t, ok, true, "Couldn't find")
	assert.Equal(t, v, Some("one"), "Incorrect value")
	assert.Equal(t, v.Render(), "one", "Incorrect rendering")
}

func TestInsertMiss(t *testing.T) {
	trie := NewTrie(usvc.GetLogger(true).WithName("trie"))
	trie.Insert(Some("one"), "foo")

	_, ok := trie.Get("bar")
	assert.Equal(t, ok, false, "Incorrectly found")
}

func TestInsertPath(t *testing.T) {
	trie := NewTrie(usvc.GetLogger(true).WithName("trie"))
	trie.Insert(Some("one"), "foo", "bar", "baz")

	v, ok := trie.Get("foo", "bar", "baz")
	assert.Equal(t, ok, true, "Couldn't find")
	assert.Equal(t, v, Some("one"), "Incorrect value")
	assert.Equal(t, v.Render(), "one", "Incorrect rendering")
}

func TestInsertPathMiss(t *testing.T) {
	trie := NewTrie(usvc.GetLogger(true).WithName("trie"))
	trie.Insert(Some("one"), "foo", "bar", "baz")

	_, ok := trie.Get("foo", "bar", "barry")
	assert.Equal(t, ok, false, "Incorrectly found")
}

// func TestMergeSimple(t *testing.T) {
// 	t.Skip("Not implemented")

// 	trie := NewTrie()
// 	trie.Insert(Some("one"), "foo")
// 	trie.Insert(Some("two"), "bar")

// 	v, ok := trie.Get("foo")
// 	if ok != true || v != Some("one") {
// 		t.Errorf("Incorrect value")
// 	}

// 	v, ok = trie.Get("bar")
// 	if ok != true || v != "two" {
// 		t.Errorf("Incorrect value")
// 	}
// }

// func TestMergePaths(t *testing.T) {
// 	t.Skip("Not implemented")

// 	trie := NewTrie()
// 	trie.Insert(Some("one"), "foo", "bar", "baz")
// 	trie.Insert(Some("two"), "foo", "bar", "barry")
// 	trie.Insert(Some("three"), "foo", "lol", "rofl")
// 	trie.Insert(Some("zero"), "foo")

// 	v, ok := trie.Get("foo")
// 	if ok != true || v != Some("zero") {
// 		t.Errorf("Incorrect value")
// 	}

// 	v, ok = trie.Get("foo", "bar", "baz")
// 	if ok != true || v != "one" {
// 		t.Errorf("Incorrect value")
// 	}

// 	v, ok = trie.Get("foo", "bar", "barry")
// 	if ok != true || v != "two" {
// 		t.Errorf("Incorrect value")
// 	}

// 	v, ok = trie.Get("foo", "lol", "rofl")
// 	if ok != true || v != "three" {
// 		t.Errorf("Incorrect value")
// 	}
// }

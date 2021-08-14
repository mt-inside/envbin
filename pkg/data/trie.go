package data

import (
	"fmt"
	"strings"
	"time"
)

type TrieValue interface {
	Render() string
}

type Some struct {
	Value string
}

func (s Some) Render() string {
	return s.Value
}

type None struct{}

func (n None) Render() string {
	panic("Don't render me")
}

type NotPresent struct{}

func (np NotPresent) Render() string {
	return "NotPresent"
}

type Error struct {
	err error
}

func (e Error) Render() string {
	return fmt.Sprintf("Error: %v", e.err)
}

type Timeout struct {
	d time.Duration
}

func (t Timeout) Render() string {
	return fmt.Sprintf("Timed Out (waited %v)", t.d)
}

type Forbidden struct{}

func (f Forbidden) Render() string {
	return "Forbidden"
}

type Trie struct {
	Value    TrieValue
	children map[string]*Trie
}

func NewTrie() *Trie {
	return &Trie{Value: None{}, children: make(map[string]*Trie)}
}

func (t *Trie) Insert(value TrieValue, path ...string) {
	n := t
	for _, name := range path {
		if child, ok := n.children[name]; ok {
			n = child
		} else {
			n.children[name] = NewTrie()
			n = n.children[name]
		}
	}
	n.Value = value
}

func (t *Trie) Get(path ...string) (TrieValue, bool) {
	n := t
	for _, name := range path {
		if child, ok := n.children[name]; ok {
			n = child
		} else {
			return Some{""}, false
		}
	}
	return n.Value, true
}

type Entry struct {
	Path  []string
	Value TrieValue
}

func (e Entry) String() string {
	return strings.Join(e.Path, "/") + ": " + e.Value.Render()
}

func (t *Trie) Collect() []Entry {
	return t.collectInt([]Entry{}, []string{})
}

func (t *Trie) collectInt(entries []Entry, path []string) []Entry {
	entries = append(entries, Entry{path, t.Value})
	for name, c := range t.children {
		entries = c.collectInt(entries, append(path, name))
	}
	return entries
}

func (t *Trie) Walk(cb func(entry Entry)) {
	t.walkInt(cb, []string{})
}
func (t *Trie) walkInt(cb func(entry Entry), path []string) {
	cb(Entry{path, t.Value})
	for name, c := range t.children {
		c.walkInt(cb, append(path, name))
	}
}

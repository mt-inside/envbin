package data

import "strings"

type Trie struct {
	children map[string]*Trie
	Value    *string // TODO should be a type field: Value (implies this field too), None, Timeout, Forbidden, OtherError
}

func NewTrie() *Trie {
	return newTrie()
}

func newTrie() *Trie {
	return &Trie{children: make(map[string]*Trie)}
}

func (t *Trie) Insert(value string, path ...string) {
	n := t
	for _, name := range path {
		if child, ok := n.children[name]; ok {
			n = child
		} else {
			n.children[name] = newTrie()
			n = n.children[name]
		}
	}
	n.Value = &value
}

func (t *Trie) Get(path ...string) (string, bool) {
	n := t
	for _, name := range path {
		if child, ok := n.children[name]; ok {
			n = child
		} else {
			return "", false
		}
	}
	return *n.Value, true
}

type Entry struct {
	Path  []string
	Value string
}

func (e Entry) String() string {
	return strings.Join(e.Path, "/") + ": " + e.Value
}

func (t *Trie) Walk() []Entry {
	return t.walkInt([]Entry{}, []string{})
}

func (t *Trie) walkInt(entries []Entry, path []string) []Entry {
	if t.Value != nil {
		entries = append(entries, Entry{path, *t.Value})
	}
	for name, c := range t.children {
		entries = c.walkInt(entries, append(path, name))
	}
	return entries
}

func (t *Trie) Render() string {
	return t.renderInt("", 0, "/")
}

func (t *Trie) renderInt(s string, depth int, name string) string {
	s += strings.Repeat("  ", depth) + name
	if t.Value != nil {
		s += ": " + *t.Value
	}
	s += "\n"

	for n, c := range t.children {
		s = c.renderInt(s, depth+1, n)
	}

	return s
}

func (t *Trie) Merge(t2 *Trie) (tOut *Trie) {
	tOut = NewTrie()

	// TODO

	return tOut
}

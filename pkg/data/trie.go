package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-logr/logr"
)

type Trie struct {
	log      logr.Logger
	leaf     bool // zero-value => non-leaf
	children map[string]*Trie
	value    Value
}

func NewTrie(log logr.Logger) *Trie {
	return &Trie{
		log: log,
	}
}
func (t *Trie) Insert(value Value, path ...string) {
	log := t.log.WithName("Insert")
	log.V(1).Info("called", "path", path, "value", value)

	if len(path) == 0 {
		log.V(1).Info("making leaf")
		t.leaf = true
		t.value = value
		return
	}

	if t.children == nil {
		t.children = map[string]*Trie{}
	}

	name := path[0]
	if _, ok := t.children[name]; !ok {
		t.children[name] = NewTrie(t.log)
	}
	t.children[name].Insert(value, path[1:]...)
}
func (t *Trie) Get(path ...string) (Value, bool) {
	log := t.log.WithName("Get")
	log.V(1).Info("called", "path", path, "leaf?", t.leaf)

	if len(path) == 0 {
		if !t.leaf {
			log.Info("trie fuckup")
			return nil, false
		} else {
			log.V(1).Info("leaf, Some", "value", t.value)
			return t.value, true
		}
	} else {
		if t.leaf {
			switch t.value.(type) {
			case Some:
				log.Info("trie fuckup")
				return nil, false
			default:
				log.V(1).Info("leaf, !Some", "value", t.value)
				return t.value, true
			}
		} else {
			if t.children == nil {
				log.Info("trie fuckup")
				return nil, false
			}

			if _, ok := t.children[path[0]]; !ok {
				log.Info("trie fuckup")
				return nil, false
			}

			log.V(1).Info("recursing")
			return t.children[path[0]].Get(path[1:]...)
		}
	}
}
func (t *Trie) MarshalJSON() ([]byte, error) {
	if t.leaf {
		return json.Marshal(t.value.Render())
	} else {
		return json.Marshal(t.children)
	}
}

type Value interface {
	Render() string
}

type Some struct {
	Value string `json:"value"`
}

func (s Some) Render() string {
	return s.Value
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

// TODO
// whole thing to pkg/trie
// this to file walk
// as func () Walk(Trie t) { switch type }

func (t *Trie) Walk(cb func(path []string, value Value)) {
	t.walkInternal(cb, []string{})
}
func (t *Trie) walkInternal(cb func(path []string, value Value), path []string) {
	log := t.log.WithName("walkInternal")
	log.V(1).Info("called", "path", path, "leaf?", t.leaf)

	if t.leaf {
		cb(path, t.value)
	} else {
		if t.children == nil {
			panic("Invalid trie")
		}

		cb(path, Some{""})

		log.V(1).Info("recursing")
		for name, c := range t.children {
			c.walkInternal(cb, append(path, name))
		}
	}
}

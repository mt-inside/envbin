package trie

import (
	"encoding/json"
	"errors"
	"strings"
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

func (t *Trie) InsertTree(child *Trie, path ...string) {
	log := t.log.WithName("InsertTree")
	log.V(1).Info("called", "path", path)

	if t.children == nil {
		t.children = map[string]*Trie{}
	}

	if len(path) == 1 {
		name := path[0]
		t.children[name] = child
		return
	}

	name := path[0]
	if _, ok := t.children[name]; !ok {
		t.children[name] = NewTrie(t.log)
	}

	t.children[name].InsertTree(child, path[1:]...)
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
			case some:
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

func (t *Trie) UnmarshalJSON(bs []byte) error {
	if bs[0] == '{' { // lol hack - TODO: unmarshal into an interface{} and type-assert
		return json.Unmarshal(bs, &t.children)
	} else {
		t.leaf = true
		if strings.HasPrefix(string(bs), "NotPresent") {
			t.value = notPresent{}
		} else if strings.HasPrefix(string(bs), "Forbidden") {
			t.value = forbidden{}
		} else if strings.HasPrefix(string(bs), "Error") {
			t.value = erro{errors.New("lost - TODO render with structure")}
		} else if strings.HasPrefix(string(bs), "Timed Out") {
			t.value = timeout{time.Second} // TODO fixme also
		} else {
			var val string
			err := json.Unmarshal(bs, &val)
			if err != nil {
				return err
			}
			t.value = Some(val)
		}
	}
	return nil
}

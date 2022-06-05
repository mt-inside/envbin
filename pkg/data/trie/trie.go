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
	log.V(2).Info("called", "path", path, "value", value)

	if len(path) == 0 {
		log.V(2).Info("making leaf")
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
	log.V(2).Info("called", "path", path)

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
	log.V(2).Info("called", "path", path, "leaf?", t.leaf)

	if len(path) == 0 { // Asking for a leaf
		if !t.leaf {
			log.Info("trie fuckup")
			return nil, false
		} else {
			log.V(2).Info("leaf, Some", "value", t.value)
			return t.value, true
		}
	} else { // Asking to resurse
		if t.leaf { // At a leaf
			switch t.value.(type) {
			case some:
				// Not ok if this is a non-error leaf (a node shouldn't have a value and children - TODO enforce at insert time)
				log.Info("trie fuckup")
				return nil, false
			default:
				// we found a sub-tree error on the way to your key
				log.V(2).Info("leaf, !Some", "value", t.value)
				return t.value, true
			}
		} else { // Not at a leaf
			if t.children == nil {
				// empty sub-tree; no value, no childen (TODO fix up this damn data type. The type itself shouldn't allow it (golang permitting), else enforce at insert time. Abstract the .leaf into a .Leaf() which should just be able to look at run-time type or children length?)
				log.Info("trie fuckup")
				return nil, false
			}

			if _, ok := t.children[path[0]]; !ok {
				// sub-tree doesn't contain the path you asked for
				log.Info("trie fuckup")
				return nil, false
			}

			log.V(2).Info("recursing")
			return t.children[path[0]].Get(path[1:]...)
		}
	}
}

func (t *Trie) GetSubTree(path ...string) (*Trie, bool) {
	log := t.log.WithName("Get")
	log.V(2).Info("called", "path", path, "leaf?", t.leaf)

	if len(path) == 0 { // Asking for a leaf
		// makes no sense
		log.Info("trie fuckup")
		return nil, false
	} else { // Asking to resurse
		if t.leaf { // At a leaf
			// makes no sense
			log.Info("trie fuckup")
			return nil, false
		} else { // Not at a leaf
			if t.children == nil {
				// empty sub-tree; no value, no childen (TODO fix up this damn data type. The type itself shouldn't allow it (golang permitting), else enforce at insert time. Abstract the .leaf into a .Leaf() which should just be able to look at run-time type or children length?)
				log.Info("trie fuckup")
				return nil, false
			}

			if _, ok := t.children[path[0]]; !ok {
				// sub-tree doesn't contain the path you asked for
				log.Info("trie fuckup")
				return nil, false
			}

			return t.children[path[0]], true
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

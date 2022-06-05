package trie

func (t *Trie) Walk(cb func(path []string, value Value)) {
	t.walkInternal(cb, []string{})
}
func (t *Trie) walkInternal(cb func(path []string, value Value), path []string) {
	log := t.log.WithName("walkInternal")
	log.V(2).Info("called", "path", path, "leaf?", t.leaf)

	if t.leaf {
		cb(path, t.value)
	} else {
		if t.children == nil {
			panic("Invalid trie")
		}

		cb(path, some{""})
	}
	// ie it's ok for a tree node to have both a value and children
	// - this is used eg for one data source to give a basic value for something (like total core count) and another one to add children (like performance & efficiency core count)
	// - TODO: make this explicit in the Type if possible, and certainly in the docs + insert + get methods

	log.V(2).Info("recursing")
	for name, c := range t.children {
		c.walkInternal(cb, append(path, name))
	}
}

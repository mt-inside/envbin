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

		log.V(2).Info("recursing")
		for name, c := range t.children {
			c.walkInternal(cb, append(path, name))
		}
	}
}

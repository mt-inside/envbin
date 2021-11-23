package trie

import "github.com/go-logr/logr"

type InsertMsg struct {
	value Value
	path  []string
}

func Insert(value Value, path ...string) InsertMsg {
	return InsertMsg{value, path}
}

func BuildFromInsertMsgs(log logr.Logger, msgs <-chan InsertMsg) *Trie {
	t := NewTrie(log)

	for msg := range msgs {
		t.Insert(msg.value, msg.path...)
	}

	return t
}

func BuildFromSyncFn(log logr.Logger, f func(chan<- InsertMsg)) *Trie {
	c := make(chan InsertMsg)

	go func() {
		f(c)
		close(c)
	}()

	return BuildFromInsertMsgs(log, c)
}

func PrefixChan(vals chan<- InsertMsg, prefix ...string) chan<- InsertMsg {
	prefixed := make(chan InsertMsg)

	go func() {
		for v := range prefixed {
			vals <- InsertMsg{v.value, append(prefix, v.path...)}
		}
	}()

	return prefixed
}

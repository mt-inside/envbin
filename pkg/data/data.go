package data

import (
	"context"
	"time"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data/trie"
)

type pluginfn func(ctx context.Context, log logr.Logger, t *trie.Trie)

var (
	plugins []pluginfn
)

func RegisterPlugin(p pluginfn) {
	plugins = append(plugins, p)
}

func GetData(ctx context.Context, log logr.Logger) *trie.Trie {
	t := trie.NewTrie(log.WithName("trie"))

	for _, p := range plugins {
		go p(ctx, log, t)
	}
	<-ctx.Done()
	time.Sleep(1 * time.Second) // FIXME: HACK!!!
	// MVP: The datas should return tries, which are merged in series up here (avoid the race too)
	// Ideally: The datas should punt arrays of Element over a channel, and this func loops over the channel and inserts them. Not worth the optimisation, cause merging is instant, unless we wanna stream the result to the client (and stream partial updates)

	return t
}

package data

import (
	"context"
	"net/http"
	"time"

	"github.com/go-logr/logr"
)

var (
	plugins []func(ctx context.Context, log logr.Logger, t *Trie)
)

func GetData(ctx context.Context, log logr.Logger) *Trie {
	t := NewTrie(log.WithName("trie"))

	for _, p := range plugins {
		go p(ctx, log, t)
	}
	<-ctx.Done()
	time.Sleep(1 * time.Second) // FIXME: HACK!!!
	// MVP: The datas should return tries, which are merged in series up here (avoid the race too)
	// Ideally: The datas should punt arrays of Element over a channel, and this func loops over the channel and inserts them. Not worth the optimisation, cause merging is instant, unless we wanna stream the result to the client (and stream partial updates)

	return t
}
func GetDataWithRequest(ctx context.Context, log logr.Logger, r *http.Request) *Trie {
	t := GetData(ctx, log)
	getRequestData(ctx, log, t, r)

	return t
}

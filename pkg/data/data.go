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
	t := NewTrie()

	for _, p := range plugins {
		go p(ctx, log, t)
	}
	time.Sleep(6 * time.Second) // FIXME: HACK!!! The datas should punt arrays of Element over a channel, and this func loops over the channel and inserts them

	return t
}
func GetDataWithRequest(ctx context.Context, log logr.Logger, r *http.Request) *Trie {
	t := GetData(ctx, log)
	getRequestData(ctx, log, t, r)

	return t
}

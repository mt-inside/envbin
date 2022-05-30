package fetchers

import (
	"context"
	"errors"
	"time"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getTestData)
}

func getTestData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	vals <- trie.Insert(trie.Some("Ok"), "Test", "Some")
	vals <- trie.Insert(trie.NotPresent(), "Test", "NotPresent")
	vals <- trie.Insert(trie.Forbidden(), "Test", "Forbidden")
	vals <- trie.Insert(trie.Timeout(time.Second), "Test", "Timeout")
	vals <- trie.Insert(trie.Error(errors.New("test error")), "Test", "Error")
}

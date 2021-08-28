package fetchers

import (
	"context"
	"errors"
	"time"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getTestData)
}

func getTestData(ctx context.Context, log logr.Logger, t *Trie) {
	t.Insert(Some("Ok"), "Test", "Some")
	t.Insert(NotPresent(), "Test", "NotPresent")
	t.Insert(Forbidden(), "Test", "Forbidden")
	t.Insert(Timeout(time.Second), "Test", "Timeout")
	t.Insert(Error(errors.New("test error")), "Test", "Error")
}

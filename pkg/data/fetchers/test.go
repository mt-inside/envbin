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

func getTestData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	vals <- Insert(Some("Ok"), "Test", "Some")
	vals <- Insert(NotPresent(), "Test", "NotPresent")
	vals <- Insert(Forbidden(), "Test", "Forbidden")
	vals <- Insert(Timeout(time.Second), "Test", "Timeout")
	vals <- Insert(Error(errors.New("test error")), "Test", "Error")
}

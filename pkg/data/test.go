package data

import (
	"context"
	"errors"
	"time"

	"github.com/go-logr/logr"
)

func init() {
	plugins = append(plugins, getTestData)
}

func getTestData(ctx context.Context, log logr.Logger, t *Trie) {
	t.Insert(Some{"Ok"}, "Test", "Some")
	t.Insert(NotPresent{}, "Test", "NotPresent")
	t.Insert(Forbidden{}, "Test", "Forbidden")
	t.Insert(Timeout{time.Second}, "Test", "Timeout")
	t.Insert(Error{errors.New("test error")}, "Test", "Error")
}

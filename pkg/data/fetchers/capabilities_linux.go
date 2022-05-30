package fetchers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/syndtr/gocapability/capability"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getCapsData)
}

func getCapsData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	caps, err := capability.NewPid2(0)
	if err != nil {
		log.trie.Error(err, "Can't construct caps object for current process")
		vals <- trie.Insert(trie.Error(err), "Processes", "0", "Capabilities")
		return
	}

	err = caps.Load()
	if err != nil {
		log.trie.Error(err, "Can't load caps for current process")
		vals <- trie.Insert(trie.Error(err), "Processes", "0", "Capabilities")
		return
	}

	vals <- trie.Insert(trie.Some(caps.String()), "Processes", "0", "Capabilities")
}

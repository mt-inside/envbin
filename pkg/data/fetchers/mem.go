package fetchers

import (
	"context"

	sigar "github.com/cloudfoundry/gosigar"
	"github.com/docker/go-units"
	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getMemData)
}

func getMemData(ctx context.Context, log logr.Logger, t *Trie) {
	mem := sigar.Mem{}
	err := mem.Get()
	if err != nil {
		log.Error(err, "Can't read memory information")
		return
	}

	t.Insert(Some(units.BytesSize(float64(mem.Total))), "Hardware", "Memory", "Total")
}

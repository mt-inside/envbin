package data

import (
	"context"

	sigar "github.com/cloudfoundry/gosigar"
	"github.com/docker/go-units"
	"github.com/go-logr/logr"
)

func init() {
	plugins = append(plugins, getMemData)
}

func getMemData(ctx context.Context, log logr.Logger, t *Trie) {
	mem := sigar.Mem{}
	err := mem.Get()
	if err != nil {
		log.Error(err, "Can't read memory information")
		return
	}

	t.Insert(Some{units.BytesSize(float64(mem.Total))}, "Hardware", "Memory", "Total")
}

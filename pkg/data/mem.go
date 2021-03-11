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
	mem.Get()

	t.Insert(units.BytesSize(float64(mem.Total)), "Hardware", "Memory", "Total")
}

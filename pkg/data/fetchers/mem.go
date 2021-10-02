package fetchers

import (
	"context"
	"strconv"

	sigar "github.com/cloudfoundry/gosigar"
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
		t.Insert(Error(err), "Hardware", "Memory")
		return
	}

	t.Insert(Some(strconv.FormatUint(mem.Total, 10)), "Hardware", "Memory", "Total")
}

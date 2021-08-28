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
	data.RegisterPlugin(getProcsData)
}

func getProcsData(ctx context.Context, log logr.Logger, t *Trie) {
	procs := sigar.ProcList{}
	err := procs.Get()
	if err != nil {
		log.Error(err, "Can't read memory information")
		return
	}

	t.Insert(Some(strconv.Itoa(len(procs.List)-1)), "OS", "ProcessesCount")
}

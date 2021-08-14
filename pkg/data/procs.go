package data

import (
	"context"
	"strconv"

	sigar "github.com/cloudfoundry/gosigar"
	"github.com/go-logr/logr"
)

func init() {
	plugins = append(plugins, getProcsData)
}

func getProcsData(ctx context.Context, log logr.Logger, t *Trie) {
	procs := sigar.ProcList{}
	err := procs.Get()
	if err != nil {
		log.Error(err, "Can't read memory information")
		return
	}

	t.Insert(Some{strconv.Itoa(len(procs.List) - 1)}, "OS", "ProcessesCount")
}

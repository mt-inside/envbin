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
	procs.Get()

	t.Insert(strconv.Itoa(len(procs.List)-1), "OS", "ProcessesCount")
}

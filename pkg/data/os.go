package data

import (
	"context"
	"runtime"

	sigar "github.com/cloudfoundry/gosigar"
	"github.com/go-logr/logr"
	"github.com/shirou/gopsutil/v3/host"
)

func init() {
	plugins = append(plugins, getOsData)
}

func getOsData(ctx context.Context, log logr.Logger, t *Trie) {
	uptime := sigar.Uptime{}
	uptime.Get()
	is, _ := host.Info()

	t.Insert(uptime.Format(), "OS", "Uptime")
	t.Insert(runtime.GOOS, "OS", "Kernel", "Type")
	t.Insert(is.KernelVersion, "OS", "Kernel", "Version")
}

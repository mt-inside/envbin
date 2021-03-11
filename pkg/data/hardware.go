package data

import (
	"context"
	"runtime"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/klauspost/cpuid"
	"github.com/shirou/gopsutil/v3/host"
)

func init() {
	plugins = append(plugins, getHardwareData)
}

func getHardwareData(ctx context.Context, log logr.Logger, t *Trie) {
	is, _ := host.Info()

	t.Insert(runtime.GOARCH, "Hardware", "CPU", "Arch")
	t.Insert(cpuid.CPU.BrandName, "Hardware", "CPU", "Model")
	t.Insert(strconv.Itoa(cpuid.CPU.PhysicalCores), "Hardware", "CPU", "Cores")
	t.Insert(strconv.Itoa(cpuid.CPU.LogicalCores), "Hardware", "CPU", "Threads")
	t.Insert(is.VirtualizationSystem+" "+is.VirtualizationRole, "Hardware", "Virtualisation")
}

package data

import (
	"context"
	"strconv"

	"github.com/docker/go-units"
	"github.com/go-logr/logr"
	"github.com/klauspost/cpuid/v2"
)

func init() {
	plugins = append(plugins, getCpuidData)
}

func getCpuidData(ctx context.Context, log logr.Logger, t *Trie) {
	t.Insert(cpuid.CPU.BrandName, "Hardware", "CPU", "Model")
	t.Insert(strconv.Itoa(cpuid.CPU.PhysicalCores), "Hardware", "CPU", "Cores")
	t.Insert(strconv.Itoa(cpuid.CPU.LogicalCores), "Hardware", "CPU", "Threads")
	t.Insert(units.BytesSize(float64(cpuid.CPU.Cache.L1D)), "Hardware", "CPU", "Cache", "L1D")
	t.Insert(units.BytesSize(float64(cpuid.CPU.Cache.L1I)), "Hardware", "CPU", "Cache", "L1I")
	t.Insert(units.BytesSize(float64(cpuid.CPU.Cache.L2)), "Hardware", "CPU", "Cache", "L2")
	t.Insert(units.BytesSize(float64(cpuid.CPU.Cache.L3)), "Hardware", "CPU", "Cache", "L3")
}

package fetchers

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/klauspost/cpuid/v2"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/enrichments"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getCpuidData)
}

func getCpuidData(ctx context.Context, log logr.Logger, t *Trie) {
	t.InsertTree(enrichments.EnrichCpuModel(ctx, log, cpuid.CPU.BrandName), "Hardware", "CPU", "Model")
	t.Insert(Some(cpuid.CPU.BrandName), "Hardware", "CPU", "Model", "Name")
	t.Insert(Some(strconv.Itoa(cpuid.CPU.PhysicalCores)), "Hardware", "CPU", "Cores")
	t.Insert(Some(strconv.Itoa(cpuid.CPU.LogicalCores)), "Hardware", "CPU", "Threads")
	t.Insert(Some(strconv.Itoa(cpuid.CPU.Cache.L1D)), "Hardware", "CPU", "Cache", "L1D")
	t.Insert(Some(strconv.Itoa(cpuid.CPU.Cache.L1I)), "Hardware", "CPU", "Cache", "L1I")
	t.Insert(Some(strconv.Itoa(cpuid.CPU.Cache.L2)), "Hardware", "CPU", "Cache", "L2")
	t.Insert(Some(strconv.Itoa(cpuid.CPU.Cache.L3)), "Hardware", "CPU", "Cache", "L3")
}

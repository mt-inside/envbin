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

func getCpuidData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	enrichments.EnrichCpuModel(ctx, log, cpuid.CPU.BrandName, PrefixChan(vals, "Hardware", "CPU", "Model"))
	vals <- Insert(Some(cpuid.CPU.BrandName), "Hardware", "CPU", "Model", "Name")
	vals <- Insert(Some(strconv.Itoa(cpuid.CPU.PhysicalCores)), "Hardware", "CPU", "Cores")
	vals <- Insert(Some(strconv.Itoa(cpuid.CPU.LogicalCores)), "Hardware", "CPU", "Threads")
	vals <- Insert(Some(strconv.Itoa(cpuid.CPU.Cache.L1D)), "Hardware", "CPU", "Cache", "L1D")
	vals <- Insert(Some(strconv.Itoa(cpuid.CPU.Cache.L1I)), "Hardware", "CPU", "Cache", "L1I")
	vals <- Insert(Some(strconv.Itoa(cpuid.CPU.Cache.L2)), "Hardware", "CPU", "Cache", "L2")
	vals <- Insert(Some(strconv.Itoa(cpuid.CPU.Cache.L3)), "Hardware", "CPU", "Cache", "L3")
}

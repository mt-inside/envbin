package fetchers

import (
	"context"
	"strconv"

	"github.com/go-logr/logr"
	"github.com/klauspost/cpuid/v2"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/enrichments"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getCpuidData)
}

func getCpuidData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	enrichments.EnrichCpuModel(ctx, log, cpuid.CPU.BrandName, trie.PrefixChan(vals, "Hardware", "CPU", "Model"))
	vals <- trie.Insert(trie.Some(cpuid.CPU.BrandName), "Hardware", "CPU", "Model", "Name")
	vals <- trie.Insert(trie.Some(strconv.Itoa(cpuid.CPU.PhysicalCores)), "Hardware", "CPU", "Cores")
	vals <- trie.Insert(trie.Some(strconv.Itoa(cpuid.CPU.LogicalCores)), "Hardware", "CPU", "Threads")
	vals <- trie.Insert(trie.Some(strconv.Itoa(cpuid.CPU.Cache.L1D)), "Hardware", "CPU", "Cache", "L1D")
	vals <- trie.Insert(trie.Some(strconv.Itoa(cpuid.CPU.Cache.L1I)), "Hardware", "CPU", "Cache", "L1I")
	vals <- trie.Insert(trie.Some(strconv.Itoa(cpuid.CPU.Cache.L2)), "Hardware", "CPU", "Cache", "L2")
	vals <- trie.Insert(trie.Some(strconv.Itoa(cpuid.CPU.Cache.L3)), "Hardware", "CPU", "Cache", "L3")
}

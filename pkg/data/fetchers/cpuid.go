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
	if cpuid.CPU.BrandName != "" {
		vals <- trie.Insert(trie.Some(cpuid.CPU.BrandName), "Hardware", "CPU", "Model", "Name")
		enrichments.EnrichCpuModel(ctx, log, cpuid.CPU.BrandName, trie.PrefixChan(vals, "Hardware", "CPU", "Model"))
	}
	vals <- trie.Insert(trie.Some(strconv.Itoa(cpuid.CPU.PhysicalCores)), "Hardware", "CPU", "Cores")
	vals <- trie.Insert(trie.Some(strconv.Itoa(cpuid.CPU.LogicalCores)), "Hardware", "CPU", "Threads")

	l1d := cpuid.CPU.Cache.L1D
	l1i := cpuid.CPU.Cache.L1I
	vals <- trie.Insert(trie.Some(strconv.Itoa(l1d)), "Hardware", "CPU", "Cache", "Individual", "Level1Data")
	vals <- trie.Insert(trie.Some(strconv.Itoa(l1i)), "Hardware", "CPU", "Cache", "Individual", "Level1Instruction")
	vals <- trie.Insert(trie.Some(strconv.Itoa(l1d+l1i)), "Hardware", "CPU", "Cache", "Individual", "Level1")
	vals <- trie.Insert(trie.Some(strconv.Itoa(cpuid.CPU.Cache.L2)), "Hardware", "CPU", "Cache", "Individual", "Level2")
	vals <- trie.Insert(trie.Some(strconv.Itoa(cpuid.CPU.Cache.L3)), "Hardware", "CPU", "Cache", "Individual", "Level3")
}

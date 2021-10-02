package enrichments

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-logr/logr"

	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func EnrichCpuModel(ctx context.Context, log logr.Logger, name string) *Trie {
	// https://en.wikipedia.org/wiki/List_of_Intel_CPU_microarchitectures
	// TODO: really should look up the whole model number using the regexs in one of these articles, and get: (uArch, based-on, generation, process)
	intelXeonGens := []string{ // https://en.wikipedia.org/wiki/List_of_Intel_Xeon_processors
		"Unknown",
		"Nehalem (Clarkdale / Arrandale / Lynnfield / Gulftown / Bloomfield / Clarksfield)", // 1
		"Sandy Bridge",        // 2
		"Ivy Bridge",          // 3, v2
		"Haswell",             // 4, v3
		"Broadwell",           // 5, v4
		"Skylake",             // 6, v5, scalable gen 1
		"Skylake / Kaby Lake", // 7, v6
		"Kaby Lake Refresh (Skylake) / Coffee Lake (Skylake) / Amber Lake (Skylake) / Whiskey Lake (Skylake) / Cannon Lake (Palm Cove)", // 8
		"Skylake / Coffee Lake Refresh (Skylake)",                                                              // 9
		"Cascade Lake (Skylake) / Ice Lake (Sunny Cove) / Comet Lake (Skylake) / Amber Lake Refresh (Skylake)", // 10, scalable gen 2 & 3
		"Rocket Lake (Cypress Cove) / Tiger Lake (Willow Cove)",
		"Alder Lake (Golden Cove)",
	}
	intelScalableXeonGens := []string{ // https://en.wikipedia.org/wiki/List_of_Intel_Xeon_processors
		"Unknown",
		"Skylake (6th gen)",                 // 6, v5, scalable gen 1
		"Cascade Lake (10th gen)",           // 10, scalable gen 2
		"Cooper Lake / Ice Lake (10th gen)", // 10, scalable gen 3
	}
	intelCoreGens := []string{ // https://en.wikipedia.org/wiki/Intel_Core
		"Unknown",
		"Nehalem (Clarkdale / Arrandale / Lynnfield / Gulftown / Bloomfield / Clarksfield)",
		"Sandy Bridge",
		"Ivy Bridge",
		"Haswell",
		"Broadwell",
		"Skylake",
		"Skylake / Kaby Lake",
		"Kaby Lake Refresh / Coffee Lake / Amber Lake / Whiskey Lake / Cannon Lake",
		"Skylake / Coffee Lake Refresh",
		"Cascade Lake / Ice Lake / Comet Lake / Amber Lake Refresh",
		"Tiger Lake / Rocket Lake",
		"Alder Lake",
	}
	amdEpycGens := map[string]string{
		"1": "zen/+",
		"2": "zen2",
		"3": "zen3",
	}
	amdRyzenGens := map[string]string{
		"1": "zen",
		"2": "zen+",
		"3": "zen2",
		"4": "zen2 APU",
		"5": "zen3",
	}

	if strings.Contains(name, "Xeon") {
		// https://www.intel.co.uk/content/www/uk/en/processors/processor-numbers-data-center.html
		if strings.Contains(name, "Platinum") || strings.Contains(name, "Gold") || strings.Contains(name, "Silver") || strings.Contains(name, "Bronze") {
			// "Intel Xeon Scalable Processors"
			var series, series_num, generation, sku, flags string
			fmt.Sscanf(name, "Intel(R) Xeon(R) %s %1s%1s%2s%1s CPU", &series, &series_num, &generation, &sku, &flags)

			t := NewTrie(log)
			t.Insert(Some(series), "Series")
			t.Insert(Some(sku), "SKU")
			t.Insert(Some(generation), "Generation")
			t.Insert(Some(flags), "Flags")

			nGen, err := strconv.Atoi(generation)
			if err != nil {
				t.Insert(Error(err), "Mircoarchitecture")
			} else {
				t.Insert(Some(intelScalableXeonGens[nGen]), "Mircoarchitecture")
			}
			return t
		} else {
			// "Intel Xeon Processors"
			var series, ways, socket, sku, generation string
			fmt.Sscanf(name, "Intel(R) Xeon(R) CPU %2s-%1s%1s%2s v%s", &series, &ways, &socket, &sku, &generation)

			t := NewTrie(log)
			t.Insert(Some(series), "Series")
			t.Insert(Some(ways), "Ways")
			t.Insert(Some(socket), "Socket")
			t.Insert(Some(sku), "SKU")
			t.Insert(Some(generation), "Generation")

			nGen, err := strconv.Atoi(generation)
			if err != nil {
				t.Insert(Error(err), "Mircoarchitecture")
			} else {
				t.Insert(Some(intelXeonGens[nGen+1]), "Mircoarchitecture")
			}
			return t
		}
	} else if strings.Contains(name, "Core(TM)") {
		var series, sku, generation, flags string
		fmt.Sscanf(name, "Intel(R) Core(TM) CPU %2s-%1s%3s%1s", &series, &generation, &sku, &flags)

		t := NewTrie(log)
		t.Insert(Some(series), "Series")
		t.Insert(Some(sku), "SKU")
		t.Insert(Some(generation), "Generation")
		t.Insert(Some(flags), "Flags")

		nGen, err := strconv.Atoi(generation)
		if err != nil {
			t.Insert(Error(err), "Mircoarchitecture")
		} else {
			t.Insert(Some(intelCoreGens[nGen]), "Mircoarchitecture")
		}
		return t
	} else if strings.Contains(name, "EPYC") {
		var series, sku, generation string
		fmt.Sscanf(name, "AMD EPYC %1s%2s%1s", &series, &sku, &generation)

		t := NewTrie(log)
		t.Insert(Some(series), "Series")
		t.Insert(Some(sku), "SKU")
		t.Insert(Some(generation), "Generation")
		t.Insert(Some(amdEpycGens[generation]), "Mircoarchitecture")
		return t
	} else if strings.Contains(name, "Ryzen") {
		var series, sku, generation string
		fmt.Sscanf(name, "AMD Ryzen %s %1s%3s", &series, &generation, &sku)

		t := NewTrie(log)
		t.Insert(Some(series), "Series")
		t.Insert(Some(sku), "SKU")
		t.Insert(Some(generation), "Generation")
		t.Insert(Some(amdRyzenGens[generation]), "Mircoarchitecture")
		return t
	}

	return nil
}

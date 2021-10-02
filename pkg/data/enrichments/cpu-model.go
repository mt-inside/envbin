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
	intelXeonGens := []string{
		"Unknown",
		"Nehalem",
		"Sandy Bridge",
		"Ivy Bridge",
		"Haswell",
		"Broadwell",
		"Skylake",
		"Kaby Lake",
		"Coffee Lake",
		"Cascade Lake",
		"Comet Lake",
		"Rocket Lake (Cypress Cove)",
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

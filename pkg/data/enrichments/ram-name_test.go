package enrichments

import (
	"context"
	"testing"

	"github.com/mt-inside/envbin/pkg/data/trie"
	"github.com/mt-inside/go-usvc"
	"github.com/yumaojun03/dmidecode/parser/memory"
)

func TestEnrichRamSpecs(t *testing.T) {
	// TODO: test logger support
	log := usvc.GetLogger(true, 0)

	cases := []struct {
		typ                memory.MemoryDeviceType
		busTransferRateMHz uint
		busWidthBits       uint
		results            map[string]string
	}{
		{
			memory.MemoryDeviceTypeDDR, 400, 64,
			map[string]string{"Bus Speed": "200", "Standard": "DDR-400", "Module": "PC-3200"},
		},
		{
			memory.MemoryDeviceTypeDDR4, 3600, 64,
			map[string]string{"Bus Speed": "1800", "Standard": "DDR4-3600", "Module": "PC4-28800"},
		},
	}

	for _, cse := range cases {
		res := trie.BuildFromSyncFn(
			log,
			/* Ugly go syntax for partial application */
			func(c chan<- trie.InsertMsg) {
				EnrichRamSpecs(context.Background(), log, cse.typ, cse.busTransferRateMHz, cse.busWidthBits, c)
			},
		)

		// TODO: compare lens to ensure no extra values in the result

		for k, v := range cse.results {
			val, ok := res.Get(k)
			if !ok {
				t.Errorf("Key missing from result; expected: %s.", k)
				continue
			}
			if v != val.Render() {
				t.Errorf("Answer was wrong for key %s; expected: %s, got: %v.", k, v, val)
			}
		}
	}
}

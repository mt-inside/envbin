package enrichments

import (
	"context"
	"testing"

	"github.com/mt-inside/envbin/pkg/data/trie"
	"github.com/mt-inside/go-usvc"
)

func TestFormatBase10(t *testing.T) {
	// TODO: test logger support
	log := usvc.GetLogger(true, 0)

	cases := []struct {
		name    string
		results map[string]string
	}{
		{
			"Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz",
			map[string]string{"Series": "i5", "SKU": "690", "Generation": "4", "Flags": "K"},
		},
		{
			"Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz",
			map[string]string{"Series": "i9", "SKU": "880", "Generation": "9", "Flags": "H"},
		},
		{
			"Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz",
			map[string]string{"Series": "E5", "Ways": "2", "Socket": "6", "SKU": "30", "Generation": "3"},
		},
		{
			"Intel(R) Xeon(R) Platinum 8168 CPU @ 2.70GHz",
			map[string]string{"Series": "Platinum", "SKU": "68", "Generation": "1"},
		},
		{
			"AMD Ryzen 9 5900X 12-Core Processor",
			map[string]string{"Series": "9", "Generation": "5", "SKU": "900"},
		},
		{
			"AMD EPYC 7713 64-Core Processor",
			map[string]string{"Series": "7", "SKU": "71", "Generation": "3"},
		},
		{
			"Random string",
			map[string]string{"Details": "Unknown"},
		},
	}

	for _, cse := range cases {
		res := trie.BuildFromSyncFn(
			log,
			/* Ugly go syntax for partial application */
			func(c chan<- trie.InsertMsg) {
				EnrichCpuModel(context.Background(), log, cse.name, c)
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

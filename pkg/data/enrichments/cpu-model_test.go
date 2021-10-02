package enrichments

import (
	"context"
	"testing"

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
			"Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz",
			map[string]string{"Series": "E5", "Ways": "2", "Socket": "6", "SKU": "30", "Generation": "3"},
		},
		{
			"AMD EPYC 7713 64-Core Processor",
			map[string]string{"Series": "7", "SKU": "71", "Generation": "3"},
		},
		{
			"AMD Ryzen 9 5900X 12-Core Processor",
			map[string]string{"Series": "9", "Generation": "5", "SKU": "900"},
		},
	}

	for _, cse := range cases {
		res := EnrichCpuModel(context.Background(), log, cse.name)

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

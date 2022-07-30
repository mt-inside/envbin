package fetchers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/google/go-tpm/tpm"
	"github.com/google/go-tpm/tpm2"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin("tpm", getTpmData)
}

func getTpmData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	if !tryTpmV1(vals) && !tryTpmV2(vals) {
		vals <- trie.Insert(trie.NotPresent(), "Hardware", "TPM")
	}
}

func tryTpmV1(vals chan<- trie.InsertMsg) bool {
	rwc, err := tpm.OpenTPM("/dev/tpm0")
	if err != nil {
		return false
	}

	defer rwc.Close()

	vals <- trie.Insert(trie.Some("1.x"), "Hardware", "TPM", "Version")

	manuf, err := tpm.GetManufacturer(rwc)
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "TPM", "Manufacturer")
		return false
	}
	vals <- trie.Insert(trie.Some(string(manuf)), "Hardware", "TPM", "Manufacturer")

	return true
}

func tryTpmV2(vals chan<- trie.InsertMsg) bool {
	rwc, err := tpm2.OpenTPM("/dev/tpm0")
	if err != nil {
		return false
	}
	defer rwc.Close()

	vals <- trie.Insert(trie.Some("2.0"), "Hardware", "TPM", "Version")

	manuf, err := tpm2.GetManufacturer(rwc)
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Hardware", "TPM", "Manufacturer")
		return false
	}
	vals <- trie.Insert(trie.Some(string(manuf)), "Hardware", "TPM", "Manufacturer")

	return true
}

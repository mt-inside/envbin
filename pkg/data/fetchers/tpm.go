package fetchers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/google/go-tpm/tpm"
	"github.com/google/go-tpm/tpm2"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getTpmData)
}

func getTpmData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	if !tryTpmV1(vals) && !tryTpmV2(vals) {
		vals <- Insert(NotPresent(), "Hardware", "TPM")
	}
}

func tryTpmV1(vals chan<- InsertMsg) bool {
	rwc, err := tpm.OpenTPM("/dev/tpm0")
	if err != nil {
		return false
	}

	defer rwc.Close()

	vals <- Insert(Some("1.x"), "Hardware", "TPM", "Version")

	manuf, err := tpm.GetManufacturer(rwc)
	if err != nil {
		vals <- Insert(Error(err), "Hardware", "TPM", "Manufacturer")
		return false
	}
	vals <- Insert(Some(string(manuf)), "Hardware", "TPM", "Manufacturer")

	return true
}

func tryTpmV2(vals chan<- InsertMsg) bool {
	rwc, err := tpm2.OpenTPM("/dev/tpm0")
	if err != nil {
		return false
	}
	defer rwc.Close()

	vals <- Insert(Some("2.0"), "Hardware", "TPM", "Version")

	manuf, err := tpm2.GetManufacturer(rwc)
	if err != nil {
		vals <- Insert(Error(err), "Hardware", "TPM", "Manufacturer")
		return false
	}
	vals <- Insert(Some(string(manuf)), "Hardware", "TPM", "Manufacturer")

	return true
}

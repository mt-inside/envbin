package fetchers

import (
	"context"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getFirmwareData)
}

func getFirmwareData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	_, err := os.Stat("/sys/firmware/efi")
	if os.IsNotExist(err) {
		vals <- Insert(Some("BIOS"), "Hardware", "Firmware", "BootType")
	} else {
		vals <- Insert(Some("EFI"), "Hardware", "Firmware", "BootType")

		files, err := filepath.Glob("/sys/firmware/efi/efivars/SecureBoot-*")
		if err != nil || len(files) != 1 {
			vals <- Insert(Error(err), "Hardware", "Firmware", "SecureBoot")
			return
		}

		bytes, err := os.ReadFile(files[0])
		if err != nil {
			vals <- Insert(Error(err), "Hardware", "Firmware", "SecureBoot")
			return
		}

		vals <- Insert(Some(strconv.Itoa(int(bytes[4]))), "Hardware", "Firmware", "SecureBoot")
	}
}

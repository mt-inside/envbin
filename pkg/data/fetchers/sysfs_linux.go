package fetchers

import (
	"context"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin("linux sysfs firmware", getFirmwareData)
}

func getFirmwareData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	_, err := os.Stat("/sys/firmware/efi")
	if os.IsNotExist(err) {
		vals <- trie.Insert(trie.Some("BIOS"), "Hardware", "Firmware", "BootType")
	} else {
		vals <- trie.Insert(trie.Some("EFI"), "Hardware", "Firmware", "BootType")

		files, err := filepath.Glob("/sys/firmware/efi/efivars/SecureBoot-*")
		if err != nil || len(files) != 1 {
			vals <- trie.Insert(trie.Error(err), "Hardware", "Firmware", "SecureBoot")
			return
		}

		bytes, err := os.ReadFile(files[0])
		if err != nil {
			vals <- trie.Insert(trie.Error(err), "Hardware", "Firmware", "SecureBoot")
			return
		}

		vals <- trie.Insert(trie.Some(strconv.Itoa(int(bytes[4]))), "Hardware", "Firmware", "SecureBoot")
	}
}

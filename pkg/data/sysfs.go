//+build linux

package data

import (
	"context"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-logr/logr"
)

func init() {
	plugins = append(plugins, getFirmwareData)
}

func getFirmwareData(ctx context.Context, log logr.Logger, t *Trie) {
	_, err := os.Stat("/sys/firmware/efi")
	if os.IsNotExist(err) {
		t.Insert("BIOS", "Hardware", "Firmware", "BootType")
	} else {
		t.Insert("EFI", "Hardware", "Firmware", "BootType")

		files, err := filepath.Glob("/sys/firmware/efi/efivars/SecureBoot-*")
		if err != nil || len(files) != 1 {
			t.Insert("Error", "Hardware", "Firmware", "SecureBoot")
			return
		}

		bytes, err := os.ReadFile(files[0])
		if err != nil {
			t.Insert("Error", "Hardware", "Firmware", "SecureBoot")
			return
		}

		t.Insert(strconv.Itoa(int(bytes[4])), "Hardware", "Firmware", "SecureBoot")
	}
}

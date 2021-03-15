//+build linux

package data

import (
	"context"
	"os"

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
	}
}

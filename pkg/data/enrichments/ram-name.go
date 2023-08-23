//go:generate stringer -type=memory.MemoryDeviceType -trimprefix=MemoryDeviceType

package enrichments

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/yumaojun03/dmidecode/parser/memory"

	"github.com/mt-inside/envbin/pkg/data/trie"
)

func EnrichRamSpecs(
	ctx context.Context, log logr.Logger,
	typ memory.MemoryDeviceType, busTransferRateMHz uint, busWidthBits uint,
	vals chan<- trie.InsertMsg) {

	supportedTypes := map[memory.MemoryDeviceType]bool{
		memory.MemoryDeviceTypeDDR:  true,
		memory.MemoryDeviceTypeDDR2: true,
		memory.MemoryDeviceTypeDDR3: true,
		memory.MemoryDeviceTypeDDR4: true,
		//memory.MemoryDeviceTypeDDR5: true,
	}

	if supportedTypes[typ] {
		gen := typ.String()[3:]
		busClockMHz := busTransferRateMHz / 2 // DDR
		dataPerSecondMB := busTransferRateMHz * busWidthBits / 8

		vals <- trie.Insert(trie.Some(fmt.Sprintf("%d", busClockMHz)), "Bus Speed")
		vals <- trie.Insert(trie.Some(fmt.Sprintf("%s-%d", typ.String(), busTransferRateMHz)), "Standard")
		vals <- trie.Insert(trie.Some(fmt.Sprintf("PC%s-%d", gen, dataPerSecondMB)), "Module")
	}
}

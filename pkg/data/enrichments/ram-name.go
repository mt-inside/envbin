//go:generate stringer -type=memory.MemoryDeviceType -trimprefix=MemoryDeviceType

package enrichments

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data/trie"
	"github.com/yumaojun03/dmidecode/parser/memory"
)

func EnrichRamSpecs(
	ctx context.Context, log logr.Logger,
	typ memory.MemoryDeviceType, busClockMHz uint, busWidthBits uint,
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
		transfersPerSecondMHz := busClockMHz * 2 // DDR
		dataPerSecondMB := transfersPerSecondMHz * busWidthBits / 8

		vals <- trie.Insert(trie.Some(fmt.Sprintf("%s-%d", typ.String(), transfersPerSecondMHz)), "Standard")
		vals <- trie.Insert(trie.Some(fmt.Sprintf("PC%s-%d", gen, dataPerSecondMB)), "Module")
	}
}

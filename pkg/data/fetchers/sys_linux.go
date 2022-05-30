package fetchers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/sys/unix"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getSysUnixData)
}

func getSysUnixData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	var si unix.Sysinfo_t
	err := unix.Sysinfo(&si)
	if err != nil {
		return
	}

	vals <- trie.Insert(trie.Some(time.Duration(time.Duration(si.Uptime)*time.Second).String()), "OS", "Uptime")
}

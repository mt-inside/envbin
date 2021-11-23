package fetchers

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/sys/unix"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getSysUnixData)
}

func getSysUnixData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	var si unix.Sysinfo_t
	err := unix.Sysinfo(&si)
	if err != nil {
		log.Error(err, "Can't read sysinfo information")
		return
	}

	vals <- Insert(Some(time.Duration(time.Duration(si.Uptime)*time.Second).String()), "OS", "Uptime")
}

//go:build linux
// +build linux

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

func getSysUnixData(ctx context.Context, log logr.Logger, t *Trie) {
	var si unix.Sysinfo_t
	err := unix.Sysinfo(&si)
	if err != nil {
		log.Error(err, "Can't read sysinfo information")
		return
	}

	t.Insert(Some(time.Duration(time.Duration(si.Uptime)*time.Second).String()), "OS", "Uptime")
}

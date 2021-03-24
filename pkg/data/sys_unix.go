// +build linux

package data

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/sys/unix"
)

func init() {
	plugins = append(plugins, getSysUnixData)
}

func getSysUnixData(ctx context.Context, log logr.Logger, t *Trie) {
	var si unix.Sysinfo_t
	unix.Sysinfo(&si)

	t.Insert(time.Duration(time.Duration(si.Uptime)*time.Second).String(), "OS", "Uptime")
}

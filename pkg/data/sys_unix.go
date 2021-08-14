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
	err := unix.Sysinfo(&si)
	if err != nil {
		log.Error(err, "Can't read memory information")
		return
	}

	t.Insert(Some{time.Duration(time.Duration(si.Uptime) * time.Second).String()}, "OS", "Uptime")
}

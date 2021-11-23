package fetchers

import (
	"context"
	"os"

	"github.com/go-logr/logr"
	"github.com/mattn/go-isatty"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getTerminalData)
}

func getTerminalData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	var tty string
	if isatty.IsTerminal(os.Stdout.Fd()) {
		if ttyDev, err := os.Readlink("/proc/self/fd/1"); err == nil {
			tty = ttyDev
		} else {
			tty = err.Error()
		}
	} else {
		tty = "n/a"
	}
	vals <- Insert(Some(tty), "Processes", "0", "Session", "TTY")
}

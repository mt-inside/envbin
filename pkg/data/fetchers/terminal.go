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

func getTerminalData(ctx context.Context, log logr.Logger, t *Trie) {
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

	t.Insert(Some(os.Getenv("XDG_SESSION_ID")), "Process", "Session", "ID")
	t.Insert(Some(os.Getenv("XDG_SESSION_CLASS")), "Process", "Session", "Class")
	t.Insert(Some(os.Getenv("XDG_SESSION_TYPE")), "Process", "Session", "Type")
	t.Insert(Some(os.Getenv("XDG_SEAT")), "Process", "Session", "Seat")
	t.Insert(Some(os.Getenv("XDG_VTNR")), "Process", "Session", "VT")
	t.Insert(Some(tty), "Process", "Session", "TTY")
}

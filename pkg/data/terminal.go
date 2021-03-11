package data

import (
	"context"
	"os"

	"github.com/go-logr/logr"
	"github.com/mattn/go-isatty"
)

func init() {
	plugins = append(plugins, getTerminalData)
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

	t.Insert(os.Getenv("XDG_SESSION_ID"), "Process", "Session", "ID")
	t.Insert(os.Getenv("XDG_SESSION_CLASS"), "Process", "Session", "Class")
	t.Insert(os.Getenv("XDG_SESSION_TYPE"), "Process", "Session", "Type")
	t.Insert(os.Getenv("XDG_SEAT"), "Process", "Session", "Seat")
	t.Insert(os.Getenv("XDG_VTNR"), "Process", "Session", "VT")
	t.Insert(tty, "Process", "Session", "TTY")
}

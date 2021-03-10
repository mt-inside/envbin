package data

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

func init() {
	plugins = append(plugins, getTerminalData)
}

func getTerminalData() map[string]string {
	data := map[string]string{}

	var tty string
	if isatty.IsTerminal(os.Stdout.Fd()) {
		if ttyDev, err := os.Readlink("/proc/self/fd/1"); err == nil {
			tty = ttyDev
		} else {
			tty = err.Error()
		}
	} else {
		tty = "no"
	}

	data["Session"] = fmt.Sprintf(
		"id %s, class %s, type %s, seat %s, vt %s, tty %s",
		os.Getenv("XDG_SESSION_ID"),
		os.Getenv("XDG_SESSION_CLASS"),
		os.Getenv("XDG_SESSION_TYPE"),
		os.Getenv("XDG_SEAT"),
		os.Getenv("XDG_VTNR"),
		tty,
	)

	return data
}

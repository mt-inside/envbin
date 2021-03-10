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

	data["Tty"] = fmt.Sprintf("%t", isatty.IsTerminal(os.Stdout.Fd()))
	data["Session"] = fmt.Sprintf(
		"id %s, seat %s, vt %s, class %s, type %s",
		os.Getenv("XDG_SESSION_ID"),
		os.Getenv("XDG_SEAT"),
		os.Getenv("XDG_VTNR"),
		os.Getenv("XDG_SESSION_CLASS"),
		os.Getenv("XDG_SESSION_TYPE"),
	)

	return data
}

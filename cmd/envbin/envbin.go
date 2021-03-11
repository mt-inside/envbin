package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

var mainOpts struct {
	DevMode bool `long:"dev-mode"`
}
var (
	flagParser = flags.NewParser(&mainOpts, flags.Default)
)

func main() {
	_, err := flagParser.Parse()
	if err != nil {
		fmt.Println("err")
		var e *flags.Error
		if errors.As(err, &e) {
			if e.Type == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		}
		panic(err)
	}
}

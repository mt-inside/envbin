package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/go-usvc"
)

var opts struct{}
var (
	flagParser = flags.NewParser(&opts, flags.Default)
)

func main() {
	log := usvc.GetLogger(false)
	log.Info(data.RenderBuildData())

	serveOpts.log = log

	//.oneshotOpts.log = log
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

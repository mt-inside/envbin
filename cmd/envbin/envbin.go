package main

import (
	"os"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/go-usvc"
	"github.com/urfave/cli/v2"
)

func main() {
	log := usvc.GetLogger(true)

	app := &cli.App{
		Name:     "envbin",
		Usage:    "Inspects and makes available information about its runtime environment",
		Version:  data.Version,
		Compiled: data.BuildTime(),

		UseShortOptionHandling: true,
		EnableBashCompletion:   true, // TODO not working

		Flags: []cli.Flag{},

		Metadata: map[string]interface{}{
			"log": log,
		},

		Commands: []*cli.Command{
			Oneshot,
			Serve,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

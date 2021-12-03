package main

import (
	"os"

	"github.com/mt-inside/envbin/pkg/data/fetchers"
	"github.com/mt-inside/go-usvc"
	"github.com/urfave/cli/v2"
)

func main() {
	log := usvc.GetLogger(true, 1)

	app := &cli.App{
		Name:     "envbin",
		Usage:    "Inspects and makes available information about its runtime environment",
		Version:  fetchers.Version,
		Compiled: fetchers.BuildTime(),

		UseShortOptionHandling: true,
		EnableBashCompletion:   true, // TODO not working

		Flags: []cli.Flag{},

		Metadata: map[string]interface{}{
			"log": log,
		},

		Commands: []*cli.Command{
			Dump,
			Serve,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

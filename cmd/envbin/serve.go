package main

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/mt-inside/envbin/pkg/middleware"
	"github.com/urfave/cli/v2"
)

var Serve = &cli.Command{
	Name:  "serve",
	Usage: "Serve data over the network",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Value: ":8080",
			Usage: "Listen address",
		},
	},

	Action: serve,
}

func serve(c *cli.Context) error {
	log := c.App.Metadata["log"].(logr.Logger)

	listenAddr := c.String("addr")

	rootMux := mux.NewRouter()
	rootMux.Use(middleware.LoggingMiddleware)

	//rootMux.Path("/health").HandlerFunc(healthHandler) TODO merge with badpod; proper probes
	//rootMux.Path("/ready").HandlerFunc(healthHandler) TODO recall use a struct and Handler() to get a log to these things

	log.Info("Serving", "addr", listenAddr)
	err := http.ListenAndServe(listenAddr, rootMux)

	// TODO: graceful shutdown (lower readiness - combine with badpod first)

	return err
}

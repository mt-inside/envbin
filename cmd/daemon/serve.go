package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
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

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	data := data.GetData(ctx, log) // TODO refresh on GET. TODO push updates to web UI (gin seems to support push)

	listenAddr := c.String("addr")

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, data)
	})

	err := r.Run(listenAddr)

	//rootMux.Path("/health").HandlerFunc(healthHandler) TODO merge with badpod; proper probes
	//rootMux.Path("/ready").HandlerFunc(healthHandler) TODO recall use a struct and Handler() to get a log to these things

	// TODO: graceful shutdown (lower readiness - combine with badpod first)

	return err
}

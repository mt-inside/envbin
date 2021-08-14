package main

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/renderers"
	"github.com/urfave/cli/v2"
)

var Oneshot = &cli.Command{
	Name:  "oneshot",
	Usage: "Write data to terminal",

	Action: oneshot,
}

func oneshot(c *cli.Context) error {
	log := c.App.Metadata["log"].(logr.Logger)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	data := data.GetData(ctx, log)
	renderers.RenderTTY(log, data)

	return nil
}

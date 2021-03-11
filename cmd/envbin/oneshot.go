package main

import (
	"context"
	"time"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/renderers"
	"github.com/mt-inside/go-usvc"
)

type oneshotCmd struct{}

var oneshotOpts oneshotCmd

func init() {
	if _, err := flagParser.AddCommand(
		"oneshot",
		"Print info to the terminal and exit",
		"Print info to the terminal and exit",
		&oneshotOpts,
	); err != nil {
		panic(err)
	}
}

func (*oneshotCmd) Execute(args []string) error {
	log := usvc.GetLogger(mainOpts.DevMode)
	log.Info(data.RenderBuildData())

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	data := data.GetData(ctx, log)
	renderers.RenderTTY(log, data)

	return nil
}

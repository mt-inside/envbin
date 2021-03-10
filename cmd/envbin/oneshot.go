package main

import (
	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/renderers"
)

type oneshotCmd struct {
	log logr.Logger
}

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
	renderers.RenderTTY()

	return nil
}

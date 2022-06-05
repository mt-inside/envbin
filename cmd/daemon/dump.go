package main

import (
	"context"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-logr/logr"
	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
	"github.com/urfave/cli/v2"
)

var Dump = &cli.Command{
	Name:  "dump",
	Usage: "Write data to terminal",

	Action: dump,
}

func dump(c *cli.Context) error {
	log := c.App.Metadata["log"].(logr.Logger)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	data := data.GetData(ctx, log)
	// TODO: cuirated list, and --all
	data.Walk(renderCb)

	return nil
}

var (
	whiteBold = color.New(color.FgHiWhite).Add(color.Bold)
	//white = color.New(color.FgHiWhite)
	norm = color.New(color.FgWhite)
	//grey      = color.New(color.FgHiBlack)
)

func renderCb(path []string, value trie.Value) {
	depth := len(path)
	if depth == 0 {
		return
	}
	norm.Print(strings.Repeat("  ", depth-1))
	whiteBold.Printf("%s: ", path[depth-1])
	norm.Printf("%s\n", value.Render())
}

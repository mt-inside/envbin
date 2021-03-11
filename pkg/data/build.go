package data

import (
	"context"
	"fmt"
	"runtime"

	"github.com/go-logr/logr"
)

const Binary = "envbin"

var (
	Version   string
	BuildTime string
)

func init() {
	plugins = append(plugins, getBuildData)
}

func getBuildData(ctx context.Context, log logr.Logger, t *Trie) {
	t.Insert(Version, "Build", "Version")
	t.Insert(BuildTime, "Build", "Time")
	t.Insert(runtime.Version(), "Build", "Runtime")
}

func RenderBuildData() string {
	return fmt.Sprintf("%s %s, built at %s with %s", Binary, Version, BuildTime, runtime.Version())
}

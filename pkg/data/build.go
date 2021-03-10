package data

import (
	"fmt"
	"runtime"
)

const Binary = "envbin"

var (
	Version   string
	BuildTime string
)

func init() {
	plugins = append(plugins, getBuildData)
}

func getBuildData() map[string]string {
	data := map[string]string{}

	data["Version"] = Version
	data["BuildTime"] = BuildTime

	return data
}

func RenderBuildData() string {
	return fmt.Sprintf("%s %s, built at %s with %s", Binary, Version, BuildTime, runtime.Version())
}

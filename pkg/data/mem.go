package data

import (
	sigar "github.com/cloudfoundry/gosigar"
	"github.com/docker/go-units"
)

func init() {
	plugins = append(plugins, getMemData)
}

func getMemData() map[string]string {
	data := map[string]string{}

	mem := sigar.Mem{}
	mem.Get()

	data["MemTotal"] = units.BytesSize(float64(mem.Total))

	return data
}

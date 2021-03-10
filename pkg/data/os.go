package data

import (
	"runtime"

	sigar "github.com/cloudfoundry/gosigar"
	"github.com/shirou/gopsutil/v3/host"
)

func init() {
	plugins = append(plugins, getOsData)
}

func getOsData() map[string]string {
	data := map[string]string{}

	uptime := sigar.Uptime{}
	uptime.Get()
	is, _ := host.Info()

	data["OsUptime"] = uptime.Format()
	data["OsType"] = runtime.GOOS
	data["KernelVersion"] = is.KernelVersion
	data["GoVersion"] = runtime.Version()

	return data
}

package data

import (
	sigar "github.com/cloudfoundry/gosigar"
	"github.com/shirou/gopsutil/host"
	"runtime"
)

func getOsData() map[string]string {
	data := map[string]string{}


	uptime := sigar.Uptime{}
	uptime.Get()
	is, _ := host.Info()

	data["OsUptime"] = uptime.Format()
	data["OsType"] = runtime.GOOS
	data["OsVersion"] = is.KernelVersion
	data["GoVersion"] = runtime.Version()

	return data
}
package data

import (
	"github.com/shirou/gopsutil/v3/host"
)

func getOsDistributionData() map[string]string {
	data := map[string]string{}

	is, _ := host.Info()

	data["OsFamily"] = is.PlatformFamily
	data["OsDistro"] = is.Platform
	data["OsVersion"] = is.PlatformVersion

	return data
}

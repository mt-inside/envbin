package data

import (
	"github.com/joho/godotenv"
	"github.com/shirou/gopsutil/v3/host"
)

func init() {
	plugins = append(plugins, getOsDistributionData)
}

func getOsDistributionData() map[string]string {
	data := map[string]string{}

	is, _ := host.Info()

	// NB this is the distro in the CONTAINER. Distroless shows up as debian
	data["OsFamily"] = is.PlatformFamily
	data["OsDistro"] = is.Platform
	data["OsVersion"] = is.PlatformVersion

	osRelease, err := godotenv.Read("/etc/os-release")
	if err != nil {
		panic(err)
	}
	// /etc/os-release:PRETTY_VERSION seems to be universal
	data["OsReal"] = osRelease["PRETTY_NAME"]

	return data
}

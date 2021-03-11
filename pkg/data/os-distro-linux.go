package data

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/joho/godotenv"
	"github.com/shirou/gopsutil/v3/host"
)

func init() {
	plugins = append(plugins, getOsDistributionData)
}

func getOsDistributionData(ctx context.Context, log logr.Logger, t *Trie) {
	is, _ := host.Info()

	// NB this is the distro in the CONTAINER. Distroless shows up as debian
	t.Insert(is.PlatformFamily, "OS", "Distro", "Family")
	t.Insert(is.Platform, "OS", "Distro", "Name")
	t.Insert(is.PlatformVersion, "OS", "Distro", "Version")

	osRelease, err := godotenv.Read("/etc/os-release")
	if err != nil {
		panic(err)
	}
	// /etc/os-release:PRETTY_VERSION seems to be universal
	t.Insert(osRelease["PRETTY_NAME"], "OS", "Distro", "Release")
}

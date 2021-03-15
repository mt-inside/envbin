package data

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/shirou/gopsutil/v3/host"
)

func init() {
	plugins = append(plugins, getPsutilData)
}

func getPsutilData(ctx context.Context, log logr.Logger, t *Trie) {
	is, _ := host.Info()

	t.Insert(is.VirtualizationSystem+" "+is.VirtualizationRole, "Hardware", "Virtualisation")

	t.Insert(is.KernelVersion, "OS", "Kernel", "Version")

	// NB this is the distro in the CONTAINER. Distroless shows up as debian
	t.Insert(is.PlatformFamily, "OS", "Distro", "Family")
	t.Insert(is.Platform, "OS", "Distro", "Name")
	t.Insert(is.PlatformVersion, "OS", "Distro", "Version")
}

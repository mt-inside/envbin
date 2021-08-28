package fetchers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/shirou/gopsutil/v3/host"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getPsutilData)
}

func getPsutilData(ctx context.Context, log logr.Logger, t *Trie) {
	is, _ := host.Info()

	t.Insert(Some(is.VirtualizationSystem+" "+is.VirtualizationRole), "Hardware", "Virtualisation")

	t.Insert(Some(is.KernelVersion), "OS", "Kernel", "Version")

	// NB this is the distro in the CONTAINER. Distroless shows up as debian
	t.Insert(Some(is.PlatformFamily), "OS", "Distro", "Family")
	t.Insert(Some(is.Platform), "OS", "Distro", "Name")
	t.Insert(Some(is.PlatformVersion), "OS", "Distro", "Version")
}

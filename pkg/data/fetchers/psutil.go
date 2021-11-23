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

func getPsutilData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	is, err := host.Info()
	if err != nil {
		log.Error(err, "Can't read PsUtil data")
		return
	}

	vals <- Insert(Some(is.VirtualizationSystem+" "+is.VirtualizationRole), "Hardware", "Virtualisation")

	vals <- Insert(Some(is.KernelVersion), "OS", "Kernel", "Version")

	// NB this is the distro in the CONTAINER. Distroless shows up as debian
	vals <- Insert(Some(is.PlatformFamily), "OS", "Distro", "Family")
	vals <- Insert(Some(is.Platform), "OS", "Distro", "Name")
	vals <- Insert(Some(is.PlatformVersion), "OS", "Distro", "Version")
}

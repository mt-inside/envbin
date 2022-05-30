package fetchers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/shirou/gopsutil/v3/host"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getPsutilData)
}

func getPsutilData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	is, err := host.Info()
	if err != nil {
		return
	}

	vals <- trie.Insert(trie.Some(is.VirtualizationSystem+" "+is.VirtualizationRole), "Hardware", "Virtualisation")

	vals <- trie.Insert(trie.Some(is.KernelVersion), "OS", "Kernel", "Version")

	// NB this is the distro in the CONTAINER. Distroless shows up as debian
	vals <- trie.Insert(trie.Some(is.PlatformFamily), "OS", "Distro", "Family")
	vals <- trie.Insert(trie.Some(is.Platform), "OS", "Distro", "Name")
	vals <- trie.Insert(trie.Some(is.PlatformVersion), "OS", "Distro", "Version")
}

package fetchers

import (
	"context"
	"runtime"

	"github.com/go-logr/logr"
	"github.com/shirou/gopsutil/v3/host"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin("psutil", getPsutilData)
}

func getPsutilData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	is, err := host.Info()
	if err != nil {
		return
	}

	vals <- trie.Insert(trie.Some(is.VirtualizationSystem+" "+is.VirtualizationRole), "Hardware", "Virtualisation")

	vals <- trie.Insert(trie.Some(is.KernelVersion), "OS", "Kernel", "Version")

	// NB this is the distro in the CONTAINER. Distroless shows up as debian
	if runtime.GOOS != "darwin" {
		// These aren't so useful / accurate on macos, and the macos-specific code pulls a lot of them in
		vals <- trie.Insert(trie.Some(is.PlatformFamily), "OS", "Distro", "Family")
		vals <- trie.Insert(trie.Some(is.Platform), "OS", "Distro", "Name")
	}
	vals <- trie.Insert(trie.Some(is.PlatformVersion), "OS", "Distro", "Version")
}

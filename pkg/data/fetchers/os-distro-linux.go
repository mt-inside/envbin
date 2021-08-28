//go:build linux
// +build linux

package fetchers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/joho/godotenv"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getOsDistributionData)
}

func getOsDistributionData(ctx context.Context, log logr.Logger, t *Trie) {
	osRelease, err := godotenv.Read("/etc/os-release")
	if err != nil {
		panic(err)
	}
	// /etc/os-release:PRETTY_VERSION seems to be universal
	t.Insert(Some(osRelease["PRETTY_NAME"]), "OS", "Distro", "Release")
}

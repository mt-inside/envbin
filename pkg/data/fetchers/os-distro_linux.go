package fetchers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/joho/godotenv"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin("os", getOsDistributionData)
}

func getOsDistributionData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	osRelease, err := godotenv.Read("/etc/os-release")
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "OS", "Distro")
		return
	}
	// /etc/os-release:PRETTY_VERSION seems to be universal
	vals <- trie.Insert(trie.Some(osRelease["PRETTY_NAME"]), "OS", "Distro", "Release")
}

package fetchers

import (
	"context"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getGcpData)
}

func getGcpData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	if !metadata.OnGCE() {
		vals <- trie.Insert(trie.NotPresent(), "Cloud", "GCP")
		return
	}

	vals <- trie.Insert(unwrapGcp(metadata.ProjectID()), "Cloud", "GCP", "AccountID")
	vals <- trie.Insert(unwrapGcp(metadata.Zone()), "Cloud", "GCP", "Zone")
	vals <- trie.Insert(unwrapGcp(metadata.InstanceID()), "Cloud", "GCP", "Instance", "ID")
	vals <- trie.Insert(unwrapGcp(metadata.InstanceName()), "Cloud", "GCP", "Instance", "Name")
	vals <- trie.Insert(unwrapGcpSlice(metadata.InstanceTags()), "Cloud", "GCP", "Instance", "Tags")
}

func unwrapGcp(s string, err error) trie.Value {
	if err != nil {
		return trie.Error(err)
	}
	return trie.Some(s)
}

func unwrapGcpSlice(s []string, err error) trie.Value {
	if err != nil {
		return trie.Error(err)
	}
	return trie.Some(strings.Join(s, ", "))
}

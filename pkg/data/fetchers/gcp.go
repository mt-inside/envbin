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
	data.RegisterPlugin("gcp", getGcpData)
}

func getGcpData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	if !metadata.OnGCE() {
		vals <- trie.Insert(trie.NotPresent(), "Cloud", "GCP")
		return
	}

	vals <- trie.Insert(unwrapGcp(metadata.ProjectIDWithContext(ctx)), "Cloud", "GCP", "AccountID")
	vals <- trie.Insert(unwrapGcp(metadata.ZoneWithContext(ctx)), "Cloud", "GCP", "Zone")
	vals <- trie.Insert(unwrapGcp(metadata.InstanceIDWithContext(ctx)), "Cloud", "GCP", "Instance", "ID")
	vals <- trie.Insert(unwrapGcp(metadata.InstanceNameWithContext(ctx)), "Cloud", "GCP", "Instance", "Name")
	vals <- trie.Insert(unwrapGcpSlice(metadata.InstanceTagsWithContext(ctx)), "Cloud", "GCP", "Instance", "Tags")
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

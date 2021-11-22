package fetchers

import (
	"context"
	"strings"

	"cloud.google.com/go/compute/metadata"
	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getGcpData)
}

func getGcpData(ctx context.Context, log logr.Logger, t *Trie) {
	if !metadata.OnGCE() {
		t.Insert(NotPresent(), "Cloud", "GCP")
		return
	}
	t.Insert(Some("GCP"), "Cloud", "Provider")

	t.Insert(unwrapGcp(metadata.ProjectID()), "Cloud", "AccountID")
	t.Insert(unwrapGcp(metadata.Zone()), "Cloud", "Zone")
	t.Insert(unwrapGcp(metadata.InstanceID()), "Cloud", "InstanceID")
	t.Insert(unwrapGcp(metadata.InstanceName()), "Cloud", "InstanceName")
	t.Insert(unwrapGcpSlice(metadata.InstanceTags()), "Cloud", "InstanceTags")
}

func unwrapGcp(s string, err error) Value {
	if err != nil {
		return Error(err)
	}
	return Some(s)
}

func unwrapGcpSlice(s []string, err error) Value {
	if err != nil {
		return Error(err)
	}
	return Some(strings.Join(s, ", "))
}

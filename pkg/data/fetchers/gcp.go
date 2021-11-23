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

func getGcpData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	if !metadata.OnGCE() {
		vals <- Insert(NotPresent(), "Cloud", "GCP")
		return
	}
	vals <- Insert(Some("GCP"), "Cloud", "Provider")

	vals <- Insert(unwrapGcp(metadata.ProjectID()), "Cloud", "AccountID")
	vals <- Insert(unwrapGcp(metadata.Zone()), "Cloud", "Zone")
	vals <- Insert(unwrapGcp(metadata.InstanceID()), "Cloud", "InstanceID")
	vals <- Insert(unwrapGcp(metadata.InstanceName()), "Cloud", "InstanceName")
	vals <- Insert(unwrapGcpSlice(metadata.InstanceTags()), "Cloud", "InstanceTags")
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

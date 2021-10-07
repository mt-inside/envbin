package fetchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	. "github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin(getAwsData)
}

func getAwsData(ctx context.Context, log logr.Logger, t *Trie) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		// Technically, no client configured (this will pass on a workstation that interracts with AWS)
		t.Insert(NotPresent(), "Cloud", "AWS")
		return
	}

	aws := imds.NewFromConfig(cfg)

	iid, err := aws.GetInstanceIdentityDocument(ctx, nil)
	if err != nil {
		t.Insert(NotPresent(), "Cloud", "AWS")
		return
	}
	t.Insert(Some("AWS"), "Cloud", "Provider")
	t.Insert(Some(iid.AccountID), "Cloud", "AccountID")
	t.Insert(Some(iid.Region), "Cloud", "Region")
	t.Insert(Some(iid.AvailabilityZone), "Cloud", "Zone")
	t.Insert(Some(iid.InstanceType), "Cloud", "Instance", "Type")
	t.Insert(Some(iid.ImageID), "Cloud", "Instance", "ImageID")
}

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

func getAwsData(ctx context.Context, log logr.Logger, vals chan<- InsertMsg) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		// Technically, no client configured (this will pass on a workstation that interracts with AWS)
		vals <- Insert(NotPresent(), "Cloud", "AWS")
		return
	}

	aws := imds.NewFromConfig(cfg)

	iid, err := aws.GetInstanceIdentityDocument(ctx, nil)
	if err != nil {
		vals <- Insert(NotPresent(), "Cloud", "AWS")
		return
	}
	vals <- Insert(Some(iid.AccountID), "Cloud", "AWS", "AccountID")
	vals <- Insert(Some(iid.Region), "Cloud", "AWS", "Region")
	vals <- Insert(Some(iid.AvailabilityZone), "Cloud", "AWS", "Zone")
	vals <- Insert(Some(iid.InstanceType), "Cloud", "AWS", "Instance", "Type")
	vals <- Insert(Some(iid.ImageID), "Cloud", "AWS", "Instance", "ImageID")
}

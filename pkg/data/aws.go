package data

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/go-logr/logr"
)

func init() {
	plugins = append(plugins, getAwsData)
}

func getAwsData(ctx context.Context, log logr.Logger, t *Trie) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return
	}

	aws := imds.NewFromConfig(cfg)

	iid, err := aws.GetInstanceIdentityDocument(ctx, nil)
	if err != nil {
		return
	}
	t.Insert("AWS", "Cloud", "Provider")
	t.Insert(iid.AccountID, "Cloud", "AccountID")
	t.Insert(iid.Region, "Cloud", "Region")
	t.Insert(iid.AvailabilityZone, "Cloud", "Zone")
	t.Insert(iid.InstanceType, "Cloud", "Instance", "Type")
	t.Insert(iid.ImageID, "Cloud", "Instance", "ImageID")
}

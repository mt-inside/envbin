package fetchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/go-logr/logr"

	"github.com/mt-inside/envbin/pkg/data"
	"github.com/mt-inside/envbin/pkg/data/trie"
)

func init() {
	data.RegisterPlugin("aws", getAwsData)
}

func getAwsData(ctx context.Context, log logr.Logger, vals chan<- trie.InsertMsg) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		// Technically, no client configured (this will pass on a workstation that interracts with AWS)
		vals <- trie.Insert(trie.NotPresent(), "Cloud", "AWS")
		return
	}

	aws := imds.NewFromConfig(cfg)

	iid, err := aws.GetInstanceIdentityDocument(ctx, nil)
	if err != nil {
		vals <- trie.Insert(trie.NotPresent(), "Cloud", "AWS")
		return
	}
	vals <- trie.Insert(trie.Some(iid.AccountID), "Cloud", "AWS", "AccountID")
	vals <- trie.Insert(trie.Some(iid.Region), "Cloud", "AWS", "Region")
	vals <- trie.Insert(trie.Some(iid.AvailabilityZone), "Cloud", "AWS", "Zone")
	vals <- trie.Insert(trie.Some(iid.InstanceType), "Cloud", "AWS", "Instance", "Type")
	vals <- trie.Insert(trie.Some(iid.ImageID), "Cloud", "AWS", "Instance", "ImageID")

	client := sts.New(sts.Options{})
	cid, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		vals <- trie.Insert(trie.Error(err), "Cloud", "AWS", "Identity")
		return
	}
	vals <- trie.Insert(trie.Some(*cid.Arn), "Cloud", "AWS", "Identity", "ARN")
	vals <- trie.Insert(trie.Some(*cid.UserId), "Cloud", "AWS", "Identity", "User ID")
	vals <- trie.Insert(trie.Some(*cid.Account), "Cloud", "AWS", "Identity", "Account") // The account that the identity is in, might not be the same as the account this compute resource is in
}

package data

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

func init() {
	plugins = append(plugins, getAwsData)
}

func getAwsData() map[string]string {
	data := map[string]string{}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	aws := imds.NewFromConfig(cfg)

	region, err := aws.GetRegion(context.TODO(), nil)
	if err != nil {
		return nil // TODO proper errors
	}
	data["AwsRegion"] = region.Region

	iid, err := aws.GetInstanceIdentityDocument(context.TODO(), nil)
	if err != nil {
		return nil // TODO proper errors
	}
	data["AwsAccountID"] = iid.AccountID
	data["AwsRegion"] = iid.Region
	data["AwsZone"] = iid.AvailabilityZone
	data["AwsInstanceType"] = iid.InstanceType
	data["AwsImageID"] = iid.ImageID

	return data
}

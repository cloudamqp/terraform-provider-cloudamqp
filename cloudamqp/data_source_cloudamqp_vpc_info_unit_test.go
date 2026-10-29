package cloudamqp

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Unit tests for the cloudamqp_vpc_info data source against the fake API in
// data_source_cloudamqp_account_unit_test.go. The info response has the same shape as the real
// API for an AWS VPC, with the security group as a nested object.

const vpcInfoDataSourceName = "data.cloudamqp_vpc_info.info"

var awsVpcInfo = map[string]any{
	"id":       "vpc-0a1b2c3d4e5f",
	"name":     "vpc-aws-20",
	"subnet":   "10.56.72.0/24",
	"owner_id": "123456789012",
	"security_group": map[string]any{
		"id":          "sg-0194712dff1eb5860",
		"name":        "default",
		"description": "default VPC security group",
		"owner_id":    "123456789012",
	},
}

func checkAwsVpcInfo() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(vpcInfoDataSourceName, "id", "vpc-0a1b2c3d4e5f"),
		resource.TestCheckResourceAttr(vpcInfoDataSourceName, "name", "vpc-aws-20"),
		resource.TestCheckResourceAttr(vpcInfoDataSourceName, "vpc_subnet", "10.56.72.0/24"),
		resource.TestCheckResourceAttr(vpcInfoDataSourceName, "owner_id", "123456789012"),
		resource.TestCheckResourceAttr(vpcInfoDataSourceName, "security_group_id", "sg-0194712dff1eb5860"),
	)
}

// TestVpcInfoDataSource_InstanceID: the info is read with instance_id.
func TestVpcInfoDataSource_InstanceID(t *testing.T) {
	t.Parallel()
	fake, server := newFakeAccountAPI(t)
	fake.vpcInfo["/api/instances/1234/vpc-peering/info"] = awsVpcInfo

	accountUnitTest(t, []resource.TestStep{
		{
			Config: accountUnitTestConfig(server.URL, `
data "cloudamqp_vpc_info" "info" {
  instance_id = 1234
}
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAwsVpcInfo(),
				resource.TestCheckResourceAttr(vpcInfoDataSourceName, "instance_id", "1234"),
			),
		},
	})
}

// TestVpcInfoDataSource_VpcID: the info is read with vpc_id.
func TestVpcInfoDataSource_VpcID(t *testing.T) {
	t.Parallel()
	fake, server := newFakeAccountAPI(t)
	fake.vpcInfo["/api/vpcs/20/vpc-peering/info"] = awsVpcInfo

	accountUnitTest(t, []resource.TestStep{
		{
			Config: accountUnitTestConfig(server.URL, `
data "cloudamqp_vpc_info" "info" {
  vpc_id = "20"
}
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkAwsVpcInfo(),
				resource.TestCheckResourceAttr(vpcInfoDataSourceName, "vpc_id", "20"),
			),
		},
	})
}

// TestVpcInfoDataSource_MissingID: neither instance_id nor vpc_id gives an error.
func TestVpcInfoDataSource_MissingID(t *testing.T) {
	t.Parallel()
	_, server := newFakeAccountAPI(t)

	accountUnitTest(t, []resource.TestStep{
		{
			Config:      accountUnitTestConfig(server.URL, `data "cloudamqp_vpc_info" "info" {}`),
			ExpectError: regexp.MustCompile(`(?i)specify either instance_id or vpc_id`),
		},
	})
}

// TestVpcInfoDataSource_NotAWS: info without an AWS VPC identifier, e.g. a GCP VPC, gives an
// error.
func TestVpcInfoDataSource_NotAWS(t *testing.T) {
	t.Parallel()
	fake, server := newFakeAccountAPI(t)
	fake.vpcInfo["/api/vpcs/10/vpc-peering/info"] = map[string]any{
		"name": "vpc-gcp-10", "network": "projects/p/global/networks/n", "subnet": "10.56.73.0/24",
	}

	accountUnitTest(t, []resource.TestStep{
		{
			Config: accountUnitTestConfig(server.URL, `
data "cloudamqp_vpc_info" "info" {
  vpc_id = "10"
}
`),
			ExpectError: regexp.MustCompile(`(?i)find .*VPC identifier`),
		},
	})
}

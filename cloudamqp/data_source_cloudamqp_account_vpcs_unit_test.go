package cloudamqp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Unit tests for the cloudamqp_account_vpcs data source against the fake API in
// data_source_cloudamqp_account_unit_test.go.

const accountVpcsDataSourceName = "data.cloudamqp_account_vpcs.vpcs"

// TestAccountVpcsDataSource_Read: VPCs are listed in the API order, newest first, with the VPC
// name from the cloud provider. A VPC without tags gets an empty list, so length() works on it.
func TestAccountVpcsDataSource_Read(t *testing.T) {
	t.Parallel()
	fake, server := newFakeAccountAPI(t)
	fake.vpcs = []map[string]any{
		{"id": 20, "name": "newest", "region": "amazon-web-services::us-east-1", "subnet": "10.56.72.0/24",
			"tags": []string{"prod"}},
		{"id": 10, "name": "oldest", "region": "google-compute-engine::europe-west1", "subnet": "10.56.73.0/24",
			"tags": nil},
	}
	fake.vpcInfo["/api/vpcs/20/vpc-peering/info"] = map[string]any{"name": "vpc-aws-20"}
	fake.vpcInfo["/api/vpcs/10/vpc-peering/info"] = map[string]any{"name": "vpc-gcp-10"}

	accountUnitTest(t, []resource.TestStep{
		{
			Config: accountUnitTestConfig(server.URL, `
data "cloudamqp_account_vpcs" "vpcs" {}

output "oldest_tags_count" {
  value = length(data.cloudamqp_account_vpcs.vpcs.vpcs[1].tags)
}
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(accountVpcsDataSourceName, "id", "noId"),
				resource.TestCheckResourceAttr(accountVpcsDataSourceName, "vpcs.#", "2"),
				resource.TestCheckResourceAttr(accountVpcsDataSourceName, "vpcs.0.id", "20"),
				resource.TestCheckResourceAttr(accountVpcsDataSourceName, "vpcs.0.name", "newest"),
				resource.TestCheckResourceAttr(accountVpcsDataSourceName, "vpcs.0.region", "amazon-web-services::us-east-1"),
				resource.TestCheckResourceAttr(accountVpcsDataSourceName, "vpcs.0.subnet", "10.56.72.0/24"),
				resource.TestCheckResourceAttr(accountVpcsDataSourceName, "vpcs.0.tags.#", "1"),
				resource.TestCheckResourceAttr(accountVpcsDataSourceName, "vpcs.0.tags.0", "prod"),
				resource.TestCheckResourceAttr(accountVpcsDataSourceName, "vpcs.0.vpc_name", "vpc-aws-20"),
				resource.TestCheckResourceAttr(accountVpcsDataSourceName, "vpcs.1.id", "10"),
				resource.TestCheckResourceAttr(accountVpcsDataSourceName, "vpcs.1.name", "oldest"),
				resource.TestCheckResourceAttr(accountVpcsDataSourceName, "vpcs.1.vpc_name", "vpc-gcp-10"),
				resource.TestCheckOutput("oldest_tags_count", "0"),
			),
		},
	})
}

// TestAccountVpcsDataSource_NoVpcs: an account without VPCs gives an empty list.
func TestAccountVpcsDataSource_NoVpcs(t *testing.T) {
	t.Parallel()
	fake, server := newFakeAccountAPI(t)
	fake.vpcs = []map[string]any{}

	accountUnitTest(t, []resource.TestStep{
		{
			Config: accountUnitTestConfig(server.URL, `
data "cloudamqp_account_vpcs" "vpcs" {}

output "vpcs_count" {
  value = length(data.cloudamqp_account_vpcs.vpcs.vpcs)
}
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(accountVpcsDataSourceName, "id", "noId"),
				resource.TestCheckOutput("vpcs_count", "0"),
			),
		},
	})
}

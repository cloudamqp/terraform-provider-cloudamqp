package cloudamqp

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Unit tests for the cloudamqp_vpc_gcp_info data source against the fake API in
// data_source_cloudamqp_account_unit_test.go.

const vpcGcpInfoDataSourceName = "data.cloudamqp_vpc_gcp_info.info"

var gcpVpcInfo = map[string]any{
	"name":    "vpc-gcp-10",
	"network": "https://www.googleapis.com/compute/v1/projects/p/global/networks/vpc-gcp-10",
	"subnet":  "10.56.73.0/24",
}

func checkGcpVpcInfo() resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(vpcGcpInfoDataSourceName, "id", "vpc-gcp-10"),
		resource.TestCheckResourceAttr(vpcGcpInfoDataSourceName, "name", "vpc-gcp-10"),
		resource.TestCheckResourceAttr(vpcGcpInfoDataSourceName, "network",
			"https://www.googleapis.com/compute/v1/projects/p/global/networks/vpc-gcp-10"),
		resource.TestCheckResourceAttr(vpcGcpInfoDataSourceName, "vpc_subnet", "10.56.73.0/24"),
	)
}

// TestVpcGcpInfoDataSource_InstanceID: the info is read with instance_id, and sleep and timeout
// get their defaults.
func TestVpcGcpInfoDataSource_InstanceID(t *testing.T) {
	t.Parallel()
	fake, server := newFakeAccountAPI(t)
	fake.vpcInfo["/api/instances/1234/vpc-peering/info"] = gcpVpcInfo

	accountUnitTest(t, []resource.TestStep{
		{
			Config: accountUnitTestConfig(server.URL, `
data "cloudamqp_vpc_gcp_info" "info" {
  instance_id = 1234
}
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkGcpVpcInfo(),
				resource.TestCheckResourceAttr(vpcGcpInfoDataSourceName, "instance_id", "1234"),
				resource.TestCheckResourceAttr(vpcGcpInfoDataSourceName, "sleep", "10"),
				resource.TestCheckResourceAttr(vpcGcpInfoDataSourceName, "timeout", "1800"),
			),
		},
	})
}

// TestVpcGcpInfoDataSource_VpcIDNotReady: the info is read with vpc_id, and 400 responses while
// the VPC is not ready are retried with the configured sleep.
func TestVpcGcpInfoDataSource_VpcIDNotReady(t *testing.T) {
	t.Parallel()
	fake, server := newFakeAccountAPI(t)
	fake.vpcInfo["/api/vpcs/10/vpc-peering/info"] = gcpVpcInfo
	fake.notReady = 1

	accountUnitTest(t, []resource.TestStep{
		{
			Config: accountUnitTestConfig(server.URL, `
data "cloudamqp_vpc_gcp_info" "info" {
  vpc_id  = "10"
  sleep   = 1
  timeout = 30
}
`),
			Check: resource.ComposeAggregateTestCheckFunc(
				checkGcpVpcInfo(),
				resource.TestCheckResourceAttr(vpcGcpInfoDataSourceName, "vpc_id", "10"),
				resource.TestCheckResourceAttr(vpcGcpInfoDataSourceName, "sleep", "1"),
				resource.TestCheckResourceAttr(vpcGcpInfoDataSourceName, "timeout", "30"),
			),
		},
	})
}

// TestVpcGcpInfoDataSource_MissingID: neither instance_id nor vpc_id gives an error.
func TestVpcGcpInfoDataSource_MissingID(t *testing.T) {
	t.Parallel()
	_, server := newFakeAccountAPI(t)

	accountUnitTest(t, []resource.TestStep{
		{
			Config:      accountUnitTestConfig(server.URL, `data "cloudamqp_vpc_gcp_info" "info" {}`),
			ExpectError: regexp.MustCompile(`(?i)specify either instance_id or vpc_id`),
		},
	})
}

// TestVpcGcpInfoDataSource_NotFound: a 404 for the info gives an error instead of a crash. SDK v2
// crashed here, so this test only passes on the framework version.
func TestVpcGcpInfoDataSource_NotFound(t *testing.T) {
	t.Parallel()
	_, server := newFakeAccountAPI(t)

	accountUnitTest(t, []resource.TestStep{
		{
			Config: accountUnitTestConfig(server.URL, `
data "cloudamqp_vpc_gcp_info" "info" {
  vpc_id = "10"
}
`),
			ExpectError: regexp.MustCompile(`No VPC info found`),
		},
	})
}

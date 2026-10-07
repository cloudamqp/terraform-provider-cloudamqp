package cloudamqp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccVPCGCPInfoDataSource_Basic: Read GCP VPC info.
func TestAccVPCGCPInfoDataSource_Basic(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_vpc" "vpc" {
            name   = "TestAccVPCGCPInfoDataSource_Basic"
            region = "google-compute-engine::europe-north2"
            subnet = "10.56.73.0/24"
            tags   = ["vcr-test"]
          }

          data "cloudamqp_vpc_gcp_info" "vpc_info" {
            	vpc_id = cloudamqp_vpc.vpc.id
          }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("cloudamqp_vpc.vpc", "vpc_name", "data.cloudamqp_vpc_gcp_info.vpc_info", "name"),
					resource.TestCheckResourceAttr("data.cloudamqp_vpc_gcp_info.vpc_info", "vpc_subnet", "10.56.73.0/24"),
				),
			},
		},
	})
}

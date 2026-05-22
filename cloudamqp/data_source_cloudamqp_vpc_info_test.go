package cloudamqp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccVPCInfoDataSource_Basic: Read AWS VPC info.
func TestAccVPCInfoDataSource_Basic(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_vpc" "vpc" {
            name   = "TestAccVPCInfoDataSource_Basic"
            region = "amazon-web-services::eu-central-1"
            subnet = "10.56.72.0/24"
            tags   = ["aws"]
          }

          data "cloudamqp_vpc_info" "vpc_info" {
            	vpc_id = cloudamqp_vpc.vpc.id
          }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("cloudamqp_vpc.vpc", "vpc_name", "data.cloudamqp_vpc_info.vpc_info", "name"),
					resource.TestCheckResourceAttr("data.cloudamqp_vpc_info.vpc_info", "vpc_subnet", "10.56.72.0/24"),
				),
			},
		},
	})
}

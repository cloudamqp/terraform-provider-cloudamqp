package cloudamqp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccAccountVPCs_Basic: Read VPCs of an account.
func TestAccAccountVPCs_Basic(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
				  resource "cloudamqp_vpc" "vpc-01" {
					  name   = "TestAccAccountVPCs_Basic-01"
					  region = "amazon-web-services::us-east-1"
					  subnet = "10.56.72.0/24"
						tags   = ["vcr-test"]
          }

					resource "cloudamqp_vpc" "vpc-02" {
					  name   = "TestAccAccountVPCs_Basic-02"
					  region = "amazon-web-services::us-east-1"
					  subnet = "10.56.72.0/24"
						tags   = ["vcr-test"]
          }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc-01", "name", "TestAccAccountVPCs_Basic-01"),
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc-01", "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc-01", "subnet", "10.56.72.0/24"),
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc-01", "tags.#", "1"),
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc-01", "tags.0", "vcr-test"),

					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc-02", "name", "TestAccAccountVPCs_Basic-02"),
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc-02", "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc-02", "subnet", "10.56.72.0/24"),
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc-02", "tags.#", "1"),
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc-02", "tags.0", "vcr-test"),
				),
			},
			{
				Config: `
				  resource "cloudamqp_vpc" "vpc-01" {
					  name   = "TestAccAccountVPCs_Basic-01"
					  region = "amazon-web-services::us-east-1"
					  subnet = "10.56.72.0/24"
						tags   = ["vcr-test"]
          }

					resource "cloudamqp_vpc" "vpc-02" {
					  name   = "TestAccAccountVPCs_Basic-02"
					  region = "amazon-web-services::us-east-1"
					  subnet = "10.56.72.0/24"
						tags   = ["vcr-test"]
          }

          data "cloudamqp_account_vpcs" "this" {
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudamqp_account_vpcs.this", "vpcs.#", "2"),
					resource.TestCheckResourceAttr("data.cloudamqp_account_vpcs.this", "vpcs.0.name", "TestAccAccountVPCs_Basic-01"),
					resource.TestCheckResourceAttr("data.cloudamqp_account_vpcs.this", "vpcs.0.region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr("data.cloudamqp_account_vpcs.this", "vpcs.0.subnet", "10.56.72.0/24"),
					resource.TestCheckResourceAttr("data.cloudamqp_account_vpcs.this", "vpcs.0.tags.#", "1"),
					resource.TestCheckResourceAttr("data.cloudamqp_account_vpcs.this", "vpcs.0.tags.0", "vcr-test"),

					resource.TestCheckResourceAttr("data.cloudamqp_account_vpcs.this", "vpcs.1.name", "TestAccAccountVPCs_Basic-02"),
					resource.TestCheckResourceAttr("data.cloudamqp_account_vpcs.this", "vpcs.1.region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr("data.cloudamqp_account_vpcs.this", "vpcs.1.subnet", "10.56.72.0/24"),
					resource.TestCheckResourceAttr("data.cloudamqp_account_vpcs.this", "vpcs.1.tags.#", "1"),
					resource.TestCheckResourceAttr("data.cloudamqp_account_vpcs.this", "vpcs.1.tags.0", "vcr-test"),
				),
			},
		},
	})
}

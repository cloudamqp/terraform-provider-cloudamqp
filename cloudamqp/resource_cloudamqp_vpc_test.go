package cloudamqp

import (
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccVpc_Basic: Accept VPC and import.
func TestAccVpc_Basic(t *testing.T) {
	var (
		fileNames       = []string{"vpc"}
		vpcResourceName = "cloudamqp_vpc.vpc"

		params = map[string]string{
			"VpcName":   "TestAccVpc_Basic",
			"VpcRegion": "amazon-web-services::us-east-1",
			"VpcSubnet": "10.56.72.0/24",
			"VpcTags":   `["Terraform", "VCR-Test"]`,
		}

		paramsUpdated = map[string]string{
			"VpcName":   "TestAccVpc_Basic_Updated",
			"VpcRegion": "amazon-web-services::us-east-1",
			"VpcSubnet": "10.56.72.0/24",
			"VpcTags":   `["Terraform", "VCR-Test", "Updated"]`,
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vpcResourceName, "name", params["VpcName"]),
					resource.TestCheckResourceAttr(vpcResourceName, "region", params["VpcRegion"]),
					resource.TestCheckResourceAttr(vpcResourceName, "subnet", params["VpcSubnet"]),
					resource.TestCheckResourceAttr(vpcResourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(vpcResourceName, "tags.0", "Terraform"),
					resource.TestCheckResourceAttr(vpcResourceName, "tags.1", "VCR-Test"),
				),
			},
			{
				ResourceName:      vpcResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ExpectNonEmptyPlan: true,
				Config:             configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vpcResourceName, "name", paramsUpdated["VpcName"]),
					resource.TestCheckResourceAttr(vpcResourceName, "region", paramsUpdated["VpcRegion"]),
					resource.TestCheckResourceAttr(vpcResourceName, "subnet", paramsUpdated["VpcSubnet"]),
					resource.TestCheckResourceAttr(vpcResourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(vpcResourceName, "tags.0", "Terraform"),
					resource.TestCheckResourceAttr(vpcResourceName, "tags.1", "VCR-Test"),
					resource.TestCheckResourceAttr(vpcResourceName, "tags.2", "Updated"),
				),
			},
		},
	})
}

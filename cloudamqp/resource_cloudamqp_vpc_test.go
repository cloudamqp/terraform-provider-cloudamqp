package cloudamqp

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccVpc_Import: Test VPC creation and import
func TestAccVpc_Import(t *testing.T) {
	var (
		fileNames       = []string{"vpc"}
		vpcResourceName = "cloudamqp_vpc.vpc"

		params = map[string]string{
			"VpcName":   "TestAccVpc_Import",
			"VpcRegion": "amazon-web-services::us-east-1",
			"VpcSubnet": "10.56.72.0/24",
			"VpcTags":   `["Terraform", "VCR-Test"]`,
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
		},
	})
}

// TestAccVpc_DifferentRegions: Test VPC creation in different regions.
func TestAccVpc_DifferentRegions(t *testing.T) {
	regions := []string{
		"amazon-web-services::us-east-1",
		"amazon-web-services::us-west-2",
		"amazon-web-services::eu-west-1",
		"google-compute-engine::us-central1",
		"azure-arm::eastus",
	}

	for _, region := range regions {
		t.Run(region, func(t *testing.T) {
			params := map[string]string{
				"VpcName":   "TestAccVpc_Region_" + strings.ReplaceAll(region, "::", "_"),
				"VpcRegion": region,
				"VpcSubnet": "10.56.72.0/24",
				"VpcTags":   `["Terraform", "Region-Test"]`,
			}

			cloudamqpResourceTest(t, resource.TestCase{
				PreCheck:                  func() { testAccPreCheck(t) },
				PreventPostDestroyRefresh: true,
				Steps: []resource.TestStep{
					{
						Config: configuration.GetTemplatedConfig(t, []string{"vpc"}, params),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("cloudamqp_vpc.vpc", "region", region),
						),
					},
				},
			})
		})
	}
}

// TestAccVpc_DifferentSubnets: Test VPC creation with different CIDR blocks.
func TestAccVpc_DifferentSubnets(t *testing.T) {
	subnets := []string{
		"10.0.0.0/24",
		"172.16.0.0/24",
		"192.168.1.0/24",
		"10.56.72.0/24",
	}

	for i, subnet := range subnets {
		t.Run(subnet, func(t *testing.T) {
			params := map[string]string{
				"VpcName":   fmt.Sprintf("TestAccVpc_Subnet_%d", i),
				"VpcRegion": "amazon-web-services::us-east-1",
				"VpcSubnet": subnet,
				"VpcTags":   `["Terraform", "Subnet-Test"]`,
			}

			cloudamqpResourceTest(t, resource.TestCase{
				PreCheck:                  func() { testAccPreCheck(t) },
				PreventPostDestroyRefresh: true,
				Steps: []resource.TestStep{
					{
						Config: configuration.GetTemplatedConfig(t, []string{"vpc"}, params),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("cloudamqp_vpc.vpc", "subnet", subnet),
						),
					},
				},
			})
		})
	}
}

// TestAccVpc_TagsUpdate: Test updating VPC tags.
func TestAccVpc_TagsUpdate(t *testing.T) {
	var (
		fileNames       = []string{"vpc"}
		vpcResourceName = "cloudamqp_vpc.vpc"

		params = map[string]string{
			"VpcName":   "TestAccVpc_TagsUpdate",
			"VpcRegion": "amazon-web-services::us-east-1",
			"VpcSubnet": "10.56.72.0/24",
			"VpcTags":   `["Terraform", "VCR-Test"]`,
		}

		paramsUpdated = map[string]string{
			"VpcName":   "TestAccVpc_TagsUpdate",
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
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
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

// TestAccVpc_NoTags: Test VPC creation without tags.
func TestAccVpc_NoTags(t *testing.T) {
	params := map[string]string{
		"VpcName":   "TestAccVpc_NoTags",
		"VpcRegion": "amazon-web-services::us-east-1",
		"VpcSubnet": "10.56.72.0/24",
		"VpcTags":   `[]`,
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, []string{"vpc"}, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc", "name", params["VpcName"]),
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc", "tags.#", "0"),
				),
			},
		},
	})
}

// TestAccVpc_NameUpdate: Test updating VPC name only.
func TestAccVpc_NameUpdate(t *testing.T) {
	initialParams := map[string]string{
		"VpcName":   "TestAccVpc_NameUpdate_Initial",
		"VpcRegion": "amazon-web-services::us-east-1",
		"VpcSubnet": "10.56.72.0/24",
		"VpcTags":   `["Terraform", "Name-Test"]`,
	}

	updatedParams := map[string]string{
		"VpcName":   "TestAccVpc_NameUpdate_Updated",
		"VpcRegion": "amazon-web-services::us-east-1",
		"VpcSubnet": "10.56.72.0/24",
		"VpcTags":   `["Terraform", "Name-Test"]`,
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, []string{"vpc"}, initialParams),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc", "name", initialParams["VpcName"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, []string{"vpc"}, updatedParams),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc", "name", updatedParams["VpcName"]),
				),
			},
		},
	})
}

// TestAccVpc_InvalidCIDR: Test VPC creation with invalid CIDR blocks.
func TestAccVpc_InvalidCIDR(t *testing.T) {
	invalidCIDRs := []string{
		"invalid-cidr",
		"10.0.0.0/33",
		"256.256.256.256/24",
		"10.0.0.0",
	}

	for _, cidr := range invalidCIDRs {
		t.Run(cidr, func(t *testing.T) {
			params := map[string]string{
				"VpcName":   "TestAccVpc_InvalidCIDR",
				"VpcRegion": "amazon-web-services::us-east-1",
				"VpcSubnet": cidr,
				"VpcTags":   `["Terraform", "Invalid-CIDR-Test"]`,
			}

			cloudamqpResourceTest(t, resource.TestCase{
				PreCheck:                  func() { testAccPreCheck(t) },
				PreventPostDestroyRefresh: true,
				Steps: []resource.TestStep{
					{
						Config:      configuration.GetTemplatedConfig(t, []string{"vpc"}, params),
						ExpectError: regexp.MustCompile("Invalid CIDR"),
					},
				},
			})
		})
	}
}

// TestAccVpc_ReplacementOnRegionChange: Test VPC replacement when region changes.
func TestAccVpc_ReplacementOnRegionChange(t *testing.T) {
	initialParams := map[string]string{
		"VpcName":   "TestAccVpc_Replacement",
		"VpcRegion": "amazon-web-services::us-east-1",
		"VpcSubnet": "10.56.72.0/24",
		"VpcTags":   `["Terraform", "Replacement-Test"]`,
	}

	updatedParams := map[string]string{
		"VpcName":   "TestAccVpc_Replacement",
		"VpcRegion": "amazon-web-services::us-west-2",
		"VpcSubnet": "10.56.72.0/24",
		"VpcTags":   `["Terraform", "Replacement-Test"]`,
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, []string{"vpc"}, initialParams),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc", "region", initialParams["VpcRegion"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, []string{"vpc"}, updatedParams),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc", "region", updatedParams["VpcRegion"]),
				),
			},
		},
	})
}

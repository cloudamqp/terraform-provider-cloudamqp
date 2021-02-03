package cloudamqp

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceVpcInfo_Basic(t *testing.T) {
	instanceName := "cloudamqp_instance.instance"
	resourceName := "data.cloudamqp_vpc_info.vpc_info"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcInfoConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceVpcInfoExists(instanceName),
					resource.TestCheckResourceAttr(resourceName, "vpc_subnet", "10.56.72.0/24"),
				),
			},
		},
	})
}

func testAccCheckDataSourceVpcInfoExists(instanceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[instanceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", instanceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No resource id set")
		}
		instanceID, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		_, err := api.ReadVpcInfo(instanceID)
		if err != nil {
			return fmt.Errorf("Failed to fetch instance: %v", err)
		}

		return nil
	}
}

func testAccDataSourceVpcInfoConfigBasic() string {
	return `
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-vpc-info-ds-test"
			nodes 			= 1
			plan  			= "bunny-1"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
			vpc_subnet 	= "10.56.72.0/24"
		}

		data "cloudamqp_vpc_info" "vpc_info" {
			instance_id = cloudamqp_instance.instance.id
		}
	`
}

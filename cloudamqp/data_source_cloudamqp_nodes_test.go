package cloudamqp

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceNodes_Basic(t *testing.T) {
	instanceName := "cloudamqp_instance.instance"
	resourceName := "data.cloudamqp_nodes.nodes"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNodesConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceNodesExists(instanceName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nodes.0.rabbitmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resourceName, "nodes.0.running", "true"),
					resource.TestCheckResourceAttr(resourceName, "nodes.0.hipe", "false"),
				),
			},
		},
	})
}

func testAccCheckDataSourceNodesExists(instanceName, resourceName string) resource.TestCheckFunc {
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
		_, err := api.ReadNodes(instanceID)
		if err != nil {
			return fmt.Errorf("Failed to fetch instance: %v", err)
		}

		return nil
	}
}

func testAccDataSourceNodesConfigBasic() string {
	return `
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-nodes-ds-test"
			nodes 			= 1
			plan  			= "bunny-1"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		data "cloudamqp_nodes" "nodes" {
			instance_id = cloudamqp_instance.instance.id
		}
	`
}

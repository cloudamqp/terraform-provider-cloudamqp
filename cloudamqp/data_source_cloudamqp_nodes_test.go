package cloudamqp

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNodesDataSource_Basic(t *testing.T) {
	instance_name := "cloudamqp_instance.instance"
	resource_name := "data.cloudamqp_nodes.nodes"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNodesDataSourceConfig_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNodesDataSourceExists(instance_name, resource_name),
					resource.TestCheckResourceAttr(resource_name, "nodes.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "nodes.0.rabbitmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resource_name, "nodes.0.running", "true"),
					resource.TestCheckResourceAttr(resource_name, "nodes.0.hipe", "false"),
				),
			},
		},
	})
}

func testAccCheckNodesDataSourceExists(instance_name, resource_name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Resource %s not found", instance_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No resource id set")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)
		api := testAccProvider.Meta().(*api.API)
		_, err := api.ReadNodes(instance_id)
		if err != nil {
			return fmt.Errorf("Failed to fetch instance: %v", err)
		}

		return nil
	}
}

func testAccNodesDataSourceConfig_Basic() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-nodes-ds-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		data "cloudamqp_nodes" "nodes" {
			instance_id = cloudamqp_instance.instance.id
		}
		`)
}

package cloudamqp

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Basic instance test case. Creating dedicated AWS instance and do some minor updates.
func TestAccInstance_Basics(t *testing.T) {
	var (
		resourceName = "cloudamqp_instance.instance"
		name         = "terraform-before"
		region       = "amazon-web-services::us-east-1"
		plan         = "bunny-1"

		newName = "terraform-after"
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigBasic(name, region, plan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", plan),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
				),
			},
			{
				Config: testAccInstanceConfigBasic(newName, region, plan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", newName),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", plan),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
				),
			},
		},
	})
}

func testAccInstanceConfigBasic(name, region, plan string) string {
	log.Printf("[DBEUG] resource_instance::testAccInstanceConfig_Basic name: %s, region: %s, plan: %s", name, region, plan)
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "%s"
			nodes 			= 1
			plan 				= "%s"
			region 			= "%s"
			tags 				= ["terraform"]
		}
	`, name, plan, region)
}

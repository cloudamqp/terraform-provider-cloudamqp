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

func TestAccInstance_PlanChange(t *testing.T) {
	var (
		resourceName = "cloudamqp_instance.instance"
		name         = "Instance plan change"
		plan         = "squirrel-1"
		region       = "amazon-web-services::us-east-1"

		newPlan = "bunny-1"
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
				),
			},
			{
				Config: testAccInstanceConfigBasic(name, region, newPlan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", newPlan),
					resource.TestCheckResourceAttr(resourceName, "region", region),
				),
			},
		},
	})
}

func TestAccInstance_Upgrade(t *testing.T) {
	var (
		resourceName = "cloudamqp_instance.instance"
		name         = "Instance plan changes"
		plan         = "bunny-1"
		region       = "amazon-web-services::us-east-1"

		upgradePlan = "bunny-3"
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigBasic(name, region, plan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Instance plan changes"),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", plan),
				),
			},
			{
				Config: testAccInstanceConfigBasic(name, region, upgradePlan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodes", "3"),
					resource.TestCheckResourceAttr(resourceName, "plan", upgradePlan),
				),
			},
		},
	})
}

func TestAccInstance_Downgrade(t *testing.T) {
	var (
		resourceName = "cloudamqp_instance.instance"
		name         = "Instance plan changes"
		plan         = "bunny-3"
		region       = "amazon-web-services::us-east-1"

		downgradePlan = "bunny-1"
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigBasic(name, region, plan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "Instance plan changes"),
					resource.TestCheckResourceAttr(resourceName, "nodes", "3"),
					resource.TestCheckResourceAttr(resourceName, "plan", plan),
				),
			},
			{
				Config: testAccInstanceConfigBasic(name, region, downgradePlan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", downgradePlan),
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
			plan 				= "%s"
			region 			= "%s"
			tags 				= ["terraform"]
		}
	`, name, plan, region)
}

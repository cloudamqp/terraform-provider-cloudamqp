package cloudamqp

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Basic instance test case. Creating dedicated AWS instance and do some minor updates.
func TestAccInstance_Basics(t *testing.T) {
	resource_name := "cloudamqp_instance.instance"
	name := acctest.RandomWithPrefix("terraform")
	new_name := acctest.RandomWithPrefix("terraform")
	region := "amazon-web-services::us-east-1"
	plan := "bunny"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy(resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig_Basic(name, region, plan),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", name),
					resource.TestCheckResourceAttr(resource_name, "nodes", "1"),
					resource.TestCheckResourceAttr(resource_name, "plan", plan),
					resource.TestCheckResourceAttr(resource_name, "region", region),
					resource.TestCheckResourceAttr(resource_name, "rmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resource_name, "tags.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "tags.0", "terraform"),
				),
			},
			{
				Config: testAccInstanceConfig_Basic(new_name, region, plan),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", new_name),
					resource.TestCheckResourceAttr(resource_name, "nodes", "1"),
					resource.TestCheckResourceAttr(resource_name, "plan", plan),
					resource.TestCheckResourceAttr(resource_name, "region", region),
					resource.TestCheckResourceAttr(resource_name, "rmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resource_name, "tags.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "tags.0", "terraform"),
				),
			},
		},
	})
}

func testAccCheckInstanceExists(resource_name string) resource.TestCheckFunc {
	log.Printf("[DEBUG] resource_instance::testAccCheckInstanceExists resource: %s", resource_name)
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resource_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for instance resource")
		}
		instance_id := rs.Primary.ID

		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadInstance(instance_id)
		log.Printf("[DEBUG] resource_instance::testAccCheckInstanceExists data: %v", data)
		if err != nil {
			return fmt.Errorf("Error fetching item with resource %s. %s", resource_name, err)
		}
		return nil
	}
}

func testAccCheckInstanceDestroy(resource_name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Printf("[DEBUG] resource_instance::testAccCheckInstanceDestroy")
		api := testAccProvider.Meta().(*api.API)

		rs, ok := state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resource_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for instance resource")
		}
		instance_id := rs.Primary.ID

		_, err := api.ReadInstance(instance_id)
		if err == nil {
			return fmt.Errorf("Instance resource still exists")
		}
		invalidIdErr := "Invalid ID"
		expectedErr := regexp.MustCompile(invalidIdErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("Expected %s, got %s", invalidIdErr, err)
		}

		return nil
	}
}

func testAccInstanceConfig_Basic(name, region, plan string) string {
	log.Printf("[DBEUG] resource_instance::testAccInstanceConfig_Basic name: %s, region: %s, plan: %s", name, region, plan)
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "%s"
			nodes 			= 1
			plan 				= "%s"
			region 			= "%s"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
			vpc_subnet 	= "192.168.0.1/24"
		}
	`, name, plan, region)
}

func testAccInstanceConfig_Basic_Without_VPC(instance_name, name, region, plan, version string) string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "%s" {
			name 				= "%s"
			nodes 			= 1
			plan 				= "%s"
			region 			= "%s"
			rmq_version = "%s"
			tags 				= ["terraform"]
		}
	`, instance_name, name, plan, region, version)
}

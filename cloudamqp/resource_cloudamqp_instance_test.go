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
	resourceName := "cloudamqp_instance.instance"
	name := acctest.RandomWithPrefix("terraform")
	newName := acctest.RandomWithPrefix("terraform")
	region := "amazon-web-services::us-east-1"
	plan := "bunny-1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfigBasic(name, region, plan),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", plan),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "rmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
				),
			},
			{
				Config: testAccInstanceConfigBasic(newName, region, plan),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", newName),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", plan),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckResourceAttr(resourceName, "rmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
				),
			},
		},
	})
}

func testAccCheckInstanceExists(resourceName string) resource.TestCheckFunc {
	log.Printf("[DEBUG] resource_instance::testAccCheckInstanceExists resource: %s", resourceName)
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for instance resource")
		}
		instanceID := rs.Primary.ID

		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadInstance(instanceID)
		log.Printf("[DEBUG] resource_instance::testAccCheckInstanceExists data: %v", data)
		if err != nil {
			return fmt.Errorf("Error fetching item with resource %s. %s", resourceName, err)
		}
		return nil
	}
}

func testAccCheckInstanceDestroy(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Printf("[DEBUG] resource_instance::testAccCheckInstanceDestroy")
		api := testAccProvider.Meta().(*api.API)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for instance resource")
		}
		instanceID := rs.Primary.ID

		_, err := api.ReadInstance(instanceID)
		if err == nil {
			return fmt.Errorf("Instance resource still exists")
		}
		invalidIDErr := "Invalid ID"
		expectedErr := regexp.MustCompile(invalidIDErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("Expected %s, got %s", invalidIDErr, err)
		}

		return nil
	}
}

func testAccInstanceConfigBasic(name, region, plan string) string {
	log.Printf("[DBEUG] resource_instance::testAccInstanceConfig_Basic name: %s, region: %s, plan: %s", name, region, plan)
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "%s"
			nodes 			= 1
			plan 				= "%s"
			region 			= "%s"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}
	`, name, plan, region)
}

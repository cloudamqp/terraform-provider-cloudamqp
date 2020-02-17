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

func TestAccInstance_Basics(t *testing.T) {
	name := acctest.RandomWithPrefix("terraform")
	new_name := acctest.RandomWithPrefix("terraform")
	resource_name := "cloudamqp_instance.instance"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", name),
					resource.TestCheckResourceAttr(resource_name, "nodes", "1"),
					resource.TestCheckResourceAttr(resource_name, "plan", "bunny"),
					resource.TestCheckResourceAttr(resource_name, "region", "amazon-web-services::eu-north-1"),
					resource.TestCheckResourceAttr(resource_name, "rmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resource_name, "tags.#", "2"),
					resource.TestCheckResourceAttr(resource_name, "tags.0", "terraform"),
					resource.TestCheckResourceAttr(resource_name, "tags.1", "test"),
				),
			},
			{
				Config: testAccInstanceConfig_basic_update(new_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", new_name),
					resource.TestCheckResourceAttr(resource_name, "nodes", "1"),
					resource.TestCheckResourceAttr(resource_name, "plan", "bunny"),
					resource.TestCheckResourceAttr(resource_name, "region", "amazon-web-services::eu-north-1"),
					resource.TestCheckResourceAttr(resource_name, "rmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resource_name, "tags.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "tags.0", "terraform"),
				),
			},
		},
	})
}

func testAccCheckInstanceExists(resource string) resource.TestCheckFunc {
	log.Printf("[DEBUG] resource_instance::testAccCheckInstanceExists resource: %s", resource)
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		id := rs.Primary.ID

		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadInstance(id)
		log.Printf("[DEBUG] resource_instance::testAccCheckInstanceExists data: %v", data)
		if err != nil {
			return fmt.Errorf("Error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	log.Printf("[DEBUG] resource_instance::testAccCheckInstanceDestroy")
	api := testAccProvider.Meta().(*api.API)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudamqp_instance" {
			continue
		}

		_, err := api.ReadInstance(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Resource still exists")
		}
		invalidIdErr := "Invalid ID"
		expectedErr := regexp.MustCompile(invalidIdErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("Expected %s, got %s", invalidIdErr, err)
		}
	}

	return nil
}

func testAccInstanceConfig_basic(name string) string {
	log.Printf("[DEBUG] resource_instance::testAccInstanceConfig_basic name: %s", name)
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 					= "%s"
			nodes 				= 1
			plan 					= "bunny"
			region 				= "amazon-web-services::eu-north-1"
			rmq_version 	= "3.8.2"
			tags 					= ["terraform", "test"]
			vpc_subnet 		= "192.168.0.1/24"
		}
	`, name)
}

func testAccInstanceConfig_basic_update(new_name string) string {
	log.Printf("[DEBUG] resource_instance::testAccInstanceConfig_basic_update")
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "%s"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
			vpc_subnet  = "192.168.0.1/24"
		}
	`, new_name)
}

// func testAccInstanceConfig_Dedicated_AWS region: amazon-web-services::us-east-1, plan: bunny
// func testAccInstanceConfig_Shared_AWS region: amazon-web-services::us-east-1, plan

// func testAccInstanceConfig_AWS_Scale_Up region: amazon-web-services::us-east-1, plan: bunny -> rabbit, nodes: 1 -> 3
// func testAccInstanceConfig_AWS_Scale_Down region: amazon-web-services::us-east-1, plan rabbit -> bunny, nodes: 3 -> 1

// func testAccInstanceConfig_Dedicated_Azure, region: azure-arm::east-us
// func testAccInstanceConfig_Shared_Azure, region: azure-arm::east-us
// func testAccInstanceConfig_Dedicated_GCE, region: google-compute-engine::us-central1
// func testAccInstanceConfig_Shared_GCE, region: google-compute-engine::us-central1

// func testAccInstanceConfig_Dedicated_DO, region: digital-ocean::nyc3

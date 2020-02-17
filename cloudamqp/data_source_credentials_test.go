package cloudamqp

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCredentialsDataSource(t *testing.T) {
	resource_name := "cloudamqp_credentials.credentials"
	// var instance_id, username, password string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckCredentialsDataSource_destroy(resource_name),
		Steps: []resource.TestStep{
			// {
			// 	Config: testAccCredentialsDataSource_Config_Basic(),
			// },
			{
				Config: testAccCredentialsDataSource_config(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					// resource.TestCheckResourceAttr(resource_name, "instance_id", instance_id),
					resource.TestCheckResourceAttr(resource_name, "username", "hfkstucd"),
					resource.TestCheckResourceAttr(resource_name, "password", "xHm8vKLtbaGGtnEFmpL8ZXjYTfCPfzQb"),
				),
			},
			// {
			// 	ResourceName:      "cloudamqp_credentials.credentials-import",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	ImportStateIdFunc: importStateCredentials(resource_name),
			// },
		},
	})
}

func testAccCredentialsDataSource_Config_Basic() string {
	log.Printf("[DEBUG] data_source_credentials::testAccCredentialsDataSource_Config_Basic")
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance_credentials" {
			lifecycle {
				prevent_destroy = true
			}
			name 				= "terraform-credentials-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
			vpc_subnet = "192.168.0.1/24"
		}`)
}

func testAccCredentialsDataSource_config() string {
	log.Printf("[DEBUG] data_source_credentials::testAccCredentialsDataSource_configs")
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance_credentials" {
			name 				= "terraform-credentials-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
			vpc_subnet = "192.168.0.1/24"
		}

		data "cloudamqp_credentials" "credentials" {
			instance_id = cloudamqp_instance.instance_credentials.id
		}
		`)
}

// func fetchCredentials(resource_name string) resource.ImportStateIdFunc {
// 	return func(s *terraform.State) (string, error) {
// 		rs, ok := s.RootModule().Resources[resource_name]
// 		if !ok {
// 			return "", "", fmt.Errorf("Resource %s not found", resource_name)
// 		}
// 		if rs.Primary.ID == "" {
// 			return "", "", fmt.Errorf("No record ID is set")
// 		}
// 		return "", rs.Primary.ID, nil
// 	}
// }

func testAccCheckCredentialsDataSource_destroy(resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Control if extist! Need to remove the instance to remove the credentials.
		return nil
	}
}

func importStateId(resource_name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resource_name]
		if !ok {
			return "", fmt.Errorf("Resource %s not found", resource_name)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("No Record ID is set")
		}
		instance_id := rs.Primary.ID
		return instance_id, nil
	}
}

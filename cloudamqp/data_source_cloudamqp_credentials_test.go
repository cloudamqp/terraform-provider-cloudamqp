package cloudamqp

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceCredentials_Basic(t *testing.T) {
	instanceName := "cloudamqp_instance.instance"
	resourceName := "data.cloudamqp_credentials.credentials"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCredentielsConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceCredentialsExists(instanceName, resourceName),
				),
			},
		},
	})
}

func testAccCheckDataSourceCredentialsExists(instanceName, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[instanceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", instanceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No resource id set")
		}
		instanceID := rs.Primary.ID
		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadInstance(instanceID)
		if err != nil {
			return fmt.Errorf("Failed to fetch instance: %v", err)
		}

		r := regexp.MustCompile(`^.*:\/\/(?P<username>(.*)):(?P<password>(.*))@`)
		match := r.FindStringSubmatch(data["url"].(string))
		var username, password string
		for i, name := range r.SubexpNames() {
			if name == "username" {
				username = match[i]
			}
			if name == "password" {
				password = match[i]
			}
		}

		rs, ok = state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", resourceName)
		}
		if rs.Primary.Attributes["username"] == "" {
			return fmt.Errorf("No username attribute set for resource")
		}
		if rs.Primary.Attributes["password"] == "" {
			return fmt.Errorf("No password attribute set for resource")
		}

		if username != rs.Primary.Attributes["username"] {
			return fmt.Errorf("Username not equal")
		}
		if password != rs.Primary.Attributes["password"] {
			return fmt.Errorf("Password not equal")
		}
		return nil
	}
}

func testAccDataSourceCredentielsConfigBasic() string {
	return `
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-credentials-ds-test"
			nodes 			= 1
			plan  			= "bunny-1"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		data "cloudamqp_credentials" "credentials" {
			instance_id = cloudamqp_instance.instance.id
		}`
}

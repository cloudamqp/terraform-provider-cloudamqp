package cloudamqp

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCredentialsDataSource_Basic(t *testing.T) {
	instance_name := "cloudamqp_instance.instance"
	resource_name := "data.cloudamqp_credentials.credentials"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCredentialsDataSourceConfig_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCredentialsDataSourceExists(instance_name, resource_name),
				),
			},
		},
	})
}

func testAccCheckCredentialsDataSourceExists(instance_name, resource_name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Resource %s not found", instance_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No resource id set")
		}
		instance_id := rs.Primary.ID
		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadInstance(instance_id)
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

		rs, ok = state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource %s not found", resource_name)
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

func testAccCredentialsDataSourceConfig_Basic() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-credentials-ds-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		data "cloudamqp_credentials" "credentials" {
			instance_id = cloudamqp_instance.instance.id
		}
		`)
}

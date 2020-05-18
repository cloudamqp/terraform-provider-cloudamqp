package cloudamqp

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceNotificationDefault_Basic(t *testing.T) {
	instance_name := "cloudamqp_instance.instance"
	resource_name := "data.cloudamqp_notification.default_recipient"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationDefaultDataSourceConfig_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotificationDataSourceExists(instance_name, resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", "Default"),
				),
			},
		},
	})
}

func testAccCheckNotificationDataSourceExists(instance_name, resource_name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Resource %s not found", instance_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No resource id set")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		rs, ok = state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource %s not found", resource_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No resource id set")
		}
		alarm_id := rs.Primary.ID

		api := testAccProvider.Meta().(*api.API)
		_, err := api.ReadNotification(instance_id, alarm_id)
		if err != nil {
			return fmt.Errorf("Failed to fetch instance: %v", err)
		}

		return nil
	}
}

func testAccNotificationDefaultDataSourceConfig_Basic() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-notification-ds-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		data "cloudamqp_notification" "default_recipient" {
			instance_id = cloudamqp_instance.instance.id
			name 				= "Default"
		}
	`)
}

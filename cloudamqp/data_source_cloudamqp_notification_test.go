package cloudamqp

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceNotification_Basic(t *testing.T) {
	instanceName := "cloudamqp_instance.instance"
	resourceName := "data.cloudamqp_notification.default_recipient"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNotificationConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceNotificationExists(instanceName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "Default"),
				),
			},
		},
	})
}

func testAccCheckDataSourceNotificationExists(instanceName, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[instanceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", instanceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No resource id set")
		}
		instanceID, _ := strconv.Atoi(rs.Primary.ID)

		rs, ok = state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No resource id set")
		}
		alarmID := rs.Primary.ID

		api := testAccProvider.Meta().(*api.API)
		_, err := api.ReadNotification(instanceID, alarmID)
		if err != nil {
			return fmt.Errorf("Failed to fetch instance: %v", err)
		}

		return nil
	}
}

func testAccDataSourceNotificationConfigBasic() string {
	return `
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-notification-ds-test"
			nodes 			= 1
			plan  			= "bunny-1"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		data "cloudamqp_notification" "default_recipient" {
			instance_id = cloudamqp_instance.instance.id
			name 				= "Default"
		}
	`
}

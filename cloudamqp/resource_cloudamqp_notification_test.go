package cloudamqp

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNotificiaiton_Basic(t *testing.T) {
	instanceName := "cloudamqp_instance.instance"
	resourceName := "cloudamqp_notification.recipient_01"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNotificationDestroy(instanceName, resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotificationExists(instanceName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "email"),
					resource.TestCheckResourceAttr(resourceName, "value", "test@example.com"),
				),
			},
			{
				Config: testAccNotificationConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotificationExists(instanceName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "email"),
					resource.TestCheckResourceAttr(resourceName, "value", "notification@example.com"),
				),
			},
		},
	})
}

func testAccCheckNotificationExists(instanceName, resourceName string) resource.TestCheckFunc {
	log.Printf("[DEBUG] resource_notification::testAccCheckNotificationExists resource: %s", resourceName)

	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for notification resource")
		}
		recipientID := rs.Primary.ID

		rs, ok = state.RootModule().Resources[instanceName]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for instance resource")
		}
		instanceID, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		_, err := api.ReadNotification(instanceID, recipientID)
		if err != nil {
			return fmt.Errorf("Error fetching item with resource %s. %s", resourceName, err)
		}
		return nil
	}
}

func testAccCheckNotificationDestroy(instanceName, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Printf("[DEBUG] resource_notification::testAccCheckInstanceDestroy")
		api := testAccProvider.Meta().(*api.API)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for notification resource")
		}
		recipientID := rs.Primary.ID

		rs, ok = state.RootModule().Resources[instanceName]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for instance resource")
		}
		instanceID, _ := strconv.Atoi(rs.Primary.ID)

		data, err := api.ReadNotification(instanceID, recipientID)
		if data != nil || err == nil {
			return fmt.Errorf("Recipient still exists")
		}

		return nil
	}
}

func testAccNotificationConfigBasic() string {
	return `
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-notification-test"
			nodes 			= 1
			plan  			= "bunny-1"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		resource "cloudamqp_notification" "recipient_01" {
			instance_id = cloudamqp_instance.instance.id
			type = "email"
			value = "test@example.com"
			name = "test"
		}
		`
}

func testAccNotificationConfigUpdate() string {
	return `
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-notification-test"
			nodes 			= 1
			plan  			= "bunny-1"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		resource "cloudamqp_notification" "recipient_01" {
			instance_id = cloudamqp_instance.instance.id
			type = "email"
			value = "notification@example.com"
			name = "test"
		}
		`
}

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
	instance_name := "cloudamqp_instance.instance"
	resource_name := "cloudamqp_notification.recipient_01"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNotificationDestroy(instance_name, resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationConfig_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotificationExists(instance_name, resource_name),
					resource.TestCheckResourceAttr(resource_name, "type", "email"),
					resource.TestCheckResourceAttr(resource_name, "value", "test@example.com"),
				),
			},
			{
				Config: testAccNotificationConfig_Update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotificationExists(instance_name, resource_name),
					resource.TestCheckResourceAttr(resource_name, "type", "email"),
					resource.TestCheckResourceAttr(resource_name, "value", "notification@example.com"),
				),
			},
		},
	})
}

func testAccCheckNotificationExists(instance_name, resource_name string) resource.TestCheckFunc {
	log.Printf("[DEBUG] resource_notification::testAccCheckNotificationExists resource: %s", resource_name)

	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resource_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for notification resource")
		}
		recipient_id := rs.Primary.ID

		rs, ok = state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for instance resource")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		_, err := api.ReadNotification(instance_id, recipient_id)
		if err != nil {
			return fmt.Errorf("Error fetching item with resource %s. %s", resource_name, err)
		}
		return nil
	}
}

func testAccCheckNotificationDestroy(instance_name, resource_name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Printf("[DEBUG] resource_notification::testAccCheckInstanceDestroy")
		api := testAccProvider.Meta().(*api.API)

		rs, ok := state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource %s not found", resource_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for notification resource")
		}
		recipient_id := rs.Primary.ID

		rs, ok = state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for instance resource")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		data, err := api.ReadNotification(instance_id, recipient_id)
		if data != nil || err == nil {
			return fmt.Errorf("Recipient still exists")
		}

		return nil
	}
}

func testAccNotificationConfig_Basic() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-notification-test"
			nodes 			= 1
			plan  			= "bunny"
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
		`)
}

func testAccNotificationConfig_Update() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-notification-test"
			nodes 			= 1
			plan  			= "bunny"
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
		`)
}

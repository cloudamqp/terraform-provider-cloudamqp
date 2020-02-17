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
	//import_name := "cloudamqp_notification.default_recipient"
	instance_id := "cloudamqp_instance.instance_notification"
	resource_name := "cloudamqp_notification.recipient_01"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNotificationDestroy(instance_id, resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccNotification_Recipient(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotificationExists(instance_id, resource_name),
					resource.TestCheckResourceAttr(resource_name, "type", "email"),
					resource.TestCheckResourceAttr(resource_name, "value", "test@example.com"),
				),
			},
			{
				Config: testAccNotification_Recipient_Update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotificationExists(instance_id, resource_name),
					resource.TestCheckResourceAttr(resource_name, "type", "webhook"),
					resource.TestCheckResourceAttr(resource_name, "value", "http://example.com/webhook"),
				),
			},
			{
				ResourceName:      import_name,
				ImportState:       true,
				ImportStateIdFunc: testAccNotification_ImportStateId(instance_name, resource_name),
				//ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "account_id", "last_status"},
			},
		},
	})
}

// func TestAccNotification_Import(t *testing.T) {
// 	import_name := "cloudamqp_notification.default_recipient"
// 	instance_name := "cloudamqp_instance.recipient_import_instance"

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckNotificationDestroy(instance_name, import_name),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccNotification_Import(),
// 			},
// 			{
// 				//Config:            testAccNotification_Import2(),
// 				ResourceName:      import_name,
// 				ImportState:       true,
// 				ImportStateIdFunc: testAccNotification_ImportStateId(instance_name),
// 				//ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"instance_id", "account_id", "last_status"},
// 			},
// 		},
// 	})
// }

func testAccNotification_ImportStateId(instance_name, resource_names string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[instance_name]
		if !ok {
			return "", fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("No Record ID is set for instance")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadNotifications(instance_id)
		if err != nil {
			return "", err
		}

		if len(data) > 0 {
			log.Printf("[DEBUG] resource_notification::testAccNotification_ImportStateId data: %v", data)
			import_state_id := fmt.Sprintf("%v,%d", data[0]["id"], instance_id)
			log.Printf("[DEBUG] resource_notification::testAccNotification_ImportStateId import_state_id: %v", import_state_id)
			return import_state_id, nil
		}

		return "", nil
	}
}

func testAccCheckNotificationExists(instance_name, resource_name string) resource.TestCheckFunc {
	log.Printf("[DEBUG] resource_notification::testAccCheckNotificationExists resource: %s", resource_name)

	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resource_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		recipient_id := rs.Primary.ID

		rs, ok = state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set for instance")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadNotification(instance_id, recipient_id)
		log.Printf("[DEBUG] resource_notification::testAccCheckNotificationExists data: %v", data)
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
			return fmt.Errorf("No Record ID is set")
		}
		recipient_id := rs.Primary.ID

		rs, ok = state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set for instance")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		data, err := api.ReadNotification(instance_id, recipient_id)
		if data != nil || err == nil {
			return fmt.Errorf("Recipient still exists")
		}

		return nil
	}
}

func testAccNotification_Import() string {
	log.Printf("[DEBUG] resource_notification::testAccNotification_Import")
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "recipient_import_instance" {
			name 				= "terraform-notification-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
			vpc_subnet = "192.168.0.1/24"
		}

		resource "cloudamqp_notification" "default_recipient" {
			instance_id = cloudamqp_instance.recipient_import_instance.id
		}
		`)
}

func testAccNotification_Import2() string {
	log.Printf("[DEBUG] resource_notification::testAccNotification_Import")
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "recipient_import_instance" {
			name 				= "terraform-notification-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
			vpc_subnet = "192.168.0.1/24"
		}

		resource "cloudamqp_notification" "default_recipient" {
			instance_id = cloudamqp_instance.recipient_import_instance.id
			type = "email"
			value = ""
		}
		`)
}

func testAccNotification_Recipient() string {
	log.Printf("[DEBUG] resource_notification::testAccNotification_Recipient")
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance_notification" {
			name 				= "terraform-notification-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
			vpc_subnet = "192.168.0.1/24"
		}

		resource "cloudamqp_notification" "recipient_01" {
			instance_id = cloudamqp_instance.instance_notification.id
			type = "email"
			value = "test@example.com"
		}
		`)
}

func testAccNotification_Recipient_Update() string {
	log.Printf("[DEBUG] resource_notification::testAccNotification_Recipient_Update")
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance_notification" {
			name 				= "terraform-notification-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
			vpc_subnet = "192.168.0.1/24"
		}

		resource "cloudamqp_notification" "recipient_01" {
			instance_id = cloudamqp_instance.instance_notification.id
			type = "webhook"
			value = "http://example.com/webhook"
		}
		`)
}

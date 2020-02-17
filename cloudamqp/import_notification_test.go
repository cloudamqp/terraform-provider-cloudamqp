package cloudamqp

// import (
// 	"fmt"
// 	"log"
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-sdk/terraform"
// )

// func TestAccImportNotification_basic(t *testing.T) {
// 	resourceName := "cloudamqp_notification.recipient_01"

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckNotificationDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccNotification_Recipient(),
// 			},
// 			{
// 				ResourceName:      resourceName,
// 				ImportStateIdFunc: testAccNotification_ImportStateId(),
// 				ImportStateVerify: true,
// 				ImportStateVerifyIgnore: []string{
// 					"account_id",
// 					"last_status",
// 				},
// 			},
// 		},
// 	})
// }

// func testAccNotification_ImportStateId() resource.ImportStateIdFunc {
// 	return func(state *terraform.State) (string, error) {
// 		rs, ok := state.RootModule().Resources["cloudamqp_instance.instance_notification"]
// 		if !ok {
// 			return "", fmt.Errorf("") //Resource: %s not found", resource)
// 		}
// 		if rs.Primary.ID == "" {
// 			return "", fmt.Errorf("No Record ID is set")
// 		}
// 		instance_id := rs.Primary.ID

// 		rs, ok = state.RootModule().Resources["cloudamqp_notification.recipient_default"]
// 		if !ok {
// 			return "", fmt.Errorf("") //Resource: %s not found", resource)
// 		}
// 		if rs.Primary.ID == "" {
// 			return "", fmt.Errorf("No Record ID is set")
// 		}
// 		resource_id := rs.Primary.ID

// 		import_id := fmt.Sprintf("%s,%s", resource_id, instance_id)
// 		log.Printf("[DEBUG] cloudamqp::import_notification_test import_id: %v", import_id)
// 		return import_id, nil
// 	}
// }

// func testAccNotification_Recipient() string {
// 	log.Printf("[DEBUG] resource_notification::testAccNotification_Recipient")
// 	return fmt.Sprintf(`
// 		resource "cloudamqp_instance" "instance_notification" {
// 			name 				= "terraform-notification-test"
// 			nodes 			= 1
// 			plan  			= "bunny"
// 			region 			= "amazon-web-services::eu-north-1"
// 			rmq_version = "3.8.2"
// 			tags 				= ["terraform"]
// 			vpc_subnet = "192.168.0.1/24"
// 		}

// 		resource "cloudamqp_notification" "recipient_01" {
// 			instance_id = cloudamqp_instance.instance_notification.id
// 			type = "email"
// 			value = "test@example.com"
// 		}
// 		`)
// }

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

func TestAccSecurityFirewall_Basic(t *testing.T) {
	instance_name := "cloudamqp_instance.instance_firewall"
	resource_name := "cloudamqp_security_firewall.firewall"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSecurityFirewallDestroy(instance_name, resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityFirewallConfig_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityFirewallExists(instance_name, resource_name),
					resource.TestCheckResourceAttr(resource_name, "rule", "rabbitmq_web_mqtt"),
				),
			},
			// {
			// 	Config: testAccNotificationConfig_Update(),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckNotificationExists(instance_name, resource_name),
			// 		resource.TestCheckResourceAttr(resource_name, "type", "webhook"),
			// 		resource.TestCheckResourceAttr(resource_name, "value", "http://example.com/webhook"),
			// 	),
			// },
		},
	})
}

func testAccCheckSecurityFirewallExists(instance_name, resource_name string) resource.TestCheckFunc {
	log.Printf("[DEBUG] resource_plugin::testAccCheckPluginEnabled resource: %s", resource_name)

	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resource_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		plugin_name := rs.Primary.Attributes["name"]

		rs, ok = state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set for instance")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadPlugin(instance_id, plugin_name)
		log.Printf("[DEBUG] resource_plugin::testAccCheckPluginEnabled data: %v", data)
		if err != nil {
			return fmt.Errorf("Error fetching item with resource %s. %s", resource_name, err)
		}
		if data["enabled"] == false {
			return fmt.Errorf("Error resource: %s not enabled", resource_name)
		}
		return nil
	}
}

func testAccSecurityFirewallDestroy(instance_name, resource_name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Printf("[DEBUG] resource_plugins::testAccCheckPluginDisable")
		api := testAccProvider.Meta().(*api.API)

		rs, ok := state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource %s not found", resource_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		plugin_name := rs.Primary.Attributes["name"]

		rs, ok = state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set for instance")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		data, err := api.ReadNotification(instance_id, plugin_name)
		if data != nil || err == nil {
			return fmt.Errorf("Recipient still exists")
		}
		if data["enabled"] == true {
			return fmt.Errorf("Error resource: %s not disabled", resource_name)
		}
		return nil
	}
}

func testAccSecurityFirewallConfig_Basic() string {
	log.Printf("[DEBUG] resource_plugins::testAccPluginConfig_Basic")
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance_plugin" {
			name 				= "terraform-plugin-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		resource "cloudamqp_plugin" "mqtt_plugin" {
			instance_id = cloudamqp_instance.instance_plugin.id
			name = "rabbitmq_web_mqtt"
			enabled = true
		}
		`)
}

// func testAccNotificationConfig_Update() string {
// 	log.Printf("[DEBUG] resource_notification::testAccNotificationConfig_Update")
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
// 			type = "webhook"
// 			value = "http://example.com/webhook"
// 		}
// 		`)
// }

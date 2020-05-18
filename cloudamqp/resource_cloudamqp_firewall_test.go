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
	if testing.Short() {
		t.Skip("Skipping TestAccSecurityFirewall_Basic, since test is in short mode")
	}

	instance_name := "cloudamqp_instance.instance"
	resource_name := "cloudamqp_security_firewall.firewall"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSecurityFirewallDestroy(instance_name, resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityFirewallConfig_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityFirewallExists(instance_name),
					resource.TestCheckResourceAttr(resource_name, "rules.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "rules.0.ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resource_name, "rules.0.ports.#", "0"),
					resource.TestCheckResourceAttr(resource_name, "rules.0.services.#", "6"),
				),
			},
			{
				Config: testAccSecurityFirewallConfig_Update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityFirewallExists(instance_name),
					resource.TestCheckResourceAttr(resource_name, "rules.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "rules.0.ip", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(resource_name, "rules.0.ports.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "rules.0.ports.0", "4567"),
					resource.TestCheckResourceAttr(resource_name, "rules.0.services.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "rules.0.services.0", "AMQPS"),
				),
			},
		},
	})
}

func testAccCheckSecurityFirewallExists(instance_name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set for instance")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadFirewallSettings(instance_id)
		log.Printf("[DEBUG] resource_plugin::testAccCheckPluginEnabled data: %v", data)
		if err != nil {
			return fmt.Errorf("Error fetching item %s", err)
		}
		if data != nil {
			return fmt.Errorf("Error security firewall doesn't exists")
		}
		return nil
	}
}

func testAccSecurityFirewallDestroy(instance_name, resource_name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set for instance")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadFirewallSettings(instance_id)
		if data != nil || err == nil {
			return fmt.Errorf("Firewall still exists")
		}
		return nil
	}
}

func testAccSecurityFirewallConfig_Basic() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-security-firewall-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		resource "cloudamqp_instance" "security_firewall" {
			instance_id = cloudamqp_instance.instance.id
			rules {
				ip = "0.0.0.0/24"
				ports = []
				services = ["STOMP", "AMQP", "MQTTS", "STOMPS", "MQTT", "AMQPS"]
			}
		}
		`)
}

func testAccSecurityFirewallConfig_Update() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-security-firewall-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		resource "cloudamqp_instance" "security_firewall" {
			instance_id = cloudamqp_instance.instance.id
			rules {
				ip = "192.168.0.0/24"
				ports = [4567]
				services = ["AMQPS"]
			}
		}
		`)
}

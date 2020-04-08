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

func TestAccPlugin_Basic(t *testing.T) {
	instance_name := "cloudamqp_instance.instance"
	resource_name := "cloudamqp_plugin.mqtt_plugin"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPluginDisable(instance_name, resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccPluginConfig_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPluginEnabled(instance_name, resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", "rabbitmq_web_mqtt"),
					resource.TestCheckResourceAttr(resource_name, "enabled", "true"),
				),
			},
			{
				Config: testAccPluginConfig_Update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPluginDisable(instance_name, resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", "rabbitmq_web_mqtt"),
					resource.TestCheckResourceAttr(resource_name, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckPluginEnabled(instance_name, resource_name string) resource.TestCheckFunc {
	log.Printf("[DEBUG] resource::plugin::testAccCheckPluginEnabled resource: %s", resource_name)

	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resource_name)
		}
		if rs.Primary.Attributes["name"] == "" {
			return fmt.Errorf("No id is set for plugin resource")
		}
		plugin_name := rs.Primary.Attributes["name"]

		rs, ok = state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for instance resource")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadPlugin(instance_id, plugin_name)
		if err != nil {
			return fmt.Errorf("Error fetching item with resource %s. %s", resource_name, err)
		}
		if data["enabled"] == false {
			return fmt.Errorf("Error resource: %s not enabled", resource_name)
		}
		return nil
	}
}

func testAccCheckPluginDisable(instance_name, resource_name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Printf("[DEBUG] resource::plugins::testAccCheckPluginDisable")
		api := testAccProvider.Meta().(*api.API)

		rs, ok := state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource %s not found", resource_name)
		}
		if rs.Primary.Attributes["name"] == "" {
			return fmt.Errorf("No name is set for plugin resource")
		}
		plugin_name := rs.Primary.Attributes["name"]

		rs, ok = state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for instance resource")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		data, err := api.ReadPlugin(instance_id, plugin_name)
		if err != nil {
			return fmt.Errorf("Failed to retrieve plugin %v", err)
		}
		if data["enabled"] == true {
			return fmt.Errorf("Error resource: %s not disabled", resource_name)
		}
		return nil
	}
}

func testAccPluginConfig_Basic() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-plugin-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
			vpc_subnet = "192.168.0.1/24"
		}

		resource "cloudamqp_plugin" "mqtt_plugin" {
			instance_id = cloudamqp_instance.instance.id
			name = "rabbitmq_web_mqtt"
			enabled = true
		}
		`)
}

func testAccPluginConfig_Update() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-plugin-test"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
			vpc_subnet = "192.168.0.1/24"
		}

		resource "cloudamqp_plugin" "mqtt_plugin" {
			instance_id = cloudamqp_instance.instance.id
			name = "rabbitmq_web_mqtt"
			enabled = false
		}
		`)
}

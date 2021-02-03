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
	instanceName := "cloudamqp_instance.instance"
	resourceName := "cloudamqp_plugin.mqtt_plugin"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPluginConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPluginEnabled(instanceName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "rabbitmq_web_mqtt"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccPluginConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPluginDisable(instanceName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "rabbitmq_web_mqtt"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckPluginEnabled(instanceName, resourceName string) resource.TestCheckFunc {
	log.Printf("[DEBUG] resource::plugin::testAccCheckPluginExists resource: %s", resourceName)

	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resourceName)
		}
		if rs.Primary.Attributes["name"] == "" {
			return fmt.Errorf("No id is set for plugin resource")
		}
		pluginName := rs.Primary.Attributes["name"]

		rs, ok = state.RootModule().Resources[instanceName]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for instance resource")
		}
		instanceID, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadPlugin(instanceID, pluginName)
		if err != nil {
			return fmt.Errorf("Error fetching item with resource %s. %s", resourceName, err)
		}
		if data["enabled"] == false {
			return fmt.Errorf("Error resource: %s not enabled", resourceName)
		}
		return nil
	}
}

func testAccCheckPluginDisable(instanceName, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		log.Printf("[DEBUG] resource::plugins::testAccCheckPluginDisable")
		api := testAccProvider.Meta().(*api.API)

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", resourceName)
		}
		if rs.Primary.Attributes["name"] == "" {
			return fmt.Errorf("No name is set for plugin resource")
		}
		pluginName := rs.Primary.Attributes["name"]

		rs, ok = state.RootModule().Resources[instanceName]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for instance resource")
		}
		instanceID, _ := strconv.Atoi(rs.Primary.ID)

		data, err := api.ReadPlugin(instanceID, pluginName)
		if err != nil {
			return fmt.Errorf("Failed to retrieve plugin %v", err)
		}
		if data["enabled"] == true {
			return fmt.Errorf("Error resource: %s not disabled", resourceName)
		}
		return nil
	}
}

func testAccPluginConfigBasic() string {
	return `
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-plugin-test"
			nodes 			= 1
			plan  			= "bunny-1"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		resource "cloudamqp_plugin" "mqtt_plugin" {
			instance_id = cloudamqp_instance.instance.id
			name = "rabbitmq_web_mqtt"
			enabled = true
		}
		`
}

func testAccPluginConfigUpdate() string {
	return `
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-plugin-test"
			nodes 			= 1
			plan  			= "bunny-1"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		resource "cloudamqp_plugin" "mqtt_plugin" {
			instance_id = cloudamqp_instance.instance.id
			name = "rabbitmq_web_mqtt"
			enabled = false
		}
		`
}

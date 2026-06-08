package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccPluginsDatasource_Basic: filter plugins using enabled, required and recommended.
func TestAccPluginsDatasource_Basic(t *testing.T) {
	t.Parallel()

	dataSourceName := "data.cloudamqp_plugins.plugins"

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
          resource "cloudamqp_instance" "instance" {
            name   = "TestAccPluginsDatasource_Basic"
            plan   = "bunny-1"
            region = "amazon-web-services::eu-central-1"
          }

          data "cloudamqp_plugins" "plugins" {
            instance_id = cloudamqp_instance.instance.id
            enabled     = true
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPluginsContains(dataSourceName, "rabbitmq_consistent_hash_exchange"),
					testAccCheckPluginsContains(dataSourceName, "rabbitmq_exchange_federation"),
				),
			},
			{
				Config: `
          resource "cloudamqp_instance" "instance" {
            name   = "TestAccPluginsDatasource_Basic"
            plan   = "bunny-1"
            region = "amazon-web-services::eu-central-1"
          }

          data "cloudamqp_plugins" "plugins" {
            instance_id = cloudamqp_instance.instance.id
            enabled     = true
            required    = true
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPluginsContains(dataSourceName, "rabbitmq_management"),
					testAccCheckPluginsContains(dataSourceName, "rabbitmq_management_agent"),
					testAccCheckPluginsContains(dataSourceName, "rabbitmq_prometheus"),
					testAccCheckPluginsContains(dataSourceName, "rabbitmq_web_dispatch"),
				),
			},
			{
				Config: `
          resource "cloudamqp_instance" "instance" {
            name   = "TestAccPluginsDatasource_Basic"
            plan   = "bunny-1"
            region = "amazon-web-services::eu-central-1"
          }

          data "cloudamqp_plugins" "plugins" {
            instance_id = cloudamqp_instance.instance.id
            recommended = true
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPluginsContains(dataSourceName, "rabbitmq_shovel"),
					testAccCheckPluginsContains(dataSourceName, "rabbitmq_shovel_management"),
				),
			},
		},
	})
}

// testAccCheckPluginsContains verifies that a plugin with the given name exists in the plugins list.
func testAccCheckPluginsContains(dataSourceName, pluginName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("data source not found: %s", dataSourceName)
		}

		attrs := rs.Primary.Attributes
		countStr, ok := attrs["plugins.#"]
		if !ok {
			return fmt.Errorf("plugins.# not found in state")
		}

		count := 0
		fmt.Sscanf(countStr, "%d", &count)

		for i := 0; i < count; i++ {
			if attrs[fmt.Sprintf("plugins.%d.name", i)] == pluginName {
				return nil
			}
		}

		return fmt.Errorf("plugin %q not found in %s", pluginName, dataSourceName)
	}
}

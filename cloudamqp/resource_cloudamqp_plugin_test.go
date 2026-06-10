package cloudamqp

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccPlugin_Basic: Enabled plugin, import and disable it.
func TestAccPlugin_Basic(t *testing.T) {
	t.Parallel()

	var (
		fileNames             = []string{"instance", "plugin", "data_source/plugins"}
		instanceResourceName  = "cloudamqp_instance.instance"
		pluginResourceName    = "cloudamqp_plugin.rabbitmq_mqtt"
		dataSourcePluginsName = "data.cloudamqp_plugins.plugins"

		params = map[string]string{
			"InstanceName":  "TestAccPlugin_Basic",
			"InstanceID":    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":  "bunny-1",
			"PluginName":    "rabbitmq_mqtt",
			"PluginEnabled": "true",
		}

		paramsUpdated = map[string]string{
			"InstanceName":  "TestAccPlugin_Basic",
			"InstanceID":    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":  "bunny-1",
			"PluginName":    "rabbitmq_mqtt",
			"PluginEnabled": "false",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(pluginResourceName, "name", "rabbitmq_mqtt"),
					resource.TestCheckResourceAttr(pluginResourceName, "enabled", "true"),
					resource.TestMatchResourceAttr(dataSourcePluginsName, "plugins.#", regexp.MustCompile(`[0-9]`)),
					resource.TestCheckResourceAttr(dataSourcePluginsName, "timeout", "1800"),
				),
			},
			{
				ResourceName:            pluginResourceName,
				ImportStateIdFunc:       testAccImportCombinedStateIdFunc(instanceResourceName, pluginResourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sleep"},
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(pluginResourceName, "name", "rabbitmq_mqtt"),
					resource.TestCheckResourceAttr(pluginResourceName, "enabled", "false"),
					resource.TestMatchResourceAttr(dataSourcePluginsName, "plugins.#", regexp.MustCompile(`[0-9]`)),
					resource.TestCheckResourceAttr(dataSourcePluginsName, "timeout", "1800"),
				),
			},
		},
	})
}

// TestAccPlugin_Error: Try to enable an invalid plugin, expect job failure error.
func TestAccPlugin_Error(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccPlugin_Error"
						plan   = "bunny-1"
						region = "amazon-web-services::eu-central-1"
						tags   = ["vcr-test"]
					}

					resource "cloudamqp_plugin" "rabbitmq_invalid" {
						instance_id = cloudamqp_instance.instance.id
						name        = "rabbitmq_invalid"
						enabled     = true
					}
				`,
				ExpectError: regexp.MustCompile(`plugins_not_found`),
			},
		},
	})
}

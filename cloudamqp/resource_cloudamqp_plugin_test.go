package cloudamqp

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Recording the test benefits from having sleep set to default (10 s).
// While replaying the test, using lower sleep (1 s), will speed up the test.

// TestAccPlugin_Basic: Enabled plugin, import and disable it.
func TestAccPlugin_Basic(t *testing.T) {
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
			"PluginSleep":   "1",
		}

		paramsUpdated = map[string]string{
			"InstanceName":  "TestAccPlugin_Basic",
			"InstanceID":    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":  "bunny-1",
			"PluginName":    "rabbitmq_mqtt",
			"PluginEnabled": "false",
			"PluginSleep":   "1",
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

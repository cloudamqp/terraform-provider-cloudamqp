package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TODO: Wait until we fully support 3.13 with community plugin.
// Otherwise we earlier loaded RC/beta version if present when letting backend choose version.
// (only for emails ending with 84codes). Now we have enabled 3.13.0 since it's out.
// But no community plugin added yet.

// TestAccPluginCommunity_Basic: Install community plugin and check then disable it.
func TestAccPluginCommunity_Basic(t *testing.T) {
	var (
		fileNames    = []string{"instance", "plugin_community"}
		instanceName = "cloudamqp_instance.instance"
		resourceName = "cloudamqp_plugin_community.rabbitmq_delayed_message_exchange"

		params = map[string]string{
			"InstanceName":           "TestAccPluginCommunity_Basic",
			"InstanceID":             fmt.Sprintf("%s.id", instanceName),
			"PluginCommunityName":    "rabbitmq_delayed_message_exchange",
			"PluginCommunityEnabled": "true",
		}

		paramsUpdated = map[string]string{
			"InstanceName":           "TestAccPluginCommunity_Basic",
			"InstanceID":             fmt.Sprintf("%s.id", instanceName),
			"PluginCommunityName":    "rabbitmq_delayed_message_exchange",
			"PluginCommunityEnabled": "false",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(resourceName, "name", params["PluginCommunityName"]),
					resource.TestCheckResourceAttr(resourceName, "enabled", params["PluginCommunityEnabled"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", paramsUpdated["PluginCommunityName"]),
					resource.TestCheckResourceAttr(resourceName, "enabled", paramsUpdated["PluginCommunityEnabled"]),
				),
			},
		},
	})
}

package cloudamqp

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccPluginCommunity_Basic: Install community plugin and check then disable it.
func TestAccPluginCommunity_Basic(t *testing.T) {
	t.Parallel()

	var (
		fileNames                      = []string{"instance", "plugin_community", "data_source/plugins_community"}
		instanceResourceName           = "cloudamqp_instance.instance"
		communityPluginResourceName    = "cloudamqp_plugin_community.rabbitmq_delayed_message_exchange"
		dataSourceCommunityPluginsName = "data.cloudamqp_plugins_community.community_plugins"

		params = map[string]string{
			"InstanceName":           "TestAccPluginCommunity_Basic",
			"InstanceID":             fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":           "bunny-1",
			"PluginCommunityName":    "rabbitmq_delayed_message_exchange",
			"PluginCommunityEnabled": "true",
		}

		paramsUpdated = map[string]string{
			"InstanceName":           "TestAccPluginCommunity_Basic",
			"InstanceID":             fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":           "bunny-1",
			"PluginCommunityName":    "rabbitmq_delayed_message_exchange",
			"PluginCommunityEnabled": "false",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(communityPluginResourceName, "name", params["PluginCommunityName"]),
					resource.TestCheckResourceAttr(communityPluginResourceName, "enabled", params["PluginCommunityEnabled"]),
					resource.TestMatchResourceAttr(dataSourceCommunityPluginsName, "plugins.#", regexp.MustCompile(`[0-9]`)),
					resource.TestCheckResourceAttr(dataSourceCommunityPluginsName, "timeout", "1800"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(communityPluginResourceName, "name", paramsUpdated["PluginCommunityName"]),
					resource.TestCheckResourceAttr(communityPluginResourceName, "enabled", paramsUpdated["PluginCommunityEnabled"]),
					resource.TestMatchResourceAttr(dataSourceCommunityPluginsName, "plugins.#", regexp.MustCompile(`[0-9]`)),
					resource.TestCheckResourceAttr(dataSourceCommunityPluginsName, "timeout", "1800"),
				),
			},
		},
	})
}

// TestAccPluginCommunity_Error: Try to install an invalid community plugin, expect job failure error.
func TestAccPluginCommunity_Error(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccPluginCommunity_Error"
						plan   = "bunny-1"
						region = "amazon-web-services::eu-central-1"
						tags   = ["vcr-test"]
					}

					resource "cloudamqp_plugin_community" "rabbitmq_invalid_exchange" {
						instance_id = cloudamqp_instance.instance.id
						name        = "rabbitmq_invalid_exchange"
						enabled     = true
					}
				`,
				ExpectError: regexp.MustCompile(`is not available as a community plugin`),
			},
		},
	})
}

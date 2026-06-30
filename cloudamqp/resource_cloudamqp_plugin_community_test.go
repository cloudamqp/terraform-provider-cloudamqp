package cloudamqp

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Recording the test benefits from having sleep set to default (10 s).
// While replaying the test, using lower sleep (1 s), will speed up the test.

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
			"PluginCommunitySleep":   "1",
		}

		paramsUpdated = map[string]string{
			"InstanceName":           "TestAccPluginCommunity_Basic",
			"InstanceID":             fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":           "bunny-1",
			"PluginCommunityName":    "rabbitmq_delayed_message_exchange",
			"PluginCommunityEnabled": "false",
			"PluginCommunitySleep":   "1",
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

// TestAccPluginCommunity_Import: Importing a community plugin.
func TestAccPluginCommunity_Import(t *testing.T) {
	t.Parallel()

	var (
		instanceResourceName        = "cloudamqp_instance.instance"
		communityPluginResourceName = "cloudamqp_plugin_community.rabbitmq_delayed_message_exchange"
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
				  resource "cloudamqp_instance" "instance" {
					  name   = "TestAccPluginCommunity_Import"
					  plan   = "bunny-1"
					  region = "amazon-web-services::eu-central-1"
					  tags   = ["vcr-test"]
					}

          resource "cloudamqp_plugin_community" "rabbitmq_delayed_message_exchange" {
            instance_id = cloudamqp_instance.instance.id
            name        = "rabbitmq_delayed_message_exchange"
            enabled     = true
						sleep       = 1
          }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(communityPluginResourceName, "name", "rabbitmq_delayed_message_exchange"),
					resource.TestCheckResourceAttr(communityPluginResourceName, "enabled", "true"),
				),
			},
			{
				ResourceName: communityPluginResourceName,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					rs, ok := state.RootModule().Resources[instanceResourceName]
					if !ok {
						return "", fmt.Errorf("Resource %s not found", instanceResourceName)
					}
					return fmt.Sprintf("%s,%s", "rabbitmq_delayed_message_exchange", rs.Primary.ID), nil
				},
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"require", "sleep"},
			},
		},
	})
}

package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Playing the test benefits from having sleep set to default (10 s) while recording.
// While using lower sleep (1 s) for replaying the test. This will speed up the test.

// TestAccPluginCommunity_Basic: Install community plugin and check then disable it.
func TestAccPluginCommunity_Basic(t *testing.T) {
	var (
		fileNames                   = []string{"instance", "plugin_community"}
		instanceResourceName        = "cloudamqp_instance.instance"
		communityPluginResourceName = "cloudamqp_plugin_community.rabbitmq_delayed_message_exchange"

		params = map[string]string{
			"InstanceName":           "TestAccPluginCommunity_Basic",
			"InstanceID":             fmt.Sprintf("%s.id", instanceResourceName),
			"PluginCommunityName":    "rabbitmq_delayed_message_exchange",
			"PluginCommunityEnabled": "true",
			"PluginCommunitySleep":   "1",
		}

		paramsUpdated = map[string]string{
			"InstanceName":           "TestAccPluginCommunity_Basic",
			"InstanceID":             fmt.Sprintf("%s.id", instanceResourceName),
			"PluginCommunityName":    "rabbitmq_delayed_message_exchange",
			"PluginCommunityEnabled": "false",
			"PluginCommunitySleep":   "1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(communityPluginResourceName, "name", params["PluginCommunityName"]),
					resource.TestCheckResourceAttr(communityPluginResourceName, "enabled", params["PluginCommunityEnabled"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(communityPluginResourceName, "name", paramsUpdated["PluginCommunityName"]),
					resource.TestCheckResourceAttr(communityPluginResourceName, "enabled", paramsUpdated["PluginCommunityEnabled"]),
				),
			},
		},
	})
}

package cloudamqp

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Unit tests for cloudamqp_plugin_community against the fake API in
// resource_cloudamqp_plugin_unit_test.go. An installed community plugin shows up in the normal
// plugin list, and uninstalling removes it from there.

const pluginCommunityResourceName = "cloudamqp_plugin_community.delayed"

var (
	pluginCommunityInstallRequest = fmt.Sprintf(`POST /api/instances/%d/plugins/community {"plugin_name":"rabbitmq_delayed_message_exchange"}`, pluginTestInstanceID)
	pluginCommunityDisableRequest = fmt.Sprintf(`PUT /api/instances/%d/plugins/community {"enabled":false,"plugin_name":"rabbitmq_delayed_message_exchange"}`, pluginTestInstanceID)
	pluginCommunityDeleteRequest  = fmt.Sprintf(`DELETE /api/instances/%d/plugins/community/rabbitmq_delayed_message_exchange`, pluginTestInstanceID)
)

func newFakePluginAPIWithCommunity(t *testing.T) (*fakePluginAPI, *httptest.Server) {
	fake, server := newFakePluginAPI(t)
	fake.plugins = []*fakePlugin{
		{Name: "rabbitmq_management", Version: "4.0.2", Description: "RabbitMQ Management Console", Enabled: true},
	}
	fake.community = []map[string]any{
		{"name": "rabbitmq_delayed_message_exchange",
			"description": "A plugin that adds delayed-messaging (or scheduled-messaging) to RabbitMQ.",
			"require":     "3.4.0"},
	}
	return fake, server
}

func pluginCommunityResourceConfig(baseURL, enabled string) string {
	return pluginUnitTestConfig(baseURL, fmt.Sprintf(`
resource "cloudamqp_plugin_community" "delayed" {
  instance_id = %d
  name        = "rabbitmq_delayed_message_exchange"
  enabled     = %s
  sleep       = 1
}
`, pluginTestInstanceID, enabled))
}

func (f *fakePluginAPI) checkInstalled(name string, installed bool) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		f.mu.Lock()
		defer f.mu.Unlock()
		if (f.plugin(name) != nil) != installed {
			return fmt.Errorf("expected plugin %s installed=%t in the API", name, installed)
		}
		return nil
	}
}

// TestPluginCommunity_InstallUninstall: create installs the plugin and waits for it, changing
// enabled updates it in place, and destroy uninstalls it and waits until it is gone.
func TestPluginCommunity_InstallUninstall(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPIWithCommunity(t)

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config: pluginCommunityResourceConfig(server.URL, "true"),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(pluginCommunityInstallRequest),
				fake.checkEnabled("rabbitmq_delayed_message_exchange", true),
				resource.TestCheckResourceAttr(pluginCommunityResourceName, "id", "rabbitmq_delayed_message_exchange"),
				resource.TestCheckResourceAttr(pluginCommunityResourceName, "instance_id", fmt.Sprint(pluginTestInstanceID)),
				resource.TestCheckResourceAttr(pluginCommunityResourceName, "name", "rabbitmq_delayed_message_exchange"),
				resource.TestCheckResourceAttr(pluginCommunityResourceName, "enabled", "true"),
				resource.TestCheckResourceAttr(pluginCommunityResourceName, "require", "3.4.0"),
				resource.TestCheckResourceAttr(pluginCommunityResourceName, "description",
					"A plugin that adds delayed-messaging (or scheduled-messaging) to RabbitMQ."),
				resource.TestCheckResourceAttr(pluginCommunityResourceName, "sleep", "1"),
				resource.TestCheckResourceAttr(pluginCommunityResourceName, "timeout", "1800"),
			),
		},
		{
			Config:   pluginCommunityResourceConfig(server.URL, "true"),
			PlanOnly: true,
		},
		{
			Config: pluginCommunityResourceConfig(server.URL, "false"),
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(pluginCommunityInstallRequest, pluginCommunityDisableRequest),
				fake.checkEnabled("rabbitmq_delayed_message_exchange", false),
				resource.TestCheckResourceAttr(pluginCommunityResourceName, "enabled", "false"),
				resource.TestCheckResourceAttr(pluginCommunityResourceName, "require", "3.4.0"),
			),
		},
		{
			Config:   pluginCommunityResourceConfig(server.URL, "false"),
			PlanOnly: true,
		},
	}, resource.ComposeAggregateTestCheckFunc(
		fake.checkRequests(pluginCommunityInstallRequest, pluginCommunityDisableRequest, pluginCommunityDeleteRequest),
		fake.checkInstalled("rabbitmq_delayed_message_exchange", false),
	))
}

// TestPluginCommunity_CreateDisabled: create with enabled = false installs the plugin and
// disables it in the same apply.
func TestPluginCommunity_CreateDisabled(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPIWithCommunity(t)
	config := pluginCommunityResourceConfig(server.URL, "false")

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeAggregateTestCheckFunc(
				fake.checkRequests(pluginCommunityInstallRequest, pluginCommunityDisableRequest),
				fake.checkEnabled("rabbitmq_delayed_message_exchange", false),
				resource.TestCheckResourceAttr(pluginCommunityResourceName, "enabled", "false"),
				resource.TestCheckResourceAttr(pluginCommunityResourceName, "require", "3.4.0"),
			),
		},
		{
			Config:   config,
			PlanOnly: true,
		},
	}, nil)
}

// TestPluginCommunity_Import: import with "{name},{instance_id}" reads the installed plugin.
func TestPluginCommunity_Import(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPIWithCommunity(t)

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config: pluginCommunityResourceConfig(server.URL, "true"),
		},
		{
			ResourceName:      pluginCommunityResourceName,
			ImportState:       true,
			ImportStateId:     fmt.Sprintf("rabbitmq_delayed_message_exchange,%d", pluginTestInstanceID),
			ImportStateVerify: true,
			// Import reads the installed plugin list, which has no required RabbitMQ version
			ImportStateVerifyIgnore: []string{"sleep", "timeout", "timeouts", "require"},
		},
	}, nil)
}

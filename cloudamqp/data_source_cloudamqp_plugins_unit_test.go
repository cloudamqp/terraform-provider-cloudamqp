package cloudamqp

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Unit tests for the cloudamqp_plugins data source against the fake API in
// resource_cloudamqp_plugin_unit_test.go. The API sends recommended and required as a reason
// text or null, the data source stores them as booleans.

const pluginsDataSourceName = "data.cloudamqp_plugins.plugins"

func newFakePluginAPIWithList(t *testing.T) (*fakePluginAPI, *httptest.Server) {
	fake, server := newFakePluginAPI(t)
	required := "Used for health and metrics"
	recommended := "Used for instance logs"
	fake.plugins = []*fakePlugin{
		{Name: "rabbitmq_management", Version: "4.0.2", Description: "RabbitMQ Management Console",
			Enabled: true, Required: &required},
		{Name: "rabbitmq_mqtt", Version: "4.0.2", Description: "RabbitMQ MQTT Adapter"},
		{Name: "rabbitmq_event_exchange", Version: "4.0.2", Description: "Event Exchange Type",
			Recommended: &recommended},
		{Name: "rabbitmq_shovel", Version: "4.0.2", Description: "Data Shovel for RabbitMQ", Enabled: true},
	}
	return fake, server
}

func pluginsDataSourceConfig(baseURL, filter string) string {
	return pluginUnitTestConfig(baseURL, fmt.Sprintf(`
data "cloudamqp_plugins" "plugins" {
  instance_id = %d
  %s
}
`, pluginTestInstanceID, filter))
}

// TestPluginsDataSource_All: all plugins are listed in API order, with recommended and required
// as booleans.
func TestPluginsDataSource_All(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPIWithList(t)

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config: pluginsDataSourceConfig(server.URL, ""),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(pluginsDataSourceName, "id", fmt.Sprintf("%d.plugins", pluginTestInstanceID)),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "sleep", "10"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "timeout", "1800"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.#", "4"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.0.name", "rabbitmq_management"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.0.version", "4.0.2"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.0.description", "RabbitMQ Management Console"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.0.enabled", "true"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.0.required", "true"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.0.recommended", "false"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.1.name", "rabbitmq_mqtt"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.1.enabled", "false"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.1.required", "false"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.1.recommended", "false"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.2.name", "rabbitmq_event_exchange"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.2.recommended", "true"),
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.3.name", "rabbitmq_shovel"),
			),
		},
	}, nil)
}

// TestPluginsDataSource_Filters: each filter keeps only the matching plugins.
func TestPluginsDataSource_Filters(t *testing.T) {
	t.Parallel()

	for filter, names := range map[string][]string{
		"enabled = true":     {"rabbitmq_management", "rabbitmq_shovel"},
		"required = true":    {"rabbitmq_management"},
		"recommended = true": {"rabbitmq_event_exchange"},
		"enabled = false":    {"rabbitmq_management", "rabbitmq_mqtt", "rabbitmq_event_exchange", "rabbitmq_shovel"},
	} {
		t.Run(filter, func(t *testing.T) {
			t.Parallel()
			fake, server := newFakePluginAPIWithList(t)

			checks := []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(pluginsDataSourceName, "plugins.#", fmt.Sprint(len(names))),
			}
			for i, name := range names {
				checks = append(checks, resource.TestCheckResourceAttr(pluginsDataSourceName, fmt.Sprintf("plugins.%d.name", i), name))
			}

			pluginUnitTest(t, fake, []resource.TestStep{
				{
					Config: pluginsDataSourceConfig(server.URL, filter),
					Check:  resource.ComposeAggregateTestCheckFunc(checks...),
				},
			}, nil)
		})
	}
}

const pluginsCommunityDataSourceName = "data.cloudamqp_plugins_community.community"

// TestPluginsCommunityDataSource_Read: the available community plugins are listed in API order.
func TestPluginsCommunityDataSource_Read(t *testing.T) {
	t.Parallel()
	fake, server := newFakePluginAPI(t)
	fake.community = []map[string]any{
		{"name": "rabbitmq_delayed_message_exchange", "description": "Delayed messaging", "require": "3.4.0"},
		{"name": "rabbitmq_message_deduplication", "description": "Message deduplication", "require": "3.13.0"},
	}

	pluginUnitTest(t, fake, []resource.TestStep{
		{
			Config: pluginUnitTestConfig(server.URL, fmt.Sprintf(`
data "cloudamqp_plugins_community" "community" {
  instance_id = %d
}
`, pluginTestInstanceID)),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(pluginsCommunityDataSourceName, "id", fmt.Sprintf("%d.plugins_community", pluginTestInstanceID)),
				resource.TestCheckResourceAttr(pluginsCommunityDataSourceName, "sleep", "10"),
				resource.TestCheckResourceAttr(pluginsCommunityDataSourceName, "timeout", "1800"),
				resource.TestCheckResourceAttr(pluginsCommunityDataSourceName, "plugins.#", "2"),
				resource.TestCheckResourceAttr(pluginsCommunityDataSourceName, "plugins.0.name", "rabbitmq_delayed_message_exchange"),
				resource.TestCheckResourceAttr(pluginsCommunityDataSourceName, "plugins.0.description", "Delayed messaging"),
				resource.TestCheckResourceAttr(pluginsCommunityDataSourceName, "plugins.0.require", "3.4.0"),
				resource.TestCheckResourceAttr(pluginsCommunityDataSourceName, "plugins.1.name", "rabbitmq_message_deduplication"),
				resource.TestCheckResourceAttr(pluginsCommunityDataSourceName, "plugins.1.require", "3.13.0"),
			),
		},
	}, nil)
}

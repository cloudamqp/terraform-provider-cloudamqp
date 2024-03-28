package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccPlugin_Basic: Enabled plugin, import and disable it.
func TestAccPlugin_Basic(t *testing.T) {
	var (
		fileNames    = []string{"instance", "plugin"}
		instanceName = "cloudamqp_instance.instance"
		resourceName = "cloudamqp_plugin.rabbitmq_mqtt"

		params = map[string]string{
			"InstanceName":  "TestAccPlugin_Basic",
			"InstanceID":    fmt.Sprintf("%s.id", instanceName),
			"PluginName":    "rabbitmq_mqtt",
			"PluginEnabled": "true",
		}

		paramsUpdated = map[string]string{
			"InstanceName":  "TestAccPlugin_Basic",
			"InstanceID":    fmt.Sprintf("%s.id", instanceName),
			"PluginName":    "rabbitmq_mqtt",
			"PluginEnabled": "false",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(resourceName, "name", "rabbitmq_mqtt"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceName, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "rabbitmq_mqtt"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

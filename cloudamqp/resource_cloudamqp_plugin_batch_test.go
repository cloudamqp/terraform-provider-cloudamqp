package cloudamqp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccPluginBatch_Basic: Create instance with three enabled plugins, then add a plugin and
// disable another and finally disable all managed plugins.
func TestAccPluginBatch_Basic(t *testing.T) {
	t.Parallel()

	var (
		instanceResourceName    = "cloudamqp_instance.instance"
		pluginBatchResourceName = "cloudamqp_plugin_batch.plugins"
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccPluginBatch_Basic"
						plan   = "bunny-1"
						region = "amazon-web-services::eu-central-1"
						tags   = ["vcr-test"]
					}

					resource "cloudamqp_plugin_batch" "plugins" {
						instance_id = cloudamqp_instance.instance.id

						plugins = {
							rabbitmq_stomp           = true,
							rabbitmq_top             = true,
							rabbitmq_web_mqtt        = true
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccPluginBatch_Basic"),
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.%", "3"),
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.rabbitmq_stomp", "true"),
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.rabbitmq_top", "true"),
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.rabbitmq_web_mqtt", "true"),
				),
			},
			{
				Config: `
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccPluginBatch_Basic"
						plan   = "bunny-1"
						region = "amazon-web-services::eu-central-1"
						tags   = ["vcr-test"]
					}

					resource "cloudamqp_plugin_batch" "plugins" {
						instance_id = cloudamqp_instance.instance.id

						plugins = {
							rabbitmq_stomp           = true,
							rabbitmq_top             = false,
							rabbitmq_web_mqtt        = true,
							rabbitmq_random_exchange = true
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.%", "4"),
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.rabbitmq_stomp", "true"),
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.rabbitmq_top", "false"),
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.rabbitmq_web_mqtt", "true"),
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.rabbitmq_random_exchange", "true"),
				),
			},
			{
				Config: `
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccPluginBatch_Basic"
						plan   = "bunny-1"
						region = "amazon-web-services::eu-central-1"
						tags   = ["vcr-test"]
					}

					resource "cloudamqp_plugin_batch" "plugins" {
						instance_id = cloudamqp_instance.instance.id

						plugins = {
							rabbitmq_stomp           = false,
							rabbitmq_top             = false,
							rabbitmq_web_mqtt        = false,
							rabbitmq_random_exchange = false
						}
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.%", "4"),
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.rabbitmq_stomp", "false"),
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.rabbitmq_top", "false"),
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.rabbitmq_web_mqtt", "false"),
					resource.TestCheckResourceAttr(pluginBatchResourceName, "plugins.rabbitmq_random_exchange", "false"),
				),
			},
		},
	})
}

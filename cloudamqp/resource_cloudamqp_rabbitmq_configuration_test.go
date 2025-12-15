package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccRabbitMqConfiguration_Basic: Update RabbitMQ configuration and import.
func TestAccRabbitMqConfiguration_Basic(t *testing.T) {
	t.Parallel()

	var (
		fileNames                  = []string{"instance", "rabbitmq_configuration/config"}
		instanceResourceName       = "cloudamqp_instance.instance"
		rabbitmqConfigResourceName = "cloudamqp_rabbitmq_configuration.rabbitmq_config"

		params = map[string]string{
			"InstanceName":    "TestAccRabbitMqConfiguration_Basic",
			"InstanceID":      fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":    "bunny-1",
			"ChannelMax":      "100",
			"ConnectionMax":   "100",
			"ConsumerTimeout": "720000",
			"Heartbeat":       "60",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(rabbitmqConfigResourceName, "channel_max", params["ChannelMax"]),
					resource.TestCheckResourceAttr(rabbitmqConfigResourceName, "connection_max", params["ConnectionMax"]),
					resource.TestCheckResourceAttr(rabbitmqConfigResourceName, "consumer_timeout", params["ConsumerTimeout"]),
					resource.TestCheckResourceAttr(rabbitmqConfigResourceName, "heartbeat", params["Heartbeat"]),
				),
			},
			{
				ResourceName:            rabbitmqConfigResourceName,
				ImportStateIdFunc:       testAccImportStateIdFunc(instanceResourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sleep", "timeout"},
			},
		},
	})
}

func TestAccRabbitMqConfiguration_LogExhangeLevel(t *testing.T) {
	t.Parallel()

	var (
		fileNames                  = []string{"instance", "rabbitmq_configuration/config", "data_source/nodes", "node_actions"}
		instanceResourceName       = "cloudamqp_instance.instance"
		rabbitMqConfigResourceName = "cloudamqp_rabbitmq_configuration.rabbitmq_config"
		nodeActionResourceName     = "cloudamqp_node_actions.node_action"
		dataSourceNodesName        = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName":     "TestAccRabbitMqConfiguration_Basic",
			"InstanceID":       fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":     "bunny-1",
			"LogExchangeLevel": "info",
			"NodeName":         fmt.Sprintf("%s.nodes[0].name", dataSourceNodesName),
			"NodeAction":       "restart",
			"NodeDependsOn":    rabbitMqConfigResourceName,
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(rabbitMqConfigResourceName, "log_exchange_level", params["LogExchangeLevel"]),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.configured", "true"),
					resource.TestCheckResourceAttr(nodeActionResourceName, "action", params["NodeAction"]),
				),
			},
		},
	})
}

// TestAccRabbitMqConfiguration_ZeroValue: While using Framework 0 int values can be detected correctly.
// Issue in Terraform SDK v2, where 0 value cannot be detected due to default int value being 0.
func TestAccRabbitMqConfiguration_ZeroValue(t *testing.T) {
	t.Parallel()

	var (
		fileNames                  = []string{"instance", "rabbitmq_configuration/zero_value"}
		instanceResourceName       = "cloudamqp_instance.instance"
		rabbitmqConfigResourceName = "cloudamqp_rabbitmq_configuration.rabbitmq_config"

		params = map[string]string{
			"InstanceName": "TestAccRabbitMqConfiguration_ZeroValue",
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan": "bunny-1",
			"Heartbeat":    "0", // Set heartbeat to 0 to test handling of 0 value.
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(rabbitmqConfigResourceName, "heartbeat", params["Heartbeat"]),
				),
			},
		},
	})
}

func TestAccRabbitMqConfiguration_MqttConfiguration(t *testing.T) {
	t.Parallel()

	instanceResourceName := "cloudamqp_instance.instance"
	rabbitmqConfigResourceName := "cloudamqp_rabbitmq_configuration.rabbitmq_config"
	dataSourceNodesName := "data.cloudamqp_nodes.nodes"

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_instance" "instance" {
						name   = "TestAccRabbitMqConfiguration_MqttConfiguration"
						plan   = "bunny-1"
						region = "amazon-web-services::eu-central-1"
						tags   = []
					}

					resource "cloudamqp_rabbitmq_configuration" "rabbitmq_config" {
						instance_id                      = cloudamqp_instance.instance.id
						mqtt_vhost                       = cloudamqp_instance.instance.vhost
						mqtt_exchange                    = "amq.topic"
						mqtt_ssl_cert_login              = true
						ssl_options_fail_if_no_peer_cert = true
						ssl_options_verify               = "verify_peer"
					}

					data "cloudamqp_nodes" "nodes" {
						instance_id = cloudamqp_instance.instance.id
					}

					resource "cloudamqp_node_actions" "node_action" {
						instance_id = cloudamqp_instance.instance.id
						node_name   = data.cloudamqp_nodes.nodes.nodes[0].name
						action      = "restart"

						depends_on = [
							cloudamqp_rabbitmq_configuration.rabbitmq_config,
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccRabbitMqConfiguration_MqttConfiguration"),
					resource.TestCheckResourceAttr(rabbitmqConfigResourceName, "mqtt_exchange", "amq.topic"),
					resource.TestCheckResourceAttr(rabbitmqConfigResourceName, "mqtt_ssl_cert_login", "true"),
					resource.TestCheckResourceAttr(rabbitmqConfigResourceName, "ssl_options_fail_if_no_peer_cert", "true"),
					resource.TestCheckResourceAttr(rabbitmqConfigResourceName, "ssl_options_verify", "verify_peer"),
					resource.TestCheckResourceAttrPair(
						rabbitmqConfigResourceName, "mqtt_vhost",
						instanceResourceName, "vhost",
					),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.configured", "true"),
				),
			},
		},
	})
}

package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccRabbitMqConfiguration_Basic: Update RabbitMQ configuration and import.
func TestAccRabbitMqConfiguration_Basic(t *testing.T) {
	var (
		fileNames    = []string{"instance", "rabbitmq_configuration"}
		instanceName = "cloudamqp_instance.instance"
		resourceName = "cloudamqp_rabbitmq_configuration.rabbitmq_config"

		params = map[string]string{
			"InstanceName":    "TestAccRabbitMqConfiguration_Basic",
			"InstanceID":      fmt.Sprintf("%s.id", instanceName),
			"ChannelMax":      "100",
			"ConnectionMax":   "100",
			"ConsumerTimeout": "720000",
			"Heartbeat":       "60",
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
					resource.TestCheckResourceAttr(resourceName, "channel_max", params["ChannelMax"]),
					resource.TestCheckResourceAttr(resourceName, "connection_max", params["ConnectionMax"]),
					resource.TestCheckResourceAttr(resourceName, "consumer_timeout", params["ConsumerTimeout"]),
					resource.TestCheckResourceAttr(resourceName, "heartbeat", params["Heartbeat"]),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccImportStateIdFunc(instanceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRabbitMqConfiguration_LogExhangeLevel(t *testing.T) {
	var (
		fileNames           = []string{"instance", "rabbitmq_configuration", "nodes_data", "node_actions"}
		instanceName        = "cloudamqp_instance.instance"
		rabbitMqConfigName  = "cloudamqp_rabbitmq_configuration.rabbitmq_config"
		nodeActionsName     = "cloudamqp_node_actions.node_action"
		dataSourceNodesName = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName":     "TestAccRabbitMqConfiguration_Basic",
			"InstanceID":       fmt.Sprintf("%s.id", instanceName),
			"LogExchangeLevel": "info",
			"NodeName":         fmt.Sprintf("%s.nodes[0].name", dataSourceNodesName),
			"NodeAction":       "restart",
			"NodeDependsOn":    rabbitMqConfigName,
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
					resource.TestCheckResourceAttr(rabbitMqConfigName, "log_exchange_level", params["LogExchangeLevel"]),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.configured", "true"),
					resource.TestCheckResourceAttr(nodeActionsName, "action", params["NodeAction"]),
				),
			},
		},
	})
}

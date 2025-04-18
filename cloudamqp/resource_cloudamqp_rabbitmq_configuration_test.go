package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccRabbitMqConfiguration_Basic: Update RabbitMQ configuration and import.
func TestAccRabbitMqConfiguration_Basic(t *testing.T) {
	var (
		fileNames            = []string{"instance", "rabbitmq_configuration"}
		instanceResourceName = "cloudamqp_instance.instance"
		pluginResourceName   = "cloudamqp_rabbitmq_configuration.rabbitmq_config"

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
					resource.TestCheckResourceAttr(pluginResourceName, "channel_max", params["ChannelMax"]),
					resource.TestCheckResourceAttr(pluginResourceName, "connection_max", params["ConnectionMax"]),
					resource.TestCheckResourceAttr(pluginResourceName, "consumer_timeout", params["ConsumerTimeout"]),
					resource.TestCheckResourceAttr(pluginResourceName, "heartbeat", params["Heartbeat"]),
				),
			},
			{
				ResourceName:      pluginResourceName,
				ImportStateIdFunc: testAccImportStateIdFunc(instanceResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRabbitMqConfiguration_LogExhangeLevel(t *testing.T) {
	var (
		fileNames                  = []string{"instance", "rabbitmq_configuration", "data_source/nodes", "node_actions"}
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

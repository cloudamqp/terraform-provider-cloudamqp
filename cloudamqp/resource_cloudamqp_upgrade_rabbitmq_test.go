package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/converter"
)

// TestAccUpgradeRabbitMQ_Latest: Upgrade RabbitMQ to latest possible version, from 3.12.2 -> 3.13.2
// Extra checks are needed when comparing versions, because next step is executed before backend
// have been updated. Same reason unable to use cloudamqp_upgradable_versions data source correctly.
func TestAccUpgradeRabbitMQ_Latest(t *testing.T) {
	var (
		fileNames            = []string{"instance_with_version", "data_source/nodes"}
		instanceResourceName = "cloudamqp_instance.instance"
		dataSourceNodesName  = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName":       "TestAccUpgradeRabbitMQ_Latest",
			"InstanceTags":       converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":         fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":       "bunny-1",
			"InstanceRmqVersion": "3.12.2",
		}

		fileNamesUpgrade = []string{"instance", "data_source/nodes", "upgrade_rabbitmq_latest"}

		paramsUpgrade01 = map[string]string{
			"InstanceName":                  "TestAccUpgradeRabbitMQ_Latest",
			"InstanceTags":                  converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":                    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":                  "bunny-1",
			"UpgradeRabbitMQCurrentVersion": "3.12.2",
			"UpgradeRabbitMQNewVersion":     "3.12.13",
		}

		paramsUpgrade02 = map[string]string{
			"InstanceName":                  "TestAccUpgradeRabbitMQ_Latest",
			"InstanceTags":                  converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":                    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":                  "bunny-1",
			"UpgradeRabbitMQCurrentVersion": "3.12.13",
			"UpgradeRabbitMQNewVersion":     "3.13.2",
		}

		fileNamesCheckUpgrade = []string{"instance", "data_source/nodes"}
		paramsCheck           = map[string]string{
			"InstanceName": "TestAccUpgradeRabbitMQ_Latest",
			"InstanceTags": converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan": "bunny-1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(instanceResourceName, "plan", "bunny-1"),
					resource.TestCheckResourceAttr(instanceResourceName, "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr(instanceResourceName, "rmq_version", params["InstanceRmqVersion"]),
					resource.TestCheckResourceAttr(instanceResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(instanceResourceName, "tags.0", "terraform"),
					resource.TestCheckResourceAttr(instanceResourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", params["InstanceRmqVersion"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesUpgrade, paramsUpgrade01),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "3.12.2"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesCheckUpgrade, paramsCheck),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "3.12.13"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesUpgrade, paramsUpgrade02),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "3.12.13"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesCheckUpgrade, paramsCheck),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "3.13.2"),
				),
			},
		},
	})
}

// TestAccUpgradeRabbitMQ_Specific: Upgrade RabbitMQ to a specific version, from 3.12.2 -> 3.13.2
// Extra checks are needed when comparing versions, because next step is executed before backend
// have been updated.
func TestAccUpgradeRabbitMQ_Specific(t *testing.T) {
	var (
		fileNames            = []string{"instance_with_version", "data_source/nodes"}
		instanceResourceName = "cloudamqp_instance.instance"
		dataSourceNodesName  = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName":       "TestAccUpgradeRabbitMQ_Latest",
			"InstanceTags":       converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":         fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":       "bunny-1",
			"InstanceRmqVersion": "3.12.2",
		}

		fileNamesUpgrade = []string{"instance", "data_source/nodes", "upgrade_rabbitmq_latest"}

		paramsUpgrade01 = map[string]string{
			"InstanceName":              "TestAccUpgradeRabbitMQ_Latest",
			"InstanceTags":              converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":                fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":              "bunny-1",
			"UpgradeRabbitMQNewVersion": "3.12.13",
		}

		paramsUpgrade02 = map[string]string{
			"InstanceName":              "TestAccUpgradeRabbitMQ_Latest",
			"InstanceTags":              converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":                fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":              "bunny-1",
			"UpgradeRabbitMQNewVersion": "3.13.2",
		}

		fileNamesCheckUpgrade = []string{"instance", "data_source/nodes"}
		paramsCheck           = map[string]string{
			"InstanceName": "TestAccUpgradeRabbitMQ_Latest",
			"InstanceTags": converter.CommaStringArray([]string{"terraform"}),
			"InstancePlan": "bunny-1",
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(instanceResourceName, "plan", "bunny-1"),
					resource.TestCheckResourceAttr(instanceResourceName, "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr(instanceResourceName, "rmq_version", params["InstanceRmqVersion"]),
					resource.TestCheckResourceAttr(instanceResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(instanceResourceName, "tags.0", "terraform"),
					resource.TestCheckResourceAttr(instanceResourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", params["InstanceRmqVersion"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesUpgrade, paramsUpgrade01),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "3.12.2"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesCheckUpgrade, paramsCheck),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "3.12.13"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesUpgrade, paramsUpgrade02),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "3.12.13"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesCheckUpgrade, paramsCheck),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "3.13.2"),
				),
			},
		},
	})
}

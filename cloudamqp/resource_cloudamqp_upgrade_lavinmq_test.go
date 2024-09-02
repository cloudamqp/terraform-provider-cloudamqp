package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/converter"
)

// TestAccUpgradeLavinMQ_Latest: Upgrade LavinMQ to latest possible version, from 1.3.0 -> 2.0.0-rc.3
// Extra checks are needed when comparing versions, because next step is executed before backend
// have been updated. Same reason unable to use cloudamqp_upgradable_versions data source correctly.
func TestAccUpgradeLavinMQ_Latest(t *testing.T) {
	var (
		fileNames            = []string{"instance_with_version_lavinmq", "data_source/nodes"}
		instanceResourceName = "cloudamqp_instance.instance"
		dataSourceNodesName  = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName":       "TestAccUpgradeLavinMQ_Latest",
			"InstanceTags":       converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":         fmt.Sprintf("%s.id", instanceResourceName),
			"InstanceRmqVersion": "1.3.0",
		}

		fileNamesUpgrade = []string{"instance_lavinmq", "data_source/nodes", "upgrade_lavinmq"}

		paramsUpgrade01 = map[string]string{
			"InstanceName":             "TestAccUpgradeLavinMQ_Latest",
			"InstanceTags":             converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":               fmt.Sprintf("%s.id", instanceResourceName),
			"UpgradeLavinMQNewVersion": "2.0.0-rc.3",
		}

		fileNamesCheckUpgrade = []string{"instance_lavinmq", "data_source/nodes"}
		paramsCheck           = map[string]string{
			"InstanceName": "TestAccUpgradeLavinMQ_Latest",
			"InstanceTags": converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(instanceResourceName, "plan", "wolverine-1"),
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
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "1.3.0"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesCheckUpgrade, paramsCheck),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "2.0.0-rc.3"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesCheckUpgrade, paramsCheck),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "2.0.0-rc.3"),
				),
			},
		},
	})
}

// TestAccUpgradeLavintMQ_Specific: Upgrade LavinMQ to a specific version, from 1.3.0 -> 2.0.0-rc.3
// Extra checks are needed when comparing versions, because next step is executed before backend
// have been updated.
func TestAccUpgradeLavinMQ_Specific(t *testing.T) {
	var (
		fileNames            = []string{"instance_with_version_lavinmq", "data_source/nodes"}
		instanceResourceName = "cloudamqp_instance.instance"
		dataSourceNodesName  = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName":       "TestAccUpgradeLavinMQ_Latest",
			"InstanceTags":       converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":         fmt.Sprintf("%s.id", instanceResourceName),
			"InstanceRmqVersion": "1.3.0",
		}

		fileNamesUpgrade = []string{"instance_lavinmq", "data_source/nodes", "upgrade_lavinmq"}

		paramsUpgrade01 = map[string]string{
			"InstanceName":             "TestAccUpgradeLavinMQ_Latest",
			"InstanceTags":             converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":               fmt.Sprintf("%s.id", instanceResourceName),
			"UpgradeLavinMQNewVersion": "2.0.0-rc.3",
		}

		fileNamesCheckUpgrade = []string{"instance_lavinmq", "data_source/nodes"}
		paramsCheck           = map[string]string{
			"InstanceName": "TestAccUpgradeLavinMQ_Latest",
			"InstanceTags": converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(instanceResourceName, "plan", "wolverine-1"),
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
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "1.3.0"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesCheckUpgrade, paramsCheck),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "2.0.0-rc.3"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesCheckUpgrade, paramsCheck),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "2.0.0-rc.3"),
				),
			},
		},
	})
}

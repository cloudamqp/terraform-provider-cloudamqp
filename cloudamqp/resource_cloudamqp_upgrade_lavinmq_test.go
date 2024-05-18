package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/converter"
)

// TestAccUpgradeLavintMQ: Upgrade LavinMQ to a specific version, from 1.3.1 -> 2.0.0-rc.3
// Extra checks are needed when comparing versions, because next step is executed before backend
// have been updated.
func TestAccUpgradeLavinMQ(t *testing.T) {
	var (
		fileNames            = []string{"instance_with_version", "data_source/nodes"}
		instanceResourceName = "cloudamqp_instance.instance"
		dataSourceNodesName  = "data.cloudamqp_nodes.nodes"
		plan                 = "wolverine-1"

		params = map[string]string{
			"InstanceName":       "TestAccUpgradeLavinMQ",
			"InstanceTags":       converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":         fmt.Sprintf("%s.id", instanceResourceName),
			"InstanceRmqVersion": "1.3.1",
			"InstancePlan":       plan,
		}

		fileNamesUpgrade = []string{"instance", "data_source/nodes", "upgrade_lavinmq"}

		paramsUpgrade01 = map[string]string{
			"InstanceName":             "TestAccUpgradeLavinMQ",
			"InstanceTags":             converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":               fmt.Sprintf("%s.id", instanceResourceName),
			"UpgradeLavinMQNewVersion": "2.0.0-rc.3",
			"InstancePlan":             plan,
		}

		fileNamesCheckUpgrade = []string{"instance", "data_source/nodes"}
		paramsCheck           = map[string]string{
			"InstanceName": "TestAccUpgradeLavinMQ",
			"InstanceTags": converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan": plan,
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
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
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.rabbitmq_version", "1.3.1"),
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

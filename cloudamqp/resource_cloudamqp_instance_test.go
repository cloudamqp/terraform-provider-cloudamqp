package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/converter"
)

// TestAccInstance_Basic: Creating dedicated AWS instance, do some minor updates, import and read
// nodes data source.
func TestAccInstance_Basic(t *testing.T) {
	var (
		fileNames            = []string{"instance", "data_source/nodes"}
		instanceResourceName = "cloudamqp_instance.instance"
		dataSourceNodesName  = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName": "TestAccInstance_Basic-before",
			"InstanceTags": converter.CommaStringArray([]string{"terraform"}),
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan": "bunny-1",
		}

		paramsUpdated = map[string]string{
			"InstanceName": "TestAccInstance_Basic-after",
			"InstanceTags": converter.CommaStringArray([]string{"terraform", "acceptance-test"}),
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
					resource.TestCheckResourceAttr(instanceResourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(instanceResourceName, "tags.0", "terraform"),
					resource.TestCheckResourceAttr(instanceResourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.configured", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.running", "true"),
				),
			},
			{
				ResourceName:            instanceResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keep_associated_vpc"},
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsUpdated["InstanceName"]),
					resource.TestCheckResourceAttr(instanceResourceName, "plan", "bunny-1"),
					resource.TestCheckResourceAttr(instanceResourceName, "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr(instanceResourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(instanceResourceName, "tags.0", "terraform"),
					resource.TestCheckResourceAttr(instanceResourceName, "tags.1", "acceptance-test"),
					resource.TestCheckResourceAttr(instanceResourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.configured", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.running", "true"),
				),
			},
		},
	})
}

// TestAccInstance_Upgrade: Creating dedicated AWS instance, upgrade plan and verify.
func TestAccInstance_Upgrade(t *testing.T) {
	var (
		fileNames            = []string{"instance", "data_source/nodes"}
		instanceResourceName = "cloudamqp_instance.instance"
		dataSourceNodesName  = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName": "TestAccInstance_Upgrade",
			"InstancePlan": "bunny-1",
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
		}

		paramsUpdated = map[string]string{
			"InstanceName": "TestAccInstance_Upgrade",
			"InstancePlan": "bunny-3",
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
		}
	)
	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(instanceResourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(instanceResourceName, "plan", params["InstancePlan"]),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.configured", "true"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "nodes", "3"),
					resource.TestCheckResourceAttr(instanceResourceName, "plan", paramsUpdated["InstancePlan"]),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "3"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.configured", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.1.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.1.configured", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.2.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.2.configured", "true"),
				),
			},
		},
	})
}

// TestAccInstance_PlanChange: Creating dedicated AWS instance, change plan and verify.
func TestAccInstance_PlanChange(t *testing.T) {
	var (
		fileNames            = []string{"instance"}
		instanceResourceName = "cloudamqp_instance.instance"

		params = map[string]string{
			"InstanceName": "TestAccInstance_PlanChange",
			"InstancePlan": "squirrel-1",
		}

		paramsUpdated = map[string]string{
			"InstanceName": "TestAccInstance_PlanChange",
			"InstancePlan": "bunny-1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(instanceResourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(instanceResourceName, "plan", params["InstancePlan"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(instanceResourceName, "plan", paramsUpdated["InstancePlan"]),
				),
			},
		},
	})
}

// TestAccInstance_Downgrade: Creating dedicated AWS instance, downgrade plan and verify.
func TestAccInstance_Downgrade(t *testing.T) {
	var (
		fileNames            = []string{"instance", "data_source/nodes"}
		instanceResourceName = "cloudamqp_instance.instance"
		dataSourceNodesName  = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName": "TestAccInstance_Downgrade",
			"InstancePlan": "bunny-3",
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
		}

		paramsUpdated = map[string]string{
			"InstanceName": "TestAccInstance_Downgrade",
			"InstancePlan": "bunny-1",
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(instanceResourceName, "nodes", "3"),
					resource.TestCheckResourceAttr(instanceResourceName, "plan", params["InstancePlan"]),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "3"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.configured", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.1.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.1.configured", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.2.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.2.configured", "true"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(instanceResourceName, "plan", paramsUpdated["InstancePlan"]),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.configured", "true"),
				),
			},
		},
	})
}

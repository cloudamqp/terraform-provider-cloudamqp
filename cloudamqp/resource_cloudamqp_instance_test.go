package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/converter"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccInstance_Basic: Creating dedicated AWS instance, do some minor updates, import and read
// nodes data source.
func TestAccInstance_Basic(t *testing.T) {
	var (
		fileNames    = []string{"instance"}
		instanceName = "cloudamqp_instance.instance"

		params = map[string]string{
			"InstanceName": "TestAccInstance_Basic-before",
			"InstanceTags": converter.CommaStringArray([]string{"terraform"}),
		}

		paramsUpdated = map[string]string{
			"InstanceName": "TestAccInstance_Basic-after",
			"InstanceTags": converter.CommaStringArray([]string{"terraform", "acceptance-test"}),
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
					resource.TestCheckResourceAttr(instanceName, "nodes", "1"),
					resource.TestCheckResourceAttr(instanceName, "plan", "bunny-1"),
					resource.TestCheckResourceAttr(instanceName, "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr(instanceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(instanceName, "tags.0", "terraform"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceName, "name", paramsUpdated["InstanceName"]),
					resource.TestCheckResourceAttr(instanceName, "nodes", "1"),
					resource.TestCheckResourceAttr(instanceName, "plan", "bunny-1"),
					resource.TestCheckResourceAttr(instanceName, "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr(instanceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(instanceName, "tags.0", "terraform"),
					resource.TestCheckResourceAttr(instanceName, "tags.1", "acceptance-test"),
				),
			},
			{
				ResourceName:            instanceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keep_associated_vpc"},
			},
		},
	})
}

// TestAccInstance_PlanChange: Creating dedicated AWS instance, change plan and verify.
func TestAccInstance_PlanChange(t *testing.T) {
	var (
		fileNames    = []string{"instance"}
		instanceName = "cloudamqp_instance.instance"

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
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(instanceName, "nodes", "1"),
					resource.TestCheckResourceAttr(instanceName, "plan", params["InstancePlan"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceName, "nodes", "1"),
					resource.TestCheckResourceAttr(instanceName, "plan", paramsUpdated["InstancePlan"]),
				),
			},
		},
	})
}

// TestAccInstance_Upgrade: Creating dedicated AWS instance, upgrade plan and verify.
func TestAccInstance_Upgrade(t *testing.T) {
	var (
		fileNames           = []string{"instance", "nodes_data"}
		instanceName        = "cloudamqp_instance.instance"
		dataSourceNodesName = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName": "TestAccInstance_Upgrade",
			"InstancePlan": "bunny-1",
			"InstanceID":   fmt.Sprintf("%s.id", instanceName),
		}

		paramsUpdated = map[string]string{
			"InstanceName": "TestAccInstance_Upgrade",
			"InstancePlan": "bunny-3",
			"InstanceID":   fmt.Sprintf("%s.id", instanceName),
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(instanceName, "nodes", "1"),
					resource.TestCheckResourceAttr(instanceName, "plan", params["InstancePlan"]),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.configured", "true"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceName, "nodes", "3"),
					resource.TestCheckResourceAttr(instanceName, "plan", paramsUpdated["InstancePlan"]),
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

// TestAccInstance_Downgrade: Creating dedicated AWS instance, downgrade plan and verify.
func TestAccInstance_Downgrade(t *testing.T) {
	var (
		fileNames           = []string{"instance", "nodes_data"}
		instanceName        = "cloudamqp_instance.instance"
		dataSourceNodesName = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName": "TestAccInstance_Downgrade",
			"InstancePlan": "bunny-3",
			"InstanceID":   fmt.Sprintf("%s.id", instanceName),
		}

		paramsUpdated = map[string]string{
			"InstanceName": "TestAccInstance_Downgrade",
			"InstancePlan": "bunny-1",
			"InstanceID":   fmt.Sprintf("%s.id", instanceName),
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(instanceName, "nodes", "3"),
					resource.TestCheckResourceAttr(instanceName, "plan", params["InstancePlan"]),
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
					resource.TestCheckResourceAttr(instanceName, "nodes", "1"),
					resource.TestCheckResourceAttr(instanceName, "plan", paramsUpdated["InstancePlan"]),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.running", "true"),
					resource.TestCheckResourceAttr(dataSourceNodesName, "nodes.0.configured", "true"),
				),
			},
		},
	})
}

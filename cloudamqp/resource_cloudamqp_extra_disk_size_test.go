package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccExtraDiskSize_AWS_Basic: Add extra disk size to an instance hosted in AWS.
func TestAccExtraDiskSize_AWS_Basic(t *testing.T) {
	var (
		fileNames     = []string{"instance", "extra_disk"}
		instanceName  = "cloudamqp_instance.instance"
		dataNodesName = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName":   "TestAccExtraDiskSize_AWS_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceName),
			"InstanceRegion": "amazon-web-services::us-east-1",
			"ExtraDiskSize":  "25",
			"AllowDowntime":  "false",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(dataNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataNodesName, "nodes.0.disk_size", "10"),
					resource.TestCheckResourceAttr(dataNodesName, "nodes.0.additional_disk_size", "25"),
				),
			},
		},
	})
}

// TestAccExtraDiskSize_GCE_Basic: Add extra disk size to an instance hosted in Google.
func TestAccExtraDiskSize_GCE_Basic(t *testing.T) {
	var (
		fileNames     = []string{"instance", "extra_disk"}
		instanceName  = "cloudamqp_instance.instance"
		dataNodesName = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName":   "TestAccExtraDiskSize_GCE_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceName),
			"InstanceRegion": "google-compute-engine::us-east1",
			"ExtraDiskSize":  "25",
			"AllowDowntime":  "false",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(dataNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataNodesName, "nodes.0.disk_size", "20"),
					resource.TestCheckResourceAttr(dataNodesName, "nodes.0.additional_disk_size", "25"),
				),
			},
		},
	})
}

// TestAccExtraDiskSize_Azure_Basic: Add extra disk size to an instance hosted in Azure.
func TestAccExtraDiskSize_Azure_Basic(t *testing.T) {
	var (
		fileNames     = []string{"instance", "extra_disk"}
		instanceName  = "cloudamqp_instance.instance"
		dataNodesName = "data.cloudamqp_nodes.nodes"

		params = map[string]string{
			"InstanceName":   "TestAccExtraDiskSize_Azure_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceName),
			"InstanceRegion": "azure-arm::eastus",
			"ExtraDiskSize":  "25",
			"AllowDowntime":  "true",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(dataNodesName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataNodesName, "nodes.0.disk_size", "8"),
					resource.TestCheckResourceAttr(dataNodesName, "nodes.0.additional_disk_size", "25"),
				),
			},
		},
	})
}

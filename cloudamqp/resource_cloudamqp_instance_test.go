package cloudamqp

import (
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/converter"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccInstance_Basic: Creating dedicated AWS instance, do some minor updates and import.
func TestAccInstance_Basic(t *testing.T) {
	var (
		fileNames    = []string{"instance"}
		resourceName = "cloudamqp_instance.instance"
		params       = map[string]string{
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
					resource.TestCheckResourceAttr(resourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", "bunny-1"),
					resource.TestCheckResourceAttr(resourceName, "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", paramsUpdated["InstanceName"]),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", "bunny-1"),
					resource.TestCheckResourceAttr(resourceName, "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "acceptance-test"),
				),
			},
			{
				ResourceName:            resourceName,
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
		resourceName = "cloudamqp_instance.instance"
		params       = map[string]string{
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
					resource.TestCheckResourceAttr(resourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", params["InstancePlan"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", paramsUpdated["InstancePlan"]),
				),
			},
		},
	})
}

// TestAccInstance_Upgrade: Creating dedicated AWS instance, upgrade plan and verify.
func TestAccInstance_Upgrade(t *testing.T) {
	var (
		fileNames    = []string{"instance"}
		resourceName = "cloudamqp_instance.instance"
		params       = map[string]string{
			"InstanceName": "TestAccInstance_Upgrade",
			"InstancePlan": "bunny-1",
		}
		paramsUpdated = map[string]string{
			"InstanceName": "TestAccInstance_Upgrade",
			"InstancePlan": "bunny-3",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", params["InstancePlan"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodes", "3"),
					resource.TestCheckResourceAttr(resourceName, "plan", paramsUpdated["InstancePlan"]),
				),
			},
		},
	})
}

// TestAccInstance_Downgrade: Creating dedicated AWS instance, downgrade plan and verify.
func TestAccInstance_Downgrade(t *testing.T) {
	var (
		fileNames    = []string{"instance"}
		resourceName = "cloudamqp_instance.instance"
		params       = map[string]string{
			"InstanceName": "TestAccInstance_Downgrade",
			"InstancePlan": "bunny-3",
		}
		paramsUpdated = map[string]string{
			"InstanceName": "TestAccInstance_Downgrade",
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
					resource.TestCheckResourceAttr(resourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(resourceName, "nodes", "3"),
					resource.TestCheckResourceAttr(resourceName, "plan", params["InstancePlan"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", paramsUpdated["InstancePlan"]),
				),
			},
		},
	})
}

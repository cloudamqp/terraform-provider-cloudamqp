package cloudamqp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Basic instance test case. Creating dedicated AWS instance and do some minor updates.
func TestAccInstance_Basics(t *testing.T) {
	var (
		resourceName = "cloudamqp_instance.instance"
		params       = map[string]any{
			"InstanceName": "TestAccInstance_Basics-before",
		}
		paramsUpdated = map[string]any{
			"InstanceName": "TestAccInstance_Basics-after",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: loadTemplatedConfig(t, "cloudamqp_instance", params),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", params["InstanceName"].(string)),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", "bunny-1"),
					resource.TestCheckResourceAttr(resourceName, "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
				),
			},
			{
				Config: loadTemplatedConfig(t, "cloudamqp_instance", paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", paramsUpdated["InstanceName"].(string)),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", "bunny-1"),
					resource.TestCheckResourceAttr(resourceName, "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
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

func TestAccInstance_PlanChange(t *testing.T) {
	var (
		resourceName = "cloudamqp_instance.instance"
		params       = map[string]any{
			"InstanceName": "TestAccInstance_PlanChange",
			"InstancePlan": "squirrel-1",
		}
		paramsUpdated = map[string]any{
			"InstanceName": "TestAccInstance_PlanChange",
			"InstancePlan": "bunny-1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: loadTemplatedConfig(t, "cloudamqp_instance", params),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", params["InstanceName"].(string)),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", params["InstancePlan"].(string)),
				),
			},
			{
				Config: loadTemplatedConfig(t, "cloudamqp_instance", paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", paramsUpdated["InstancePlan"].(string)),
				),
			},
		},
	})
}

func TestAccInstance_Upgrade(t *testing.T) {
	var (
		resourceName = "cloudamqp_instance.instance"
		params       = map[string]any{
			"InstanceName": "TestAccInstance_Upgrade",
			"InstancePlan": "bunny-1",
		}
		paramsUpdated = map[string]any{
			"InstanceName": "TestAccInstance_Upgrade",
			"InstancePlan": "bunny-3",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: loadTemplatedConfig(t, "cloudamqp_instance", params),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", params["InstanceName"].(string)),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", params["InstancePlan"].(string)),
				),
			},
			{
				Config: loadTemplatedConfig(t, "cloudamqp_instance", paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodes", "3"),
					resource.TestCheckResourceAttr(resourceName, "plan", paramsUpdated["InstancePlan"].(string)),
				),
			},
		},
	})
}

func TestAccInstance_Downgrade(t *testing.T) {
	var (
		resourceName = "cloudamqp_instance.instance"
		params       = map[string]any{
			"InstanceName": "TestAccInstance_Downgrade",
			"InstancePlan": "bunny-3",
		}
		paramsUpdated = map[string]any{
			"InstanceName": "TestAccInstance_Downgrade",
			"InstancePlan": "bunny-1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: loadTemplatedConfig(t, "cloudamqp_instance", params),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", params["InstanceName"].(string)),
					resource.TestCheckResourceAttr(resourceName, "nodes", "3"),
					resource.TestCheckResourceAttr(resourceName, "plan", params["InstancePlan"].(string)),
				),
			},
			{
				Config: loadTemplatedConfig(t, "cloudamqp_instance", paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", paramsUpdated["InstancePlan"].(string)),
				),
			},
		},
	})
}

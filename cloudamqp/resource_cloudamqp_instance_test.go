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
			"InstanceName": "terraform-before",
			"InstancePlan": "bunny-1",
		}
		params_updated = map[string]any{
			"InstanceName": "terraform-after",
			"InstancePlan": "bunny-1",
		}
	)

	config, err := loadTemplatedConfig("cloudamqp_instance", "main.tf", params)
	if err != nil {
		t.Fatalf("failed to load configuration, err: %v", err)
	}
	config_update, err := loadTemplatedConfig("cloudamqp_instance", "main.tf", params_updated)
	if err != nil {
		t.Fatalf("failed to load configuration, err: %v", err)
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", params["InstanceName"].(string)),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", params["InstancePlan"].(string)),
					resource.TestCheckResourceAttr(resourceName, "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
				),
			},
			{
				Config: config_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", params_updated["InstanceName"].(string)),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", params_updated["InstancePlan"].(string)),
					resource.TestCheckResourceAttr(resourceName, "region", "amazon-web-services::us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "terraform"),
				),
			},
		},
	})
}

func TestAccInstance_PlanChange(t *testing.T) {
	var (
		resourceName = "cloudamqp_instance.instance"
		params       = map[string]any{
			"InstanceName": "Instance plan change",
			"InstancePlan": "squirrel-1",
		}
		params_updated = map[string]any{
			"InstanceName": "Instance plan change",
			"InstancePlan": "bunny-1",
		}
	)

	config, err := loadTemplatedConfig("cloudamqp_instance", "main.tf", params)
	if err != nil {
		t.Fatalf("failed to load configuration, err: %v", err)
	}
	config_update, err := loadTemplatedConfig("cloudamqp_instance", "main.tf", params_updated)
	if err != nil {
		t.Fatalf("failed to load configuration, err: %v", err)
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", params["InstanceName"].(string)),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", params["InstancePlan"].(string)),
				),
			},
			{
				Config: config_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", params_updated["InstancePlan"].(string)),
				),
			},
		},
	})
}

func TestAccInstance_Upgrade(t *testing.T) {
	var (
		resourceName = "cloudamqp_instance.instance"
		params       = map[string]any{
			"InstanceName": "Instance plan changes",
			"InstancePlan": "bunny-1",
		}
		params_updated = map[string]any{
			"InstanceName": "Instance plan changes",
			"InstancePlan": "bunny-3",
		}
	)

	config, err := loadTemplatedConfig("cloudamqp_instance", "main.tf", params)
	if err != nil {
		t.Fatalf("failed to load configuration, err: %v", err)
	}
	config_update, err := loadTemplatedConfig("cloudamqp_instance", "main.tf", params_updated)
	if err != nil {
		t.Fatalf("failed to load configuration, err: %v", err)
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", params["InstanceName"].(string)),
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", params["InstancePlan"].(string)),
				),
			},
			{
				Config: config_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodes", "3"),
					resource.TestCheckResourceAttr(resourceName, "plan", params_updated["InstancePlan"].(string)),
				),
			},
		},
	})
}

func TestAccInstance_Downgrade(t *testing.T) {
	var (
		resourceName = "cloudamqp_instance.instance"
		params       = map[string]any{
			"InstanceName": "Instance plan changes",
			"InstancePlan": "bunny-3",
		}
		params_updated = map[string]any{
			"InstanceName": "Instance plan changes",
			"InstancePlan": "bunny-1",
		}
	)

	config, err := loadTemplatedConfig("cloudamqp_instance", "main.tf", params)
	if err != nil {
		t.Fatalf("failed to load configuration, err: %v", err)
	}
	config_update, err := loadTemplatedConfig("cloudamqp_instance", "main.tf", params_updated)
	if err != nil {
		t.Fatalf("failed to load configuration, err: %v", err)
	}

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", params["InstanceName"].(string)),
					resource.TestCheckResourceAttr(resourceName, "nodes", "3"),
					resource.TestCheckResourceAttr(resourceName, "plan", params["InstancePlan"].(string)),
				),
			},
			{
				Config: config_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "nodes", "1"),
					resource.TestCheckResourceAttr(resourceName, "plan", params_updated["InstancePlan"].(string)),
				),
			},
		},
	})
}

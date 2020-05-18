package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Scale up test case for dedicated AWS instance. Test if instance could be scaled up to 3 nodes.
func TestAccInstance_Scale_Up(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestAccInstance_Scale_Up, since test is in short mode")
	}
	resource_name := "cloudamqp_instance.instance"
	name := "terraform-scale-up"
	region := "amazon-web-services::us-east-1"
	plan := "bunny"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy(resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig_Custom_Scale(name, region, plan, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", name),
					resource.TestCheckResourceAttr(resource_name, "nodes", "1"),
					resource.TestCheckResourceAttr(resource_name, "plan", plan),
					resource.TestCheckResourceAttr(resource_name, "region", region),
					resource.TestCheckResourceAttr(resource_name, "rmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resource_name, "tags.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "tags.0", "terraform"),
				),
			},
			{
				Config: testAccInstanceConfig_Custom_Scale(name, region, "rabbit", 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", name),
					resource.TestCheckResourceAttr(resource_name, "nodes", "3"),
					resource.TestCheckResourceAttr(resource_name, "plan", "rabbit"),
					resource.TestCheckResourceAttr(resource_name, "region", region),
					resource.TestCheckResourceAttr(resource_name, "rmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resource_name, "tags.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "tags.0", "terraform"),
				),
			},
		},
	})
}

// Scale down test case for dedicated AWS instance. Test if instance could be scaled down to 1 nodes.
func TestAccInstance_Scale_Down(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestAccInstance_Scale_Down, since test is in short mode")
	}
	resource_name := "cloudamqp_instance.instance"
	name := "terraform-scale-down"
	region := "amazon-web-services::us-east-1"
	plan := "rabbit"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy(resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig_Custom_Scale(name, region, plan, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", name),
					resource.TestCheckResourceAttr(resource_name, "nodes", "3"),
					resource.TestCheckResourceAttr(resource_name, "plan", plan),
					resource.TestCheckResourceAttr(resource_name, "region", region),
					resource.TestCheckResourceAttr(resource_name, "rmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resource_name, "tags.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "tags.0", "terraform"),
				),
			},
			{
				Config: testAccInstanceConfig_Custom_Scale(name, region, "rabbit", 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", name),
					resource.TestCheckResourceAttr(resource_name, "nodes", "1"),
					resource.TestCheckResourceAttr(resource_name, "plan", "rabbit"),
					resource.TestCheckResourceAttr(resource_name, "region", region),
					resource.TestCheckResourceAttr(resource_name, "rmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resource_name, "tags.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "tags.0", "terraform"),
				),
			},
		},
	})
}

// Shared AWS test case. Simple test to see if instance is created and removed.
func TestAccInstance_Shared_AWS(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestAccInstance_Shared_AWS, since test is in short mode")
	}
	instance_name := "shared_aws"
	resource_name := fmt.Sprintf("cloudamqp_instance.%s", instance_name)
	name := "terraform-shared-aws"
	region := "amazon-web-services::us-east-1"
	plan := "lemur"
	version := "3.7.14" // Could change depedning on shared server beeing used.

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy(resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig_Basic_Without_VPC(instance_name, name, region, plan, version),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", name),
					resource.TestCheckResourceAttr(resource_name, "nodes", "1"),
					resource.TestCheckResourceAttr(resource_name, "plan", plan),
					resource.TestCheckResourceAttr(resource_name, "region", region),
					resource.TestCheckResourceAttr(resource_name, "tags.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "tags.0", "terraform"),
				),
			},
		},
	})
}

// Dedicated Azure test case. Simple test to see if instance is created and removed.
func TestAccInstance_Dedicated_Azure(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestAccInstance_Dedicated_Azure, since test is in short mode")
	}
	instance_name := "dedicated_azure"
	resource_name := fmt.Sprintf("cloudamqp_instance.%s", instance_name)
	name := "terraform-dedicated-azure"
	region := "azure-arm::east-us"
	plan := "bunny"
	version := "3.8.2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy(resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig_Basic_Without_VPC(instance_name, name, region, plan, version),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", name),
					resource.TestCheckResourceAttr(resource_name, "nodes", "1"),
					resource.TestCheckResourceAttr(resource_name, "plan", plan),
					resource.TestCheckResourceAttr(resource_name, "region", region),
					resource.TestCheckResourceAttr(resource_name, "rmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resource_name, "tags.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "tags.0", "terraform"),
				),
			},
		},
	})
}

// Shared Azure test case. Simple test to see if instance is created and removed.
func TestAccInstance_Shared_Azure(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestAccInstance_Shared_Azure, since test is in short mode")
	}
	instance_name := "shared_azure"
	resource_name := fmt.Sprintf("cloudamqp_instance.%s", instance_name)
	name := "terraform-shared-azure"
	region := "azure-arm::eastus"
	plan := "lemur"
	version := "3.6.12" // Could change depedning on shared server beeing used.

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy(resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig_Basic_Without_VPC(instance_name, name, region, plan, version),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", name),
					resource.TestCheckResourceAttr(resource_name, "nodes", "1"),
					resource.TestCheckResourceAttr(resource_name, "plan", plan),
					resource.TestCheckResourceAttr(resource_name, "region", region),
					resource.TestCheckResourceAttr(resource_name, "tags.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "tags.0", "terraform"),
				),
			},
		},
	})
}

// Dedicated GCE test case. Simple tet to see if instance is created and removed.
func TestAccInstance_Dedicated_GCE(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestAccInstance_Scale_Up, since test is in short mode")
	}
	instance_name := "dedicated_gce"
	resource_name := fmt.Sprintf("cloudamqp_instance.%s", instance_name)
	name := "terraform-gce"
	region := "google-compute-engine::us-central1"
	plan := "bunny"
	version := "3.8.2"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy(resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig_Basic_Without_VPC(instance_name, name, region, plan, version),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", name),
					resource.TestCheckResourceAttr(resource_name, "nodes", "1"),
					resource.TestCheckResourceAttr(resource_name, "plan", plan),
					resource.TestCheckResourceAttr(resource_name, "region", region),
					resource.TestCheckResourceAttr(resource_name, "rmq_version", "3.8.2"),
					resource.TestCheckResourceAttr(resource_name, "tags.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "tags.0", "terraform"),
				),
			},
		},
	})
}

// Shared GCE test case. Simple test to see if instance is created and removed.
func TestAccInstance_Shared_GCE(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestAccInstance_Scale_Up, since test is in short mode")
	}
	instance_name := "shared_gce"
	resource_name := fmt.Sprintf("cloudamqp_instance.%s", instance_name)
	name := "terraform-shared-gce"
	region := "google-compute-engine::us-central1"
	plan := "lemur"
	version := "3.7.5" // Could change depedning on shared server beeing used.

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy(resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceConfig_Basic_Without_VPC(instance_name, name, region, plan, version),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resource_name),
					resource.TestCheckResourceAttr(resource_name, "name", name),
					resource.TestCheckResourceAttr(resource_name, "nodes", "1"),
					resource.TestCheckResourceAttr(resource_name, "plan", plan),
					resource.TestCheckResourceAttr(resource_name, "region", region),
					resource.TestCheckResourceAttr(resource_name, "tags.#", "1"),
					resource.TestCheckResourceAttr(resource_name, "tags.0", "terraform"),
				),
			},
		},
	})
}

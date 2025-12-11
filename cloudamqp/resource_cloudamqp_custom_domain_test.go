package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Recording the test benefits from having sleep set to default (10 s).
// While replaying the test, using lower sleep (1 s), will speed up the test.

// TestAccCustomDomain_Basic: Create custom domain, import and update hostname.
func TestAccCustomDomain_Basic(t *testing.T) {
	t.Parallel()

	var (
		fileNames                = []string{"instance", "custom_domain"}
		instanceResourceName     = "cloudamqp_instance.instance"
		customDomainResourceName = "cloudamqp_custom_domain.custom_domain"

		params = map[string]string{
			"InstanceName":         "TestAccCustomDomain_Basic",
			"InstanceID":           fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":         "bunny-1",
			"CustomDomainHostname": "test.example.com",
			"CustomDomainSleep":    "1",
			"CustomDomainTimeout":  "1800",
		}

		paramsUpdated = map[string]string{
			"InstanceName":         "TestAccCustomDomain_Basic",
			"InstanceID":           fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":         "bunny-1",
			"CustomDomainHostname": "updated.example.com",
			"CustomDomainSleep":    "1",
			"CustomDomainTimeout":  "1800",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(customDomainResourceName, "hostname", "test.example.com"),
					resource.TestCheckResourceAttr(customDomainResourceName, "sleep", "1"),
					resource.TestCheckResourceAttr(customDomainResourceName, "timeout", "1800"),
				),
			},
			{
				ResourceName:            customDomainResourceName,
				ImportStateIdFunc:       testAccImportCustomDomainStateIdFunc(instanceResourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sleep", "timeout"},
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(customDomainResourceName, "hostname", "updated.example.com"),
					resource.TestCheckResourceAttr(customDomainResourceName, "sleep", "1"),
					resource.TestCheckResourceAttr(customDomainResourceName, "timeout", "1800"),
				),
			},
		},
	})
}

// TestAccCustomDomain_DefaultValues: Test that default sleep and timeout values work correctly.
func TestAccCustomDomain_DefaultValues(t *testing.T) {
	t.Parallel()

	var (
		fileNames                = []string{"instance", "custom_domain"}
		instanceResourceName     = "cloudamqp_instance.instance"
		customDomainResourceName = "cloudamqp_custom_domain.custom_domain"

		params = map[string]string{
			"InstanceName":         "TestAccCustomDomain_DefaultValues",
			"InstanceID":           fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":         "bunny-1",
			"CustomDomainHostname": "default.example.com",
			// Omit CustomDomainSleep and CustomDomainTimeout to use defaults
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(customDomainResourceName, "hostname", "default.example.com"),
					resource.TestCheckResourceAttr(customDomainResourceName, "sleep", "10"),
					resource.TestCheckResourceAttr(customDomainResourceName, "timeout", "1800"),
				),
			},
		},
	})
}

// TestAccCustomDomain_CustomValues: Test that custom sleep and timeout values can be set.
func TestAccCustomDomain_CustomValues(t *testing.T) {
	t.Parallel()

	var (
		fileNames                = []string{"instance", "custom_domain"}
		instanceResourceName     = "cloudamqp_instance.instance"
		customDomainResourceName = "cloudamqp_custom_domain.custom_domain"

		params = map[string]string{
			"InstanceName":         "TestAccCustomDomain_CustomValues",
			"InstanceID":           fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":         "bunny-1",
			"CustomDomainHostname": "custom.example.com",
			"CustomDomainSleep":    "5",
			"CustomDomainTimeout":  "3600",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(customDomainResourceName, "hostname", "custom.example.com"),
					resource.TestCheckResourceAttr(customDomainResourceName, "sleep", "5"),
					resource.TestCheckResourceAttr(customDomainResourceName, "timeout", "3600"),
				),
			},
		},
	})
}

// testAccImportCustomDomainStateIdFunc returns the instance_id for import.
func testAccImportCustomDomainStateIdFunc(instanceResourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[instanceResourceName]
		if !ok {
			return "", fmt.Errorf("Resource %s not found", instanceResourceName)
		}
		return rs.Primary.Attributes["id"], nil
	}
}

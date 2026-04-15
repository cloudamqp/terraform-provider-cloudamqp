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
			"CustomDomainHostname": "vcr-test.ddns.net",
			"CustomDomainSleep":    "1",
			"CustomDomainTimeout":  "1800",
		}

		paramsUpdated = map[string]string{
			"InstanceName":         "TestAccCustomDomain_Basic",
			"InstanceID":           fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":         "bunny-1",
			"CustomDomainHostname": "vcr-update.ddns.net",
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
					resource.TestCheckResourceAttr(customDomainResourceName, "hostname", params["CustomDomainHostname"]),
					resource.TestCheckResourceAttr(customDomainResourceName, "sleep", params["CustomDomainSleep"]),
					resource.TestCheckResourceAttr(customDomainResourceName, "timeout", params["CustomDomainTimeout"]),
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
					resource.TestCheckResourceAttr(customDomainResourceName, "hostname", paramsUpdated["CustomDomainHostname"]),
					resource.TestCheckResourceAttr(customDomainResourceName, "sleep", paramsUpdated["CustomDomainSleep"]),
					resource.TestCheckResourceAttr(customDomainResourceName, "timeout", paramsUpdated["CustomDomainTimeout"]),
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

package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/converter"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// TestAccFirewall_Basic: Create standalone VPC and instance, enable VPC Connect and import.
func TestAccFirewall_Basic(t *testing.T) {
	var (
		// fileNames = []string{"vpc_and_instance", "firewall"}
		fileNames = []string{"firewall"}
		// vpcName      = "cloudamqp_vpc.vpc"
		// instanceName = "cloudamqp_instance.instance"
		resourceName = "cloudamqp_security_firewall.firewall_settings"

		params = map[string]string{
			// "VpcName":      "TestAccFirewall_Basic",
			// "VpcRegion":    "amazon-web-services::us-east-1",
			"InstanceName": "TestAccFirewall_Basic",
			// "InstanceID":   fmt.Sprintf("%s.id", instanceName),
			"InstanceID":         "1706",
			"FirewallIP02":       "192.168.0.0/24",
			"FirewallServices02": converter.CommaStringArray([]string{"AMQP"}),
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// PreventDiskCleanup: true,
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFirewallState(resourceName),
					// resource.TestCheckResourceAttr(vpcName, "name", params["VpcName"]),
					// resource.TestCheckResourceAttr(instanceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
					// resource.TestCheckResourceAttr(resourceName, "rules.0.description", "Default"),
					// resource.TestCheckResourceAttr(resourceName, "rules.0.ip", "0.0.0.0/0"),
					// resource.TestCheckResourceAttr(resourceName, "rules.0.ports.#", "0"),
					// resource.TestCheckResourceAttr(resourceName, "rules.0.services.#", "2"),
					// resource.TestCheckResourceAttr(resourceName, "rules.0.services.0", "AMQPS"),
					// resource.TestCheckResourceAttr(resourceName, "rules.0.services.1", "HTTPS"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// ImportStateVerifyIgnore: []string{"region", "approved_subscriptions", "allowed_projects"},
			},
		},
	})
}

func testAccCheckFirewallState(resourceName string) resource.TestCheckFunc {
	fmt.Println("testAccCheckFirewallState: ", resourceName)
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No id is set for instance resource")
		}

		fmt.Println("ID: ", rs.Primary.ID)
		for key, value := range rs.Primary.Attributes {
			fmt.Println("key: ", key, "value: ", value)
		}

		return nil
	}
}

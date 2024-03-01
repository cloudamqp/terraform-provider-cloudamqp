package cloudamqp

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/converter"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// TestAccFirewall_Basic: Create standalone VPC and instance, enable VPC Connect and import.
func TestAccFirewall_Basic(t *testing.T) {
	var (
		fileNames    = []string{"vpc_and_instance", "firewall"}
		vpcName      = "cloudamqp_vpc.vpc"
		instanceName = "cloudamqp_instance.instance"
		resourceName = "cloudamqp_security_firewall.firewall_settings"

		params = map[string]string{
			"VpcName":      "TestAccFirewall_Basic",
			"InstanceName": "TestAccFirewall_Basic",
			"InstanceID":   fmt.Sprintf("%s.id", instanceName),
		}

		paramsUpdated = map[string]string{
			"VpcName":             "TestAccFirewall_Basic",
			"InstanceName":        "TestAccFirewall_Basic",
			"InstanceID":          fmt.Sprintf("%s.id", instanceName),
			"FirewallIP":          "10.56.72.0/24",
			"FirewallDescription": "VPC Subnet",
			"FirewallServices":    converter.CommaStringArray([]string{"AMQPS"}),
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vpcName, "name", params["VpcName"]),
					resource.TestCheckResourceAttr(instanceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
					testAccCheckFirewallResourcceAttr(resourceName, map[string]string{
						"rules.%s.ip":          "0.0.0.0/0",
						"rules.%s.description": "Default",
						"rules.%s.ports.#":     "0",
						"rules.%s.services.#":  "2",
						"rules.%s.services.0":  "AMQPS",
						"rules.%s.services.1":  "HTTPS",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vpcName, "name", params["VpcName"]),
					resource.TestCheckResourceAttr(instanceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
					testAccCheckFirewallResourcceAttr(resourceName, map[string]string{
						"rules.%s.ip":          "10.56.72.0/24",
						"rules.%s.description": "VPC Subnet",
						"rules.%s.ports.#":     "0",
						"rules.%s.services.#":  "1",
						"rules.%s.services.0":  "AMQPS",
					}),
				),
			},
		},
	})
}

func testAccCheckFirewallResourcceAttr(resourceName string, params map[string]string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resourceName)
		}

		fmt.Println(params)
		uniqueId, err := findFirewallRuleUnqieIdByIp(rs, params["rules.%s.ip"])
		if err != nil {
			return err
		}

		for key, value := range params {
			realKey := fmt.Sprintf(key, uniqueId)
			if rs.Primary.Attributes[realKey] != value {
				return fmt.Errorf("failed to validate key: %s, with value: %s", realKey, value)
			}
		}
		return nil
	}
}

// rules are made up by TypeSet schema type. When stored each set block gets a unique id instead
// of its position. So instead of 0,1,2,n fetch the unique id
// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/helper/schema#SchemaSetFunc
func findFirewallRuleUnqieIdByIp(resource *terraform.ResourceState, ip string) (string, error) {
	for key, value := range resource.Primary.Attributes {
		if regexp.MustCompile(`^rules.\d+.ip$`).MatchString(key) && value == ip {
			keySplit := strings.Split(key, ".")
			return keySplit[1], nil
		}
	}
	return "", fmt.Errorf("couldn't find an unqie id")
}

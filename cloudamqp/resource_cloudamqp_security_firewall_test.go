package cloudamqp

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/converter"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccFirewall_Basic: Create standalone VPC and instance with firewall rule, import and update rules.
func TestAccFirewall_Basic(t *testing.T) {
	t.Parallel()

	var (
		fileNames            = []string{"vpc_and_instance", "firewall"}
		vpcResourceName      = "cloudamqp_vpc.vpc"
		instanceResourceName = "cloudamqp_instance.instance"
		firewallResourceName = "cloudamqp_security_firewall.firewall_settings"

		params = map[string]string{
			"VpcName":      "TestAccFirewall_Basic",
			"InstanceName": "TestAccFirewall_Basic",
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan": "bunny-1",
		}

		paramsUpdated = map[string]string{
			"VpcName":             "TestAccFirewall_Basic",
			"InstanceName":        "TestAccFirewall_Basic",
			"InstanceID":          fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":        "bunny-1",
			"FirewallIP":          "10.56.72.0/24",
			"FirewallDescription": "VPC Subnet",
			"FirewallServices":    converter.CommaStringArray([]string{"AMQPS"}),
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vpcResourceName, "name", params["VpcName"]),
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(firewallResourceName, "rules.#", "1"),
					testAccCheckFirewallResourcceAttr(firewallResourceName, map[string]string{
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
				ResourceName:      firewallResourceName,
				ImportStateIdFunc: testAccImportStateIdFunc(firewallResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(vpcResourceName, "name", params["VpcName"]),
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(firewallResourceName, "rules.#", "1"),
					testAccCheckFirewallResourcceAttr(firewallResourceName, map[string]string{
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

// TestAccFirewall_Import: Create VPC and instance with two firewall rules and import firewall settings.
func TestAccFirewall_Import(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_vpc" "vpc" {
						name   = "TestAccFirewall_Import"
						region = "amazon-web-services::eu-central-1"
						subnet = "10.56.72.0/24"
						tags   = ["vcr-test"]
					}

					resource "cloudamqp_instance" "instance" {
						name                = "TestAccFirewall_Import"
						plan                = "penguin-1"
						region              = "amazon-web-services::eu-central-1"
						tags                = ["vcr-test"]
						vpc_id              = cloudamqp_vpc.vpc.id
						keep_associated_vpc = true
					}

					resource "cloudamqp_security_firewall" "this" {
						instance_id = cloudamqp_instance.instance.id

						rules {
							description = "MGMT Interface"
							ip          = "0.0.0.0/0"
							ports       = []
							services    = ["HTTPS"]
						}

						rules {
							description = "VPC Subnet"
							ip          = cloudamqp_vpc.vpc.subnet
							ports       = []
							services    = ["AMQP", "AMQPS"]
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_security_firewall.this", "rules.#", "2"),
					testAccCheckFirewallResourcceAttr("cloudamqp_security_firewall.this", map[string]string{
						"rules.%s.ip":          "0.0.0.0/0",
						"rules.%s.description": "MGMT Interface",
						"rules.%s.ports.#":     "0",
						"rules.%s.services.#":  "1",
						"rules.%s.services.0":  "HTTPS",
					}),
					testAccCheckFirewallResourcceAttr("cloudamqp_security_firewall.this", map[string]string{
						"rules.%s.ip":          "10.56.72.0/24",
						"rules.%s.description": "VPC Subnet",
						"rules.%s.ports.#":     "0",
						"rules.%s.services.#":  "2",
						"rules.%s.services.0":  "AMQP",
						"rules.%s.services.1":  "AMQPS",
					}),
				),
			},
			{
				ResourceName:      "cloudamqp_security_firewall.this",
				ImportStateIdFunc: testAccImportStateIdFunc("cloudamqp_security_firewall.this"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccFirewall_Update: Create VPC and instance with a firewall rule and add one more rule.
func TestAccFirewall_Update(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_vpc" "vpc" {
						name   = "TestAccFirewall_Update"
						region = "amazon-web-services::eu-central-1"
						subnet = "10.56.72.0/24"
						tags   = ["vcr-test"]
					}

					resource "cloudamqp_instance" "instance" {
						name                = "TestAccFirewall_Update"
						plan                = "penguin-1"
						region              = "amazon-web-services::eu-central-1"
						tags                = ["vcr-test"]
						vpc_id              = cloudamqp_vpc.vpc.id
						keep_associated_vpc = true
					}

					resource "cloudamqp_security_firewall" "this" {
						instance_id = cloudamqp_instance.instance.id

						rules {
							description = "VPC Subnet"
							ip          = cloudamqp_vpc.vpc.subnet
							ports       = []
							services    = ["AMQP", "AMQPS"]
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_security_firewall.this", "rules.#", "1"),
					testAccCheckFirewallResourcceAttr("cloudamqp_security_firewall.this", map[string]string{
						"rules.%s.ip":          "10.56.72.0/24",
						"rules.%s.description": "VPC Subnet",
						"rules.%s.ports.#":     "0",
						"rules.%s.services.#":  "2",
						"rules.%s.services.0":  "AMQP",
						"rules.%s.services.1":  "AMQPS",
					}),
				),
			},
			{
				Config: `
					resource "cloudamqp_vpc" "vpc" {
						name   = "TestAccFirewall_Update"
						region = "amazon-web-services::eu-central-1"
						subnet = "10.56.72.0/24"
						tags   = ["vcr-test"]
					}

					resource "cloudamqp_instance" "instance" {
						name                = "TestAccFirewall_Update"
						plan                = "penguin-1"
						region              = "amazon-web-services::eu-central-1"
						tags                = ["vcr-test"]
						vpc_id              = cloudamqp_vpc.vpc.id
						keep_associated_vpc = true
					}

					resource "cloudamqp_security_firewall" "this" {
						instance_id = cloudamqp_instance.instance.id

						rules {
							description = "MGMT Interface"
							ip          = "0.0.0.0/0"
							ports       = []
							services    = ["HTTPS"]
						}

						rules {
							description = "VPC Subnet"
							ip          = cloudamqp_vpc.vpc.subnet
							ports       = []
							services    = ["AMQP", "AMQPS"]
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_security_firewall.this", "rules.#", "2"),
					testAccCheckFirewallResourcceAttr("cloudamqp_security_firewall.this", map[string]string{
						"rules.%s.ip":          "0.0.0.0/0",
						"rules.%s.description": "MGMT Interface",
						"rules.%s.ports.#":     "0",
						"rules.%s.services.#":  "1",
						"rules.%s.services.0":  "HTTPS",
					}),
					testAccCheckFirewallResourcceAttr("cloudamqp_security_firewall.this", map[string]string{
						"rules.%s.ip":          "10.56.72.0/24",
						"rules.%s.description": "VPC Subnet",
						"rules.%s.ports.#":     "0",
						"rules.%s.services.#":  "2",
						"rules.%s.services.0":  "AMQP",
						"rules.%s.services.1":  "AMQPS",
					}),
				),
			},
		},
	})
}

func testAccCheckFirewallResourcceAttr(firewallResourceName string, params map[string]string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[firewallResourceName]
		if !ok {
			return fmt.Errorf("Resource: %s not found", firewallResourceName)
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

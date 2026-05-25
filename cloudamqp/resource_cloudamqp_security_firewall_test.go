package cloudamqp

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccFirewall_Import: Create VPC and instance with firewall rule and import firewall settings.
func TestAccFirewall_Import(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
          resource "cloudamqp_vpc" "vpc" {
            name    = "TestAccFirewall_Import"
            region  = "amazon-web-services::us-east-1"
            subnet  = "10.56.72.0/24"
            tags    = ["vcr-test"]
          }

          resource "cloudamqp_instance" "instance" {
            name                = "TestAccFirewall_Import"
            plan                = "penguin-1"
            region              = "amazon-web-services::us-east-1"
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
              ip        = cloudamqp_vpc.vpc.subnet
              ports     = []
              services  = ["AMQP","AMQPS"]
            }
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc", "name", "TestAccFirewall_Import"),
					resource.TestCheckResourceAttr("cloudamqp_instance.instance", "name", "TestAccFirewall_Import"),
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

// TestAccFirewall_Update: Create VPC and instance with firewall rule and update.
func TestAccFirewall_Update(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
          resource "cloudamqp_vpc" "vpc" {
            name    = "TestAccFirewall_Update"
            region  = "amazon-web-services::us-east-1"
            subnet  = "10.56.72.0/24"
            tags    = ["vcr-test"]
          }

          resource "cloudamqp_instance" "instance" {
            name                = "TestAccFirewall_Update"
            plan                = "penguin-1"
            region              = "amazon-web-services::us-east-1"
            tags                = ["vcr-test"]
            vpc_id              = cloudamqp_vpc.vpc.id
            keep_associated_vpc = true
          }

          resource "cloudamqp_security_firewall" "this" {
            instance_id = cloudamqp_instance.instance.id

            rules {
              description = "VPC Subnet"
              ip        = cloudamqp_vpc.vpc.subnet
              ports     = []
              services  = ["AMQP","AMQPS"]
            }
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc", "name", "TestAccFirewall_Update"),
					resource.TestCheckResourceAttr("cloudamqp_instance.instance", "name", "TestAccFirewall_Update"),
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
            name    = "TestAccFirewall_Update"
            region  = "amazon-web-services::us-east-1"
            subnet  = "10.56.72.0/24"
            tags    = ["vcr-test"]
          }

          resource "cloudamqp_instance" "instance" {
            name                = "TestAccFirewall_Update"
            plan                = "penguin-1"
            region              = "amazon-web-services::us-east-1"
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
              ip        = cloudamqp_vpc.vpc.subnet
              ports     = []
              services  = ["AMQP","AMQPS"]
            }
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_vpc.vpc", "name", "TestAccFirewall_Update"),
					resource.TestCheckResourceAttr("cloudamqp_instance.instance", "name", "TestAccFirewall_Update"),
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

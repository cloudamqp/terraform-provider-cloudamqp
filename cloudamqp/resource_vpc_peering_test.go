package cloudamqp

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVpcPeering_Basic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestAccVpcPeering_Basic, since test is in short mode")
	}

	accepter_instance_name := "cloudamqp_instance.accepter_instance"
	requester_instance_name := "cloudamqp_instance.requester_instance"
	accepter_name := "cloudamqp_vpc_peering.vpc_accept_peering"
	requester_name := "cloudamqp_vpc_peering.vpc"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVpc(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccVpcPeeringDestroy(instance_name, resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcPeeringConfig_Instances(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringExists(instance_name),
					resource.TestCheckResourceAttr(resource_name, "rules.#", "1"),
				),
			},
			{
				Config: testAccVpcPeeringConfig_AWS(),
			},
			{
				Config: testAccVpcPeeringConfig_Peering(),
			},
			{
				Config: testAccVpcPeering_Accept(),
			},
		},
	})
}

func testAccCheckVpcPeeringExists(instance_name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set for instance")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadFirewallSettings(instance_id)
		log.Printf("[DEBUG] resource_plugin::testAccCheckPluginEnabled data: %v", data)
		if err != nil {
			return fmt.Errorf("Error fetching item %s", err)
		}
		if data != nil {
			return fmt.Errorf("Error security firewall doesn't exists")
		}
		return nil
	}
}

func testAccSecurityFirewallDestroy(instance_name, resource_name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set for instance")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		data, err := api.ReadFirewallSettings(instance_id)
		if data != nil || err == nil {
			return fmt.Errorf("Firewall still exists")
		}
		return nil
	}
}

func testAccSecurityFirewallConfig_Instances() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "accepter_instance" {
			name 				= "terraform-vpc-accepter"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.3"
			tags 				= ["terraform"]
			vpc_subnet 	= "192.168.0.1/24"
		}

		data "cloudamqp_vpc_info" "vpc_info" {
			instance_id = cloudamqp_instance.accepter_instance.id
		}

		resource "cloudamqp_instance" "requester_instance" {
			name 				= "terraform-vpc-requester"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.3"
			tags 				= ["terraform"]
			vpc_subnet 	= "10.40.72.0/24"
		}
	`)
}

func testAccSecurityFirewallConfig_AWS() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "accepter_instance" {
			name 				= "terraform-vpc-accepter"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.3"
			tags 				= ["terraform"]
			vpc_subnet 	= "192.168.0.1/24"
		}

		data "cloudamqp_vpc_info" "vpc_info" {
			instance_id = cloudamqp_instance.accepter_instance.id
		}

		resource "cloudamqp_instance" "requester_instance" {
			name 				= "terraform-vpc-requester"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.3"
			tags 				= ["terraform"]
			vpc_subnet 	= "10.40.72.0/24"
		}

		provider "aws" {
			region = "eu-north-1"
			access_key = %s
			secret_key = %s
		}

		data "aws_instance" "aws_instance" {
			provider = aws

			instance_tags = {
				Name   = cloudamqp_instance.requester_instance.host
			}
		}

		data "aws_subnet" "subnet" {
			provider = aws
			id = data.aws_instance.aws_instance.subnet_id
		}
	`, AWS_KEY, AWS_SECRET)
}

func testAccSecurityFirewallConfig_Peering() string {
	AWS_KEY := os.Getenv("AWS_KEY")
	AWS_SECRET := os.Getenv("AWS_SECRET")
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "accepter_instance" {
			name 				= "terraform-vpc-accepter"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.3"
			tags 				= ["terraform"]
			vpc_subnet 	= "192.168.0.1/24"
		}

		data "cloudamqp_vpc_info" "vpc_info" {
			instance_id = cloudamqp_instance.accepter_instance.id
		}

		resource "cloudamqp_instance" "requester_instance" {
			name 				= "terraform-vpc-requester"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.3"
			tags 				= ["terraform"]
			vpc_subnet 	= "10.40.72.0/24"
		}

		provider "aws" {
			region = "eu-north-1"
			access_key = %s
		 	secret_key = %s
		}

		data "aws_instance" "aws_instance" {
			provider = aws

			instance_tags = {
				Name   = cloudamqp_instance.requester_instance.host
			}
		}

		data "aws_subnet" "subnet" {
			provider = aws
			id = data.aws_instance.aws_instance.subnet_id
		}

		resource "aws_vpc_peering_connection" "aws_vpc_peering" {
			provider = aws
			vpc_id = data.aws_subnet.subnet.vpc_id
			peer_vpc_id = data.cloudamqp_vpc_info.vpc_info.id
			peer_owner_id = data.cloudamqp_vpc_info.vpc_info.owner_id
			tags = { Name = "Terraform acceptence test peering" }
		}
		`, AWS_KEY, AWS_SECRET)
}

unc testAccSecurityFirewallConfig_Accept() string {
	AWS_KEY := os.Getenv("AWS_KEY")
	AWS_SECRET := os.Getenv("AWS_SECRET")
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "accepter_instance" {
			name 				= "terraform-vpc-accepter"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.3"
			tags 				= ["terraform"]
			vpc_subnet 	= "192.168.0.1/24"
		}

		data "cloudamqp_vpc_info" "vpc_info" {
			instance_id = cloudamqp_instance.accepter_instance.id
		}

		resource "cloudamqp_instance" "requester_instance" {
			name 				= "terraform-vpc-requester"
			nodes 			= 1
			plan  			= "bunny"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.3"
			tags 				= ["terraform"]
			vpc_subnet 	= "10.40.72.0/24"
		}

		provider "aws" {
			region = "eu-north-1"
			access_key = %s
		 	secret_key = %s
		}

		data "aws_instance" "aws_instance" {
			provider = aws

			instance_tags = {
				Name   = "cloudamqp_instance.requester_instance.host" + "-01"
			}
		}

		data "aws_subnet" "subnet" {
			provider = aws
			id = data.aws_instance.aws_instance.subnet_id
		}

		resource "aws_vpc_peering_connection" "aws_vpc_peering" {
			provider = aws
			vpc_id = data.aws_subnet.subnet.vpc_id
			peer_vpc_id = data.cloudamqp_vpc_info.vpc_info.id
			peer_owner_id = data.cloudamqp_vpc_info.vpc_info.owner_id
			tags = { Name = "Terraform acceptence test peering" }
		}

		resource "cloudamqp_vpc_peering" "vpc_accept_peering" {
			instance_id = cloudamqp_instance.accepter_instance.id
			peering_id = aws_vpc_peering_connection.aws_vpc_peering.id
		}

		data "aws_route_table" "route_table" {
			provider = aws
			vpc_id = data.aws_subnet.subnet.vpc_id
		}

		resource "aws_route" "accepter_route" {
			provider = aws
			route_table_id = data.aws_route_table.route_table.route_table_id
			destination_cidr_block = cloudamqp_instance.accepter_instance.vpc_subnet
			vpc_peering_connection_id = aws_vpc_peering_connection.aws_vpc_peering.id
		}
		`, AWS_KEY, AWS_SECRET)
}

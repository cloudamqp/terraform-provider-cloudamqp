// =========
// CLOUDAMQP
// =========
provider "cloudamqp" {
  apikey  = var.cloudamqp_customer_api_key
}

// === Instance resource ===
resource "cloudamqp_instance" "cloudamqp_instance" {
  name   = "terraform-vpc-accepter-test"
  plan   = "bunny"
  region = "amazon-web-services::us-east-1"
  nodes = 1
  tags   = ["test"]
  rmq_version = "3.7.21"
  vpc_subnet = "10.40.72.0/24"
}

// === VPC data source ===
data "cloudamqp_vpc_info" "vpc_info" {
  instance_id = cloudamqp_instance.cloudamqp_instance.id
}

// Requires terraform init, to initialize aws plugin.
// ===
// AWS
// ===
provider "aws" {
  region = var.aws_region
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
}

// === AWS instance resource ===
data "aws_instance" "aws_instance" {
  provider = aws

  instance_tags = {
    Name   = var.aws_instance_name
  }
}

// === AWS - Subnet data source ===
data "aws_subnet" "subnet" {
  provider = aws
  id = data.aws_instance.aws_instance.subnet_id
}

// === AWS - VPC peering connection ===
resource "aws_vpc_peering_connection" "aws_vpc_peering" {
  provider = aws
  vpc_id = data.aws_subnet.subnet.vpc_id
  peer_vpc_id = data.cloudamqp_vpc_info.vpc_info.id
  peer_owner_id = data.cloudamqp_vpc_info.vpc_info.owner_id
  // Config once VPC peering connection is created
  //requester {
  //  allow_remote_vpc_dns_resolution = true
  //}
  tags = { Name = var.aws_peering_name }
}

// === CLOUDAMQP - VPC accept peering ===
resource "cloudamqp_vpc_peering" "vpc_accept_peering" {
  instance_id = cloudamqp_instance.cloudamqp_instance.id
  peering_id = aws_vpc_peering_connection.aws_vpc_peering.id
}

// === AWS - Route table ===
data "aws_route_table" "route_table" {
  provider = aws
  vpc_id = data.aws_subnet.subnet.vpc_id
}

// === AWS - Route ===
resource "aws_route" "accepter_route" {
  provider = aws
  route_table_id = data.aws_route_table.route_table.route_table_id
  destination_cidr_block = cloudamqp_instance.cloudamqp_instance.vpc_subnet
  vpc_peering_connection_id = aws_vpc_peering_connection.aws_vpc_peering.id
}

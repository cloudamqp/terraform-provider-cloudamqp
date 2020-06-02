---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_vpc_peering"
description: |-
  Accepting VPC peering request from an AWS accepter.
---

# cloudamqp_vpc_peering

This resouce allows you to accepting VPC peering request from an AWS requester. This is only available for CloudAMQP instance hosted in AWS. This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

Only available for dedicated subscription plans.

## Example Usage

One way to manage the vpc peering is to combine CloudAMQP Terraform provider with AWS Terraform provider and run them at the same time.

```hcl
# Configure CloudAMQP provider
resource "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# CloudAMQP - new instance, need to be created with a vpc
resource "cloudamqp_instance" "instance" {
  name   = "terraform-vpc-accepter"a
  plan   = "bunny"
  region = "amazon-web-services::us-east-1"
  nodes = 1
  tags   = ["terraform"]
  rmq_version = "3.8.4"
  vpc_subnet = "10.40.72.0/24"
}

# CloudAMQP - Extract vpc information
data "cloudamqp_vpc_info" "vpc_info" {
  instance_id = cloudamqp_instance.instance.id
}

# Configure AWS provider
provider "aws" {
  region = var.aws_region
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
}

# AWS - retreive instance to get subnet identifier
data "aws_instance" "aws_instance" {
  provider = aws

  instance_tags = {
    Name   = var.aws_instance_name
  }
}

# AWS - retrieve subnet
data "aws_subnet" "subnet" {
  provider = aws
  id = data.aws_instance.aws_instance.subnet_id
}

# AWS - Create peering request
resource "aws_vpc_peering_connection" "aws_vpc_peering" {
  provider = aws
  vpc_id = data.aws_subnet.subnet.vpc_id
  peer_vpc_id = data.cloudamqp_vpc_info.vpc_info.id
  peer_owner_id = data.cloudamqp_vpc_info.vpc_info.owner_id
  tags = { Name = var.aws_peering_name }
}

# CloudAMQP - accept the peering request
resource "cloudamqp_vpc_peering" "vpc_accept_peering" {
  instance_id = cloudamqp_instance.instance.id
  peering_id = aws_vpc_peering_connection.aws_vpc_peering.id
}

# AWS - retrieve the route table created in AWS
data "aws_route_table" "route_table" {
  provider = aws
  vpc_id = data.aws_subnet.subnet.vpc_id
}

# AWS - Once the peering request is accepted, configure routing table on accepter to allow traffic
resource "aws_route" "accepter_route" {
  provider = aws
  route_table_id = data.aws_route_table.route_table.route_table_id
  destination_cidr_block = cloudamqp_instance.instance.vpc_subnet
  vpc_peering_connection_id = aws_vpc_peering_connection.aws_vpc_peering.id
}
```

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance ID.
* `peering_id` - (Required) Peering identifier created by AW peering request.

## Import

`cloudamqp_vpc_peering` can be imported using the CloudAMQP instance identifier.

`terraform import cloudamqp_vpc_peering.aws_vpc_peering <instance_id>`

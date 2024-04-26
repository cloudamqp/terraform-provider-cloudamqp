---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_vpc_peering"
description: |-
  Accepting VPC peering request from an AWS accepter.
---

# cloudamqp_vpc_peering

This resouce allows you to accepting VPC peering request from an AWS requester. This is only available for CloudAMQP instance hosted in AWS.

~> **Note:** Creating a VPC peering will automatically add firewall rules for the peered subnet.
<details>
 <summary>
    <i>Default VPC peering firewall rule</i>
  </summary>

```hcl
rules {
  Description = "VPC peer request"
  ip          = "<VPC peered subnet>"
  ports       = []
  services    = ["AMQP", "AMQPS", "HTTPS", "STREAM", "STREAM_SSL", "STOMP", "STOMPS", "MQTT", "MQTTS"]
}
```
</details>

Pricing is available at [cloudamqp.com](https://www.cloudamqp.com/plans.html).

Only available for dedicated subscription plans.

## Example Usage

One way to manage the vpc peering is to combine CloudAMQP Terraform provider with AWS Terraform provider and run them at the same time.

<details>
  <summary>
    <b>
      <i>AWS VPC peering before v1.16.0</i>
    </b>
  </summary>

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# CloudAMQP - new instance, need to be created with a vpc
resource "cloudamqp_instance" "instance" {
  name   = "terraform-vpc-accepter"
  plan   = "bunny-1"
  region = "amazon-web-services::us-east-1"
  tags   = ["terraform"]
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

  depends_on = [
    cloudamqp_vpc_peering.vpc_accept_peering
  ]
}
```
</details>

<details>
  <summary>
    <b>
      <i>AWS VPC peering from v1.16.0 (Managed VPC)</i>
    </b>
  </summary>

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# CloudAMQP - Managed VPC resource
resource "cloudamqp_vpc" "vpc" {
  name = "<VPC name>"
  region = "amazon-web-services::us-east-1"
  subnet = "10.56.72.0/24"
  tags = ["terraform"]
}

# CloudAMQP - new instance, need to be created with a vpc
resource "cloudamqp_instance" "instance" {
  name   = "terraform-vpc-accepter"
  plan   = "bunny-1"
  region = "amazon-web-services::us-east-1"
  tags   = ["terraform"]
  vpc_id = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

# CloudAMQP - Extract vpc information
data "cloudamqp_vpc_info" "vpc_info" {
  vpc_id = cloudamqp_vpc.vpc.id
  # vpc_id prefered over instance_id
  # instance_id = cloudamqp_instance.instance.id
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
  vpc_id = cloudamqp_vpc.vpc.id
  # vpc_id prefered over instance_id
  # instance_id = cloudamqp_instance.instance.id
  peering_id = aws_vpc_peering_connection.aws_vpc_peering.id
  sleep = 30
  timeout = 600
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

  depends_on = [
    cloudamqp_vpc_peering.vpc_accept_peering
  ]
}
```

 </details>

## Argument Reference

 *Note: this resource require either `instance_id` or `vpc_id` from v1.16.0*

* `instance_id` - (Optional) The CloudAMQP instance identifier.

 ***Deprecated: Changed from required to optional in v1.16.0, will be removed in next major version (v2.0)***

* `vpc_id` - (Optional) The managed VPC identifier.

 ***Note: Introduced as optional in version v1.16.0, will be required in next major version (v2.0)***

* `peering_id` - (Required) Peering identifier created by AW peering request.
* `sleep` - (Optional) Configurable sleep time (seconds) between retries for accepting or removing peering. Default set to 60 seconds.
* `timeout` - (Optional) - Configurable timeout time (seconds) for accepting or removing peering. Default set to 3600 seconds.

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.
* `status`- VPC peering status

## Depedency

*Before v1.16.0*
This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

*From v1.16.0*
This resource depends on CloudAMQP managed VPC identifier, `cloudamqp_vpc.vpc.id` or instance identifier, `cloudamqp_instance.instance.id`.

## Import

Not possible to import this resource.

## Create VPC Peering with additional firewall rules

To create a VPC peering configuration with additional firewall rules, it's required to chain the [cloudamqp_security_firewall](https://registry.terraform.io/providers/cloudamqp/cloudamqp/latest/docs/resources/security_firewall)
resource to avoid parallel conflicting resource calls. You can do this by making the firewall resource depend on the VPC peering resource (`cloudamqp_vpc_peering.vpc_accept_peering`).

Furthermore, since all firewall rules are overwritten, the otherwise automatically added rules for the VPC peering also needs to be added.

See example below.

## Example Usage with additional firewall rules

<details>
  <summary>
    <b>
      <i>VPC peering before v1.16.0</i>
    </b>
  </summary>

```hcl
# AWS - retrieve subnet
data "aws_subnet" "subnet" {
  provider = aws
  id = data.aws_instance.aws_instance.subnet_id
}

# CloudAMQP - accept the peering request
resource "cloudamqp_vpc_peering" "vpc_accept_peering" {
  instance_id = cloudamqp_instance.instance.id
  peering_id = aws_vpc_peering_connection.aws_vpc_peering.id
}

# Firewall rules
resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  # Default VPC peering rule
  rules {
    ip          =  data.aws_instance.aws_instance.subnet_id
    ports       = [15672]
    services    = ["AMQP","AMQPS", "STREAM", "STREAM_SSL"]
    description = "VPC peering for <NETWORK>"
  }

  rules {
    ip          = "192.168.0.0/24"
    ports       = [4567, 4568]
    services    = ["AMQP","AMQPS", "HTTPS"]
  }

  depends_on = [
    cloudamqp_vpc_peering.vpc_accept_peering
  ]
}
```

</details>

<details>
  <summary>
    <b>
      <i>VPC peering from v1.16.0 (Managed VPC)</i>
    </b>
  </summary>

```hcl
# AWS - retrieve subnet
data "aws_subnet" "subnet" {
  provider = aws
  id = data.aws_instance.aws_instance.subnet_id
}

# CloudAMQP - accept the peering request
resource "cloudamqp_vpc_peering" "vpc_accept_peering" {
  vpc_id = cloudamqp_vpc.vpc.id
  # vpc_id prefered over instance_id
  # instance_id = cloudamqp_instance.instance.id
  peering_id = aws_vpc_peering_connection.aws_vpc_peering.id
  sleep = 30
  timeout = 600
}

# AWS - VPC subnet for peering requester
data "aws_vpc" "requester_vpc" {
  id = data.aws_subnet.subnet.vpc_id
}

# CloudAMQP - Managed firewall rules
resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  # Default VPC peering rule
  rules {
    ip          =  data.aws_vpc.requester_vpc.cidr_block
    ports       = [15672]
    services    = ["AMQP","AMQPS", "STREAM", "STREAM_SSL"]
    description = "VPC peering for <NETWORK>"
  }

  rules {
    ip          = "0.0.0.0/0"
    ports       = []
    services    = ["HTTPS"]
    description = "MGMT interface"
  }

  depends_on = [
    cloudamqp_vpc_peering.vpc_accept_peering
  ]
}
```

</details>

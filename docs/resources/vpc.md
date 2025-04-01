---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_vpc"
description: |-
  Managed VPC resource.
---

# cloudamqp_vpc

This resource allows you to manage standalone VPC.

New Cloudamqp instances can be added to the managed VPC. Set the instance *vpc_id* attribute to the
managed vpc identifier, see example below, when creating the instance.

Only available for dedicated subscription plans.

Pricing is available at [CloudAMQP plans].

## Example Usage

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# Managed VPC resource
resource "cloudamqp_vpc" "vpc" {
  name    = "<VPC name>"
  region  = "amazon-web-services::us-east-1"
  subnet  = "10.56.72.0/24"
  tags    = []
}

#  New instance, need to be created with a vpc
resource "cloudamqp_instance" "instance" {
  name                = "<Instance name>"
  plan                = "bunny-1"
  region              = "amazon-web-services::us-east-1"
  nodes               = 1
  tags                = []
  rmq_version         = "3.9.13"
  vpc_id              = cloudamq_vpc.vpc.id
  keep_associated_vpc = true
}

# Additional VPC information
data "cloudamqp_vpc_info" "vpc_info" {
  vpc_id = cloudamqp_vpc.vpc.id
}
```

## Argument Reference

* `name`      - (Required) The name of the VPC.
* `region`    - (Required) The hosted region for the managed standalone VPC
* `subnet`    - (Required) The VPC subnet
* `tags`      - (Optional) Tag the VPC with optional tags

## Attributes Reference

All attributes reference are computed

* `id`       - The identifier for this resource.
* `vpc_name` - VPC name given when hosted at the cloud provider

## Import

`cloudamqp_vpc` can be imported using the CloudAMQP VPC identifier. To retrieve the identifier of an
instance, use [CloudAMQP API list vpcs].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_vpc.vpc
  id = <vpc_id>
}
```

Or with Terraform CLI:

`terraform import cloudamqp_vpc.vpc <vpc_id>`

[CloudAMQP API list vpcs]: https://docs.cloudamqp.com/#list-vpcs
[CloudAMQP plans]: https://www.cloudamqp.com/plans.html

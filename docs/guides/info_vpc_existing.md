---
layout: "cloudamqp"
page_title: "CloudAMQP: Create instance in existing VPC"
subcategory: "info"
description: |-
  Example of creating an instance in existing VPC.
---

# Create instance in existing VPC

If you already have a CloudAMQP instance with a VPC for a given `region`, it's possible to create another CloudAMQP instance in the same `region` and add this to the existing VPC. In order to do this, the internal CloudAMQP VPC indentifier (vpc_id) must been known and then added to new CloudAMQP instance.

## Example usage

To extract the existing VPC identifier, the CloudAMQP used to create the VPC can be loaded as an data source. Either by using the CloudAMQP instance name or identifier (instance_id).

```hcl
# by instance name
locals {
  instance_name = "<instance name>"
}
data "cloudamqp_account" "account" {}
data "cloudamqp_instance" "instance_info" {
  instance_id = [for instance in data.cloudamqp_account.account.instances : instance if instance["name"] == local.instance_name][0].id
}

# by instance ID
data "cloudamqp_instance" "instance_info" {
  instance_id = 123
}

output "vpc_id" {
  value = data.cloudamqp_instance.instance_info.vpc_id
}

output "vpc_subnet" {
  value = data.cloudamqp_instance.instance_info.vpc_subnet
}
```

From there use the extracted VPC identifier as `vpc_id` for the resource creating a new CloudAMQP instance.

```hcl
resource "cloudamqp_instance" "my_instance" {
  name       = "<instance name>"
  plan       = "squirrel-1"
  region     = "amazon-web-services::eu-north-1"
  vpc_id     = data.cloudamqp_instance.instance_info.vpc_id
}
```

Or by already knowing the internal CloudAMQP VPC identifier

```hcl
resource "cloudamqp_instance" "my_instance" {
  name       = "<instance name>"
  plan       = "squirrel-1"
  region     = "amazon-web-services::eu-north-1"
  vpc_id     = 123
}
```

You don't need to supply `vpc_subnet`, but if you supply both `vpc_id` and `vpc_subnet`. `vpc_id` has preference in the backend API (mismatch with the correct `vpc_subnet` value won't create a new VPC when `vpc_id` is supplied).

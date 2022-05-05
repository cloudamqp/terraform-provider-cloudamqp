---
layout: "cloudamqp"
page_title: "CloudAMQP: Managed VPC"
subcategory: "info"
description: |-
  Example of handle managed VPC, post v1.16.0
---

# Guide to handle managed VPC

From v1.16.0 it is possible to handle standalone VPC as a managed VPC resource.

## Managed VPC and dedicated instances

Create standalone VPC as a managed VPC resource.

```hcl
# Managed VPC
resource "cloudamqp_vpc" "vpc" {
  name   = "<vpc-name>"
  region = "amazon-web-services::us-east-1"
  subnet = "10.56.72.0/24"
  tags   = []
}
```
Create multiple instances and add them to the managed VPC.

```hcl
resource "cloudamqp_instance" "instance_01" {
  name                = "terraform-cloudamqp-instance-01"
  plan                = "squirrel-1"
  region              = "amazon-web-services::us-west-1"
  tags                = ["terraform"]
  rmq_version         = "3.9.14"
  vpc_id = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_instance" "instance_02" {
  name                = "terraform-cloudamqp-instance-02"
  plan                = "squirrel-1"
  region              = "amazon-web-services::us-west-1"
  tags                = ["terraform"]
  rmq_version         = "3.9.14"
  vpc_id = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}
```

Both instances can be deleted and the managed VPC can still be used.

## Dedicated instance and vpc_subnet

Creating dedicated instance with attribute vpc_subnet. This will both create an instance and a standalone VPC.

```hcl
# Dedicated instance with vpc_subnet also creates VPC
resource "cloudamqp_instance" "instance_01" {
  name                = "terraform-cloudamqp-instance-01"
  plan                = "squirrel-1"
  region              = "amazon-web-services::us-west-1"
  tags                = ["terraform"]
  rmq_version         = "3.9.14"
  vpc_subnet          = "10.56.72.0/24"
}
```

### Import managed VPC
Once the instance and the VPC are created, the VPC can be imported as a managed VPC.

```hcl
# Imported managed VPC
resource "cloudamqp_vpc" "vpc" { }
```
`terraform import cloudamqp_vpc.vpc <vpc_id>`

Add the correct info for the imported standalone VPC.

```hcl
# Imported standalone VPC as a managed VPC
resource "cloudamqp_vpc" "vpc" {
  name   = "<vpc-name>"
  region = "amazon-web-services::us-east-1"
  subnet = "10.56.72.0/24"
  tags   = []
}
```

### Update instance resource

```hcl
# Add vpc_id attribute
resource "cloudamqp_instance" "instance_01" {
  name                = "terraform-cloudamqp-instance-01"
  plan                = "squirrel-1"
  region              = "amazon-web-services::us-west-1"
  tags                = ["terraform"]
  rmq_version         = "3.9.14"
  vpc_subnet          = "10.56.72.0/24"
  vpc_id              = cloudamqp_vpc.vpc.id
}
```

Run `terraform apply -refresh-only` to update the state file with the correct data.

### Delete instance

When deleting the instance, the associated VPC will be deleted by default (if no other instances are added).

```hcl
resource "cloudamqp_instance" "instance_01" {
  name                = "terraform-cloudamqp-instance-01"
  plan                = "squirrel-1"
  region              = "amazon-web-services::us-west-1"
  tags                = ["terraform"]
  rmq_version         = "3.9.14"
  vpc_subnet          = "10.56.72.0/24"
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}
```

In order to keep the associated VPC the attribute `keep_associated_vpc` must be set to *true*.

Run `terraform apply` to update the state file with the correct data, then the instance can be deleted.

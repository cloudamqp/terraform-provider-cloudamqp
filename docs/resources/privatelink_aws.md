---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_privatelink_aws"
description: |-
  Enable PrivateLink for a CloudAMQP instance hosted in AWS.
---

# cloudamqp_privatelink_aws

Enable PrivateLink for a CloudAMQP instance hosted in AWS. If no existing VPC available when enable PrivateLink, a new VPC will be created with subnet `10.52.72.0/24`.

More information about [CloudAMQP Privatelink](https://www.cloudamqp.com/docs/cloudamqp-privatelink.html#aws-privatelink).

Only available for dedicated subscription plans.

## Example Usage

CloudAMQP instance without existing VPC

```hcl
resource "cloudamqp_instance" "instance" {
  name   = "Instance 01"
  plan   = "squirrel-1"
  region = "amazon-web-services::us-west-1"
  tags   = ["test"]
  rmq_version = "3.10.8"
}

resource "cloudamqp_privatelink_aws" "privatelink" {
  instance_id = cloudamqp_instance.instance.id
  allowed_principals = [
    "arn:aws:iam::aws-account-id:user/user-name"
  ]
}
```

CloudAMQP instance already in an existing VPC.

```hcl
resource "cloudamqp_vpc" "vpc" {
  name = "Standalone VPC"
  region = "amazon-web-services::us-west-1"
  subnet = "10.56.72.0/24"
  tags = ["test"]
}

resource "cloudamqp_instance" "instance" {
  name   = "Instance 01"
  plan   = "squirrel-1"
  region = "amazon-web-services::us-west-1"
  tags   = ["test"]
  rmq_version = "3.10.8"
  vpc_id = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_privatelink_aws" "privatelink" {
  instance_id = cloudamqp_instance.instance.id
  allowed_principals = [
    "arn:aws:iam::aws-account-id:user/user-name"
  ]
}
```

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier.
* `allowed_principals` - (Required) Allowed principals to access the endpoint service.
* `sleep` - (Optional) Configurable sleep time (seconds) when enable PrivateLink. Default set to 60 seconds.
* `timeout` - (Optional) - Configurable timeout time (seconds) when enable PrivateLink. Default set to 3600 seconds.

Allowed principals format: <br>
`arn:aws:iam::aws-account-id:root` <br>
`arn:aws:iam::aws-account-id:user/user-name` <br>
`arn:aws:iam::aws-account-id:role/role-name`

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.
* `status`- PrivateLink status [enable, pending, disable]
* `service_name` - Service name of the PrivateLink used when creating the endpoint from other VPC.
* `active_zones` - Covering availability zones used when creating an Endpoint from other VPC.

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_privatelink_aws` can be imported using CloudAMQP internal identifier.

`terraform import cloudamqp_privatelink_aws.privatelink <id>`

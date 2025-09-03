---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_account_action"
description: |-
  Invoke account action
---

# cloudamqp_account_action

This resource allows you to invoke an account action. Current supported actions are

* Rotate password for RabbitMQ/LavinMQ user
* Rotate API key for the CloudAMQP instance
* Enable VPC feature

## Example Usage

Invoke one of the actions with `terraform apply`

```hcl
resource "cloudamqp_account_action" "rotate-password" {
  instance_id = cloudamqp_instance.instance.id
  action      = "rotate-password"
} 
```

```hcl
resource "cloudamqp_account_action" "rotate-apikey" {
  instance_id = cloudamqp_instance.instance.id
  action      = "rotate-apikey"
} 
```

```hcl
resource "cloudamqp_account_actions" "enable_vpc" {
  instance_id = cloudamqp_instance.instance.id
  action      = "enable-vpc"
}
```

<details>
 <summary>
    <b>
      <i>Manage the enabled VPC</i>
    </b>
  </summary>

To add the enable VPC to a managed standalone VPC.

First fetch the VPC identifier <id>

1. Run `terraform refresh` the `vpc_id` will be added to the state for the `cloudamqp_instance.instance` resource.
2. Retrieve the `vpc_id` form the CloudAMQP HTTP API. Either via [list-instances] or [list-vpcs].

```hcl
import {
  to = cloudamqp_vpc.vpc
  id = <id>
}

resource "cloudamqp_vpc" "vpc" {
  name    = "enable-vpc-feature"
  region  = "amazon-web-services::us-east-1"
  subnet  = "10.56.72.0/24"
  tags    = []
}

resource "cloudamqp_instance" "instance" {
  name                = "enable-vpc-feature"
  plan                = "penguin-1"
  region              = "amazon-web-services::us-east-1"
  tags                = []
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_account_actions" "enable_vpc" {
  instance_id = cloudamqp_instance.instance.id
  action      = "enable-vpc"
}
```

</details>

After the action have been invoked, the state need to be refreshed to get the latest changes in
***cloudamqp_instance*** or ***data.cloudamqp_instance***. This can be done with
`terraform refresh`.

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance ID.
* `action`      - (Required/ForceNew) The action to be invoked. Allowed actions
                  `rotate-password`, `rotate-apikey`, `enable-vpc`.

### Actions

**rotate-password**
Initiate rotation of the user password on your instance.

**rotate-apikey**
Initiate rotation of the instance API key used for the CloudAMQP [HTTP API].

**enable-vpc**
Enables the VPC feature on existing instance not using standalone VPC. Extra cost will be applied:
https://www.cloudamqp.com/plans.html#xtr

-> NOTE: This action is irreversible, if you want to disable VPC features you will need to delete the instance and create a new one.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

Not possible to import this resource.

[list-instances]: https://docs.cloudamqp.com/#list-instances
[list-vpcs]: https://docs.cloudamqp.com/#list-vpcs
[HTTP API]: https://docs.cloudamqp.com/cloudamqp_api.html

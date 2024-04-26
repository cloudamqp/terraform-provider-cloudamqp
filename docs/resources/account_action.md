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

## Example Usage

Invoke one of the actions with `terraform apply`

```hcl
resource "cloudamqp_account_action" "rotate-password" {
  instance_id = cloudamqp_instance.instance.id
  action = "rotate-password"
} 
```

```hcl
resource "cloudamqp_account_action" "rotate-apikey" {
  instance_id = cloudamqp_instance.instance.id
  action = "rotate-apikey"
} 
```

After the action have been invoked, the state need to be refreshed to get the latest changes in ***cloudamqp_instance*** or ***data.cloudamqp_instance***. This can be done with `terraform refresh`.

## Argument Reference

The following arguments are supported:

* `instance_id`         - (Required) The CloudAMQP instance ID.
* `action`              - (Required/ForceNew) The action to be invoked. Allowed actions `rotate-password`, `rotate-apikey`.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

Not possible to import this resource.

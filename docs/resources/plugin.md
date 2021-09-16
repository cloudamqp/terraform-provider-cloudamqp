---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_plugin"
description: |-
  Enable and disable Rabbit MQ plugin.
---

# cloudamqp_plugin

This resource allows you to enable or disable Rabbit MQ plugins.

Only available for dedicated subscription plans.

**Note: Known issue when changing multiple plugins:**
Rabbit MQ can only change one plugin at a time. Therefore if multiple plugins should be changed, Terraform needs to be changed from working in parallel to be working in sequence. This can be done by adding parallelism (default set to 10) command to apply.

```
terraform apply -parallelism=1
```

## Example Usage

```hcl
resource "cloudamqp_plugin" "plugin_rabbitmq_top" {
  instance_id = cloudamqp_instance.instance.id
  name = "rabbitmq_top"
  enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance ID.
* `name`        - (Required) The name of the Rabbit MQ plugin.
* `enabled`     - (Required) Enable or disable the plugins.

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_plugin` can be imported using the name argument of the resource together with CloudAMQP instance identifier. The name and identifier are CSV separated, see example below.

`terraform import cloudamqp_plugin.rabbitmq_management rabbitmq_management,<instance_id>`

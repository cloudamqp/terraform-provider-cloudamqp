---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_plugin"
description: |-
  Enable and disable Rabbit MQ plugin.
---

# cloudamqp_plugin

This resource allows you to enable or disable Rabbit MQ plugins.

Only available for dedicated subscription plans.

⚠️  From our go API wrapper [v1.4.0](https://github.com/84codes/go-api/releases/tag/v1.4.0) there is support for multiple retries when requesting information about plugins. This was introduced to avoid `ReadPlugin error 400: Timeout talking to backend`.

**Enable multiple plugins:** Rabbit MQ can only change one plugin at a time. It will fail if multiple plugins resources are used, unless by creating dependencies with `depend_on` between the resources. Once one plugin has been enabled, the other will continue. See example below.

## Example Usage

```hcl
resource "cloudamqp_plugin" "rabbitmq_top" {
  instance_id = cloudamqp_instance.instance.id
  name = "rabbitmq_top"
  enabled = true
}
```

**Enable multiple plugins**
```hcl
resource "cloudamqp_plugin" "rabbitmq_top" {
  instance_id = cloudamqp_instance.instance.id
  name = "rabbitmq_top"
  enabled = true
}

resource "cloudamqp_plugin" "rabbitmq_amqp1_0" {
  instance_id = cloudamqp_instance.instance.id
  name = "rabbitmq_amqp1_0"
  enabled = true

  depends_on = [
    cloudamqp_plugin.rabbitmq_top
  ]
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

If multiple plugins should be enable, create dependencies between the plugin resources. See example above.

## Import

`cloudamqp_plugin` can be imported using the name argument of the resource together with CloudAMQP instance identifier. The name and identifier are CSV separated, see example below.

`terraform import cloudamqp_plugin.rabbitmq_management rabbitmq_management,<instance_id>`

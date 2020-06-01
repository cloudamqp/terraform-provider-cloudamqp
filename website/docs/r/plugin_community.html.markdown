---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_plugin_commiunity"
description: |-
  Install or uninstall community plugin.
---

# cloudamqp_plugin_community

This resource allows you to install or uninstall community plugins. Once installed the plugin gets available in `cloudamqp_plugin`.

Only available for dedicated subscription plans.

## Example Usage

```hcl
resource "cloudamqp_plugin_community" "rabbitmq_delayed_message_exchange" {
  instance_id = cloudamqp_instance.instance_01.id
  name = "rabbitmq_delayed_message_exchange"
  enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance ID.
* `name`        - (Required) The name of the Rabbit MQ plugin.
* `enabled`     - (Required) Enable or disable the plugins.


## Import

`cloudamqp_plugin` can be imported using the name argument of the resource together with CloudAMQP instance identifier. The name and identifier are CSV separated, see example below.

`terraform import cloudamqp_plugin.delayed_message_exchange rabbitmq_delayed_message_exchange,<instance_id>`

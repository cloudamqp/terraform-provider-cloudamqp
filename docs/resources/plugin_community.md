---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_plugin_commiunity"
description: |-
  Install or uninstall community plugin.
---

# cloudamqp_plugin_community

This resource allows you to install or uninstall community plugins. Once installed the plugin will be available in `cloudamqp_plugin`.

Only available for dedicated subscription plans.

⚠️  From our go API wrapper [v1.5.0](https://github.com/84codes/go-api/releases/tag/v1.5.0) there is support for multiple retries when requesting information about community plugins. This was introduced to avoid `ReadPluginCommunity error 400: Timeout talking to backend`.

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
* `name`        - (Required) The name of the Rabbit MQ community plugin.
* `enabled`     - (Required) Enable or disable the plugins.

## Attributes Reference

* `id`  - The identifier for this resource`

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_plugin` can be imported using the name argument of the resource together with CloudAMQP instance identifier. The name and identifier are CSV separated, see example below.

`terraform import cloudamqp_plugin.<resource_name> <plugin_name>,<instance_id>`

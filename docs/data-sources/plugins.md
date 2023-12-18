---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_plugins"
description: |-
  Get information installed and available plugins.
---

# cloudamqp_plugins

Use this data source to retrieve information about installed and available plugins for the CloudAMQP instance.

## Example Usage

```hcl
data "cloudamqp_plugins" "plugins" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attributes reference

All attributes reference are computed

* `id`      - The identifier for this resource.
* `plugins` - An array of plugins. Each `plugins` block consists of the fields documented below.
* `sleep` - (Optional) Configurable sleep time (seconds) for retries when requesting information
about plugins. Default set to 10 seconds. *Available from v1.29.0*
* `timeout` - (Optional) - Configurable timeout time (seconds) for retries when requesting
information about plugins. Default set to 1800 seconds. *Available from v1.29.0*

___

The `plugins` block consist of

* `name`        - The type of the recipient.
* `version`     - Rabbit MQ version that the plugins are shipped with.
* `description` - Description of what the plugin does.
* `enabled`     - Enable or disable information for the plugin.

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

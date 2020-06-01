---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_plugins"
description: |-
  Get information installed and available plugins.
---

# cloudamqp_plugins

Use this data source to retrieve information about installed and available plugins for the CloudAMQP instance. Require to know the identifier of the corresponding `cloudamqp_instance`resource or data source.

## Example Usage

```hcl
data "cloudamqp_plugins" "plugins" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attribute reference

* `plugins` - (Computed) An array of plugins. Each `plugin` block consists of the fields documented below.

  * `name`        - (Computed) The type of the recipient.
  * `version`     - (Computed) Rabbit MQ version that the plugins are shipped with.
  * `description` - (Computed) Description of what the plugin does.
  * `enabled`     - (Computed) Enable or disable information for the plugin.

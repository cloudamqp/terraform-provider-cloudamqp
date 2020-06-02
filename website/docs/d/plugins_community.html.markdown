---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_plugins_community"
description: |-
  Get information about available community plugins.
---

# cloudamqp_plugins_community

Use this data source to retrieve information about available community plugins for the CloudAMQP instance. Depens on the identifier of the corresponding `cloudamqp_instance`resource or data source.

## Example Usage

```hcl
data "cloudamqp_plugins_community" "communit_plugins" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attribute reference

* `plugins` - (Computed) An array of community plugins. Each `plugins` block consists of the fields documented below.

___

The `plugins` block consists of

* `name`        - (Computed) The type of the recipient.
* `require`     - (Computed) Min. required Rabbit MQ version to be used.
* `description` - (Computed) Description of what the plugin does.

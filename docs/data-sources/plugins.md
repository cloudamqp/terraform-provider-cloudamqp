---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_plugins"
description: |-
  Get information installed and available plugins.
---

# cloudamqp_plugins

Use this data source to retrieve information about installed and available plugins for the CloudAMQP instance.

⚠️  From our go API wrapper [v1.4.0](https://github.com/84codes/go-api/releases/tag/v1.4.0) there is support for multiple retries when requesting information about plugins. This was introduced to avoid `ReadPlugin error 400: Timeout talking to backend`.

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

___

The `plugins` block consist of

* `name`        - The type of the recipient.
* `version`     - Rabbit MQ version that the plugins are shipped with.
* `description` - Description of what the plugin does.
* `enabled`     - Enable or disable information for the plugin.

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

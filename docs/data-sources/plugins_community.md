---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_plugins_community"
description: |-
  Get information about available community plugins.
---

# cloudamqp_plugins_community

Use this data source to retrieve information about available community plugins for the CloudAMQP instance.

⚠️  From our go API wrapper [v1.5.0](https://github.com/84codes/go-api/releases/tag/v1.5.0) there is support for multiple retries when requesting information about community plugins. This was introduced to avoid `ReadPluginCommunity error 400: Timeout talking to backend`.

## Example Usage

```hcl
data "cloudamqp_plugins_community" "communit_plugins" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attributes reference

All attributes reference are computed

* `id`      - The identifier for this resource.
* `plugins` - An array of community plugins. Each `plugins` block consists of the fields documented below.

___

The `plugins` block consists of

* `name`        - The type of the recipient.
* `require`     - Min. required Rabbit MQ version to be used.
* `description` - Description of what the plugin does.

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

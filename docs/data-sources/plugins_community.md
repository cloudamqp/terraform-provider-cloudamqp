---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_plugins_community"
description: |-
  Get information about available community plugins.
---

# cloudamqp_plugins_community

Use this data source to retrieve information about available community plugins for the CloudAMQP
instance.

## Example Usage

```hcl
data "cloudamqp_plugins_community" "communit_plugins" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier.
* `sleep`       - (Optional) Configurable sleep time (seconds) for retries when requesting
                  information about community plugins. Default set to 10 seconds.
* `timeout`     - (Optional) Configurable timeout time (seconds) for retries when requesting
                  information about community plugins. Default set to 1800 seconds.

## Attributes Reference

All attributes reference are computed

* `id`      - The identifier for this resource.
* `plugins` - An array of community plugins. Each `plugins` block consists of the fields documented
              below.

___

The `plugins` block consists of

* `name`        - The type of the recipient.
* `require`     - Min. required Rabbit MQ version to be used.
* `description` - Description of what the plugin does.

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

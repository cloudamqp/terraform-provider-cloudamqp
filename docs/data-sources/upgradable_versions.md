---
layout: "cloudamqp"
page_title: "CloudAMQP: data source upgradable_versions"
description: |-
  Get information of upgradable versions for RabbitMQ and Erlang.
---

# cloudamqp_plugins

Use this data source to retrieve information about possible upgradable versions for RabbitMQ and Erlang.

## Example Usage

```hcl
data "cloudamqp_upgradable_versions" "versions" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attributes reference

All attributes reference are computed

* `new_rabbitmq_version`  - Possible upgradable version for RabbitMQ.
* `new_erlang_version`    - Possible upgradable version for Erlang.

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

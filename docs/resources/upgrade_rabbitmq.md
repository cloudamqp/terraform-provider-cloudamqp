---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_upgrade_rabbitmq"
description: |-
  Invoke upgrade to latest possible upgradable versions for RabbitMQ and Erlang.
---

# cloudamqp_upgrade_rabbitmq

This resource allows you to automatically upgrade to latest possible upgradable versions for RabbitMQ and Erlang. Depending on initial versions of RabbitMQ and Erlang of the CloudAMQP instance. Multiple runs may be needed to get to latest versions. (E.g. after completed upgrade, check data source `cloudamqp_upgradable_versions` to see if newer versions is available. Then delete `cloudamqp_upgrade_rabbitmq` and create it again to invoke the upgrade.

> :warning: **WARNING: Before using this resource.**
>
> Auto delete queues (queues that are marked AD) will be deleted during the update.
>
> Any custom plugins support has installed on your behalf will be disabled and you need to contact support@cloudamqp.com and ask to have them re-installed.
>
> TLS 1.0 and 1.1 will not be supported after the update.

Only available for dedicated subscription plans.

## Example Usage

```hcl
# Retrieve latest possible upgradable versions for RabbitMQ and Erlang
data "cloudamqp_upgradable_versions" "versions" {
  instance_id = cloudamqp_instance.instance.id
}

# Invoke automatically upgrade to latest possible upgradable versions for RabbitMQ and Erlang
resource "cloudamqp_upgrade_rabbitmq" "upgrade" {
  instance_id = cloudamqp_instance.instance.id
}
```

```hcl
# Retrieve latest possible upgradable versions for RabbitMQ and Erlang
data "cloudamqp_upgradable_versions" "versions" {
  instance_id = cloudamqp_instance.instance.id
}

# Delete the resource
# resource "cloudamqp_upgrade_rabbitmq" "upgrade" {
#   instance_id = cloudamqp_instance.instance.id
# }
```

If newer version is still available to be upgradable in the data source, re-run again.

```hcl
# Retrieve latest possible upgradable versions for RabbitMQ and Erlang
data "cloudamqp_upgradable_versions" "versions" {
  instance_id = cloudamqp_instance.instance.id
}

# Invoke automatically upgrade to latest possible upgradable versions for RabbitMQ and Erlang
resource "cloudamqp_upgrade_rabbitmq" "upgrade" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier

## Import

Not possible to import this resource.

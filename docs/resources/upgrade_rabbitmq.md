---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_upgrade_rabbitmq"
description: |-
  Invoke upgrade to latest possible upgradable versions for RabbitMQ and Erlang.
---

# cloudamqp_upgrade_rabbitmq

This resource allows you to automatically upgrade to the latest possible upgradable versions for RabbitMQ and Erlang. Depending on initial versions of RabbitMQ and Erlang of the CloudAMQP instance, multiple runs may be needed to get to the latest versions. After completed upgrade, check data source `cloudamqp_upgradable_versions` to see if newer versions is available. Then delete `cloudamqp_upgrade_rabbitmq` and create it again to invoke the upgrade.

> **Important Upgrade Information**
> - All single node upgrades will require some downtime since RabbitMQ needs a restart.
> - From RabbitMQ version 3.9, rolling upgrades between minor versions (e.g. 3.9 to 3.10), in a multi-node cluster are possible without downtime. This means that one node is upgraded at a time while the other nodes are still running. For versions older than 3.9, patch version upgrades (e.g. 3.8.x to 3.8.y) are possible without downtime in a multi-node cluster, but minor version upgrades will require downtime. 
> - Auto delete queues (queues that are marked AD) will be deleted during the update.
> - Any custom plugins support has installed on your behalf will be disabled and you need to contact support@cloudamqp.com and ask to have them re-installed.
> - TLS 1.0 and 1.1 will not be supported after the update.

Only available for dedicated subscription plans running ***RabbitMQ***.

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

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance identifier

## Import

Not possible to import this resource.

---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_upgrade_rabbitmq"
description: |-
  Invoke upgrade to latest possible upgradable versions for RabbitMQ and Erlang.
---

# cloudamqp_upgrade_rabbitmq

This resource allows you to upgrade RabbitMQ version. Depending on initial versions of RabbitMQ and Erlang of the CloudAMQP instance, multiple runs may be needed to get to the latest or wanted version. Reason for this is certain supported RabbitMQ version will also automatically upgrade Erlang version.

There is three different ways to trigger the version upgrade

> - Specify RabbitMQ version to upgrade to
> - Upgrade to latest RabbitMQ version
> - Old behaviour to upgrade to latest RabbitMQ version

See, below example usage for the difference.

Only available for dedicated subscription plans running ***RabbitMQ***.

## Example Usage

<details>
  <summary>
    <b>
      <i>Specify version upgrade, from v1.31.0</i>
    </b>
  </summary>

Specify the version to upgrade to. List available upgradable versions, use [CloudAMQP API](https://docs.cloudamqp.com/cloudamqp_api.html#get-available-versions).
After the upgrade finished, there can still be newer versions available.

```hcl
resource "cloudamqp_instance" "instance" {
  name        = "rabbitmq-version-upgrade-test"
  plan        = "bunny-1"
  region      = "amazon-web-services::us-west-1"
}

resource "cloudamqp_upgrade_rabbitmq" "upgrade" {
  instance_id     = cloudamqp_instance.instance.id
  new_version     = "3.13.2"
}
```

</details>

<details>
  <summary>
    <b>
      <i>Upgrade to latest possible version, from v1.31.0</i>
    </b>
  </summary>

This will upgrade RabbitMQ to the latest possible version detected by the data source `cloudamqp_upgradable_versions`.
Multiple runs can be needed to upgrade the version even further.

```hcl
resource "cloudamqp_instance" "instance" {
  name        = "rabbitmq-version-upgrade-test"
  plan        = "bunny-1"
  region      = "amazon-web-services::us-west-1"
}

data "cloudamqp_upgradable_versions" "upgradable_versions" {
  instance_id = cloudamqp_instance.instance.id
}

resource "cloudamqp_upgrade_rabbitmq" "upgrade" {
  instance_id     = cloudamqp_instance.instance.id
  current_version = cloudamqp_instance.instance.rmq_version
  new_version     = data.cloudamqp_upgradable_versions.upgradable_versions.new_rabbitmq_version
}
```

</details>

<details>
  <summary>
    <b>
      <i>Upgrade to latest possible version, before v1.31.0</i>
    </b>
  </summary>

Old behaviour of the upgrading the RabbitMQ version. No longer recommended.

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

</details>

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance identifier
* `current_version` - (Optional) Helper argument to change upgrade behaviour to latest possible version
* `new_version` - (Optional/ForceNew) The new version to upgrade to

## Import

Not possible to import this resource.

## Important Upgrade Information

> - All single node upgrades will require some downtime since RabbitMQ needs a restart.
> - From RabbitMQ version 3.9, rolling upgrades between minor versions (e.g. 3.9 to 3.10), in a multi-node cluster are possible without downtime. This means that one node is upgraded at a time while the other nodes are still running. For versions older than 3.9, patch version upgrades (e.g. 3.8.x to 3.8.y) are possible without downtime in a multi-node cluster, but minor version upgrades will require downtime. 
> - Auto delete queues (queues that are marked AD) will be deleted during the update.
> - Any custom plugins support has installed on your behalf will be disabled and you need to contact support@cloudamqp.com and ask to have them re-installed.
> - TLS 1.0 and 1.1 will not be supported after the update.

## Multiple runs

Depending on initial versions of RabbitMQ and Erlang of the CloudAMQP instance, multiple runs may be needed to get to the latest or wanted version.

Example steps needed when starting at RabbitMQ version 3.12.2

|  Version         | Supported upgrading versions              | Min version to upgrade Erlang |
|------------------|-------------------------------------------|-------------------------------|
| 3.12.2           | 3.12.4, 3.12.6, 3.12.10, 3.12.12, 3.12.13 | 3.12.13                       |
| 3.12.13          | 3.13.2                                    | 3.13.2                        |
| 3.13.2           | -                                         | -                             |

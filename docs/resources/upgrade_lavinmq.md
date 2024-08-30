---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_upgrade_lavinmq"
description: |-
  Invoke upgrade to latest possible upgradable versions for LavinMQ.
---

# cloudamqp_upgrade_lavinmq

This resource allows you to upgrade LavinMQ version. 

There is two different ways to trigger the version upgrade

> - Specify LavinMQ version to upgrade to
> - Upgrade to latest LavinMQ version

See, below example usage for the difference.

Only available for dedicated subscription plans running ***LavinMQ***.

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
  name        = "lavinmq-version-upgrade-test"
  plan        = "lynx-1"
  region      = "amazon-web-services::us-west-1"
}

resource "cloudamqp_upgrade_lavinmq" "upgrade" {
  instance_id     = cloudamqp_instance.instance.id
  new_version     = "1.3.1"
}
```

</details>

<details>
  <summary>
    <b>
      <i>Upgrade to latest possible version, from v1.31.0</i>
    </b>
  </summary>

This will upgrade LavinMQ to the latest possible version detected by the data source `cloudamqp_upgradable_versions`.

```hcl
resource "cloudamqp_instance" "instance" {
  name        = "lavinmq-version-upgrade-test"
  plan        = "lynx-1"
  region      = "amazon-web-services::us-west-1"
}

data "cloudamqp_upgradable_versions" "upgradable_versions" {
  instance_id = cloudamqp_instance.instance.id
}

resource "cloudamqp_upgrade_lavinmq" "upgrade" {
  instance_id     = cloudamqp_instance.instance.id
  current_version = cloudamqp_instance.instance.rmq_version
  new_version     = data.cloudamqp_upgradable_versions.upgradable_versions.new_lavinmq_version
}
```

</details>


## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance identifier
* `new_version` - (Optional/ForceNew) The new version to upgrade to

## Import

Not possible to import this resource.

## Important Upgrade Information

> - All single node upgrades will require some downtime since LavinMQ needs a restart.
> - Auto delete queues (queues that are marked AD) will be deleted during the update.

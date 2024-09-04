---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_upgrade_lavinmq"
description: |-
  Invoke upgrade for LavinMQ.
---

# cloudamqp_upgrade_lavinmq

This resource allows you to upgrade LavinMQ version. 

See below example usage.

Only available for dedicated subscription plans running ***LavinMQ***.

## Example Usage

<details>
  <summary>
    <b>
      <i>Upgrade LavinMQ, specify which version to upgrade to, from v1.31.0</i>
    </b>
  </summary>

Specify the version to upgrade to. List available upgradable versions, use [CloudAMQP API](https://docs.cloudamqp.com/cloudamqp_api.html#get-available-versions).

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


## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance identifier
* `new_version` - (Optional/ForceNew) The new version to upgrade to

## Import

Not possible to import this resource.

## Important Upgrade Information

> - All single node upgrades will require some downtime since LavinMQ needs a restart.
> - Auto delete queues (queues that are marked AD) will be deleted during the update.

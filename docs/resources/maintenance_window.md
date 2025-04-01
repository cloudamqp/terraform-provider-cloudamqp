---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_maintenance_window"
description: |-
  Update the preferred maintenance window.
---

# cloudamqp_maintenance_window

This resource allows you to set the preferred start of new scheduled maintenances. The maintenance
windows are 3 hours long and CloudAMQP attempts to begin the maintenance as close as possible to the
preferred start. A maintenance will never start before the window.

Available for dedicated subscription plans.

## Example Usage

<details>
  <summary>
    <b>Set the preferred maintenance start</b>
  </summary>

```hcl
resource "cloudamqp_maintenance_window" "this" {
  instance_id     = cloudamqp_instance.instance.id
  preferred_day   = "Monday"
  preferred_time  = "23:00"
}
```

</details>

<details>
  <summary>
    <b>Set the preferred maintenance start with automatic updates</b>
  </summary>

When setting the automatic updates to "on", a maintenance for version update will be scheduled once
a new LavinMQ version been released.

```hcl
resource "cloudamqp_maintenance_window" "this" {
  instance_id       = cloudamqp_instance.instance.id
  preferred_day     = "Monday"
  preferred_time    = "23:00"
  automatic_updates = "on"
}
```

</details>

<details>
  <summary>
    <b>Only set preferred time of day</b>
  </summary>

```hcl
resource "cloudamqp_maintenance_window" "this" {
  instance_id     = cloudamqp_instance.instance.id
  preferred_time  = "23:00"
}
```

</details>

<details>
  <summary>
    <b>Only set preferred day of week</b>
  </summary>

```hcl
resource "cloudamqp_maintenance_window" "this" {
  instance_id   = cloudamqp_instance.instance.id
  preferred_day = "Monday"
}
```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id`       - (Required) The CloudAMQP instance ID.
* `preferred_day`     - (Optional) Preferred day of the week when to schedule maintenance.
* `preferred_time`    - (Optional) Preferred time (UTC) of the day when to schedule maintenance.
* `automatic_updates` - (Optional/Computed) Allow scheduling of a maintenance for version update
                        once a new LavinMQ version been released.

### Valid preferred days

Valid days for preferred days:

["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]

### Valid preferred time

All times are in UTC and follow `hh:mm` [format].

Example: "00:00", "06:00", "12:00", "18:00", "23:00"

### Valid automatic updates

This argument can right now only be used for LavinMQ subscriptions plan. For new instances it's
automatically set to "off".

Valid values are: ["on", "off"]

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource (will be same as instance identifier).

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_maintenance_window` can be imported using CloudAMQP instance identifier. To retrieve the
identifier of an instance, use [CloudAMQP API list instances].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_maintenance_window.this
  id = cloudamqp_instance.instance.id
}
```

Or with Terraform CLI:

`terraform import cloudamqp_maintenance_window.this <id>`

[CloudAMQP API list instances]: https://docs.cloudamqp.com/#list-instances
[format]: https://developer.hashicorp.com/terraform/language/functions/formatdate#specification-syntax

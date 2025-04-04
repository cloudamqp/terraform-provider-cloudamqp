---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_alarm"
description: |-
  Creates and manages alarms to trigger notifications.
---

# cloudamqp_alarm

This resource allows you to create and manage alarms to trigger based on a set of conditions. Once
triggerd a notification will be sent to the assigned recipients. When creating a new instance, there
will also be a set of default alarms (cpu, memory and disk) created. All default alarms uses the
default recipient for notifications.

By setting `no_default_alarms` to *true* in `cloudamqp_instance`. This will create the instance
without default alarms and avoid the need to import them to get full control.

Available for all subscription plans, but `lemur`and `tiger`are limited to fewer alarm types. The
limited types supported can be seen in the table below in [Alarm Type Reference].

## Example Usage

<details>
  <summary>
    <b>
      <i>Basic example of CPU and memory alarm</i>
    </b>
  </summary>

```hcl
# New recipient
resource "cloudamqp_notification" "recipient_01" {
  instance_id = cloudamqp_instance.instance.id
  type        = "email"
  value       = "alarm@example.com"
  name        = "alarm"
}

# New cpu alarm
resource "cloudamqp_alarm" "cpu_alarm" {
  instance_id       = cloudamqp_instance.instance.id
  type              = "cpu"
  enabled           = true
  reminder_interval = 600
  value_threshold   = 95
  time_threshold    = 600
  recipients        = [cloudamqp_notification.recipient_01.id]
}

# New memory alarm
resource "cloudamqp_alarm" "memory_alarm" {
  instance_id       = cloudamqp_instance.instance.id
  type              = "memory"
  enabled           = true
  reminder_interval = 600
  value_threshold   = 95
  time_threshold    = 600
  recipients        = [cloudamqp_notification.recipient_01.id]
}
```

</details>

<details>
  <summary>
    <b>
      <i>Manage notice alarm, available from</i>
      <a href="https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.29.5">v1.29.5</a>
    </b>
  </summary>

Only one notice alarm can exists and cannot be created, instead the alarm resource will be updated.

```hcl
# New recipient
resource "cloudamqp_notification" "recipient_01" {
  instance_id = cloudamqp_instance.instance.id
  type        = "email"
  value       = "alarm@example.com"
  name        = "alarm"
}

# Update existing notice alarm
resource "cloudamqp_alarm" "notice" {
  instance_id = cloudamqp_instance.instance.id
  type        = "notice"
  enabled     = true
  recipients  = [cloudamqp_notification.recipient_01.id]
}
```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id`       - (Required) The CloudAMQP instance ID.
* `type`              - (Required) The alarm type, see valid options below.
* `enabled`           - (Required) Enable or disable the alarm to trigger.
* `reminder_interval` - (Optional) The reminder interval (in seconds) to resend the alarm if not
                        resolved. Set to 0 for no reminders. The Default is 0.
* `value_threshold`   - (Optional) The value to trigger the alarm for.
* `time_threshold`    - (Optional) The time interval (in seconds) the `value_threshold` should be
                        active before triggering an alarm.
* `queue_regex`       - (Optional) Regex for which queue to check.
* `vhost_regex`       - (Optional) Regex for which vhost to check
* `recipients`        - (Optional) Identifier for recipient to be notified. Leave empty to notify
                        all recipients.
* `message_type`      - (Optional) Message type `(total, unacked, ready)` used by queue alarm type.

Specific argument for `disk` alarm

* `value_calculation` - (Optional) Disk value threshold calculation, `fixed, percentage` of disk
                        space remaining.

Based on alarm type, different arguments are flagged as required or optional.

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Alarm type reference

Supported alarm types: `cpu, memory, disk, queue, connection, flow, consumer, netsplit,
  server_unreachable, notice`

Required arguments for all alarms: `instance_id, type, enabled`<br>
Optional argument for all alarms: `tags, queue_regex, vhost_regex`

| Name | Type | Shared | Dedicated | Required arguments |
| ---- | ---- | ---- | ---- | ---- |
| CPU | cpu | - | &#10004; | time_threshold, value_threshold |
| Memory | memory | - | &#10004;  | time_threshold, value_threshold |
| Disk space | disk | - | &#10004;  | time_threshold, value_threshold |
| Queue | queue | &#10004;  | &#10004; | time_threshold, value_threshold, queue_regex, vhost_regex, message_type |
| Connection | connection | &#10004; | &#10004; | time_threshold, value_threshold |
| Connection flow | flow | &#10004; | &#10004; | time_threshold, value_threshold |
| Consumer | consumer | &#10004; | &#10004; | time_threshold, value_threshold, queue, vhost |
| Netsplit | netsplit | - | &#10004; | time_threshold |
| Server unreachable | server_unreachable  | - | &#10004;  | time_threshold |
| Notice | notice | &#10004; | &#10004; | |

<br>

~> Notice alarm is manadatory! Only one can exists and cannot be deleted. Setting `no_default_alarm`
to true, will still create this alarm. See updated changes to [notice alarm] below.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_alarm` can be imported using the resource identifier together with the CloudAMQP instance
identifier (CSV separated). To retrieve the resource identifier, use [CloudAMQP API list alarms].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_alarm.alarm
  id = format("<id>,%s", cloudamqp_instance.instance.id)
}
```

Or use Terraform CLI:

`terraform import cloudamqp_alarm.alarm <id>,<instance_id>`

## Notice alarm

There is a limitation for notice alarm in the API backend. This alarm is mandatory, multiple
alarms cannot exists or be deleted.

From provider version [v1.29.5] it's possible to manage the notice alarm and no longer needs to be
imported. Just create the alarm resource as usually and it will be updated with given recipients.
If the alarm is deleted it will only be removed from the state file, but will still be enabled in
the backend.

[Alarm Type Reference]: #alarm-type-reference
[CloudAMQP API list alarms]: https://docs.cloudamqp.com/cloudamqp_api.html#list-alarms
[notice alarm]: #notice-alarm
[v1.29.5]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.29.5

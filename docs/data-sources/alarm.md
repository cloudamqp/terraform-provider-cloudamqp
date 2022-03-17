---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_alarm"
description: |-
  Get information on pre-defined or created alarms.
---

# cloudamqp_alarm

Use this data source to retrieve information about default or created alarms. Either use `alarm_id` or `type` to retrieve the alarm.

## Example Usage

```hcl
data "cloudamqp_alarm" "default_cpu_alarm" {
  instance_id = cloudamqp_instance.instance.id
  type = "cpu"
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.
* `alarm_id`    - (Optional) The alarm identifier. Either use this or `type` to give `cloudamqp_alarm` necessary information to retrieve the alarm.
* `type`        - (Optional) The alarm type. Either use this or `alarm_id` to give `cloudamqp_alarm` necessary information when retrieve the alarm. Supported [alarm types](#alarm-types)

## Attributes reference

All attributes reference are computed

* `id`                  - The identifier for this resource.
* `enabled`             - Enable/disable status of the alarm.
* `value_threshold`     - The value threshold that triggers the alarm.
* `reminder_internval`  - The reminder interval (in seconds) to resend the alarm if not resolved. Leave empty or set to 0 to not receive any reminders.
* `time_threshold`      - The time interval (in seconds) the `value_threshold` should be active before trigger an alarm.
* `queue_regex`         - Regular expression for which queue to check.
* `vhost_regex`         - Regular expression for which vhost to check
* `recipients`          - Identifier for recipient to be notified.
* `message_type`        - Message type `(total, unacked, ready)` used by queue alarm type.

Specific attribute for `disk` alarm

* `value_calculation`   - Disk value threshold calculation, `(fixed, percentage)` of disk space remaining.

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Alarm types

`cpu, memory, disk, queue, connection, consumer, netsplit, server_unreachable, notice`

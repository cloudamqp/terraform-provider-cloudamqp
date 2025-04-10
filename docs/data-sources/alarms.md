---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_alarms"
description: |-
  Get a list of pre-defined or created alarms.
---

# cloudamqp_alarms

Use this data source to retrieve a list of default or created alarms.

## Example Usage

```hcl
data "cloudamqp_alarms" "queue_alarms" {
  instance_id = cloudamqp_instance.instance.id
  type        = "queue"
}
```

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier.
* `type`        - (Optional) The alarm type to filter for. Supported
                  [alarm types](#alarm-types).

## Attributes Reference

All attributes reference are computed

*`alarms`               - List of alarms (see [below for nested schema](#nestedatt--alarms))

<a id="nestedatt--alarms"></a>

### Nested Schema for `alarms`

* `alarm_id`            - The alarm identifier.
* `type`                - The type of the alarm.
* `enabled`             - Enable/disable status of the alarm.
* `value_threshold`     - The value threshold that triggers the alarm.
* `reminder_interval`   - The reminder interval (in seconds) to resend the alarm if not resolved.
                          Set to 0 for no reminders.
* `time_threshold`      - The time interval (in seconds) the `value_threshold` should be active
                          before trigger an alarm.
* `queue_regex`         - Regular expression for which queue to check.
* `vhost_regex`         - Regular expression for which vhost to check
* `recipients`          - Identifier for recipient to be notified.
* `message_type`        - Message type `(total, unacked, ready)` used by queue alarm type.

Specific attribute for `disk` alarm

* `value_calculation`   - Disk value threshold calculation, `(fixed, percentage)` of disk space
                          remaining.

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Alarm Types

`cpu, memory, disk, queue, connection, flow, consumer, netsplit, server_unreachable, notice`

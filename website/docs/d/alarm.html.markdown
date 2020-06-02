---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_alarm"
description: |-
  Get information on pre-defined or created alarms.
---

# cloudamqp_alarm

Use this data source to retrieve information about pre-defined or created alarms. Depends on the identifier of the corresponding `cloudamqp_instance`resource or data source. Then either use `alarm_id` or `type` to retrieve the alarm.

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
* `type`        - (Optional) The alarm type. Either use this or `alarm_id` to give `cloudamqp_alarm` necessary information when retrieve the alarm.

## Attribute reference

* `enabled`         - (Computed) Enable/disable status of the alarm.
* `value_threshold` - (Computed) The value threshold that triggers the alarm.
* `time_threshold`  - (Computed) The time interval (in seconds) the `value_threshold` should be active before trigger an alarm.
* `queue_regex`     - (Computed) Regular expression for which queue to check.
* `vhost_regex`     - (Computed) Regular expression for which vhost to check
* `recipients`      - (Computed) Identifier for recipient to be notified.
* `message_type`    - (Computed) Message type `(total, unacked, ready)` used by queue alarm type.

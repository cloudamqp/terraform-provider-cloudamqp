---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_alarm"
description: |-
  Creates and manages alarms to trigger notifications.
---

# cloudamqp_alarm

This resource allows you to create and manage alarms to trigger and send notifications to assigned recipients. There will always be default alarms (cpu, memory, disk and notice) created upon CloudAMQP instance creation. All default alarms use the default recipient for notifications. This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

Available for all subscription plans, but `lemur`and `tiger`are limited to fewer alarm types. The limited types supported can be seen in the table below in [Alarm Type Reference](#alarm-type-reference).

## Example Usage

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
  value_threshold   = 95
  time_threshold    = 600
  recipient         = [2]
}

# New memory alarm
resource "cloudamqp_alarm" "memory_alarm" {
  instance_id       = cloudamqp_instance.instance.id
  type              = "memiry"
  enabled           = true
  value_threshold   = 95
  time_threshold    = 600
  recipient         = [2]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id`     - (Required) The CloudAMQP instance ID.
* `type`            - (Required) The alarm type, see valid options below.
* `enabled`         - (Required) Enable or disable the alarm to trigger.
* `value_threshold` - (Optional) The value to trigger the alarm for.
* `time_threshold`  - (Optional) The time interval (in seconds) the `value_threshold` should be active before triggering an alarm.
* `queue_regex`     - (Optional) Regex for which queue to check.
* `vhost_regex`     - (Optional) Regex for which vhost to check
* `recipients`      - (Optional) Identifier for recipient to be notified. Leave empty to notify all recipients.
* `message_type`    - (Optional) Message type `(total, unacked, ready)` used by queue alarm type.

Based on alarm type, different arguments are flagged as required or optional.

## Alarm Type reference

Valid options for notification type.

Required arguments for all alarms: *instance_id*, *type* and *enabled*<br>
Optional argument for all alarms: *tags*, *queue_regex*, *vhost_regex*

| Name | Type | Shared | Dedicated | Required arguments |
| ---- | ---- | ---- | ---- | ---- | ---- |
| CPU | cpu | - | &#10004; | time_threshold, value_threshold |
| Memory | memory | - | &#10004;  | time_threshold, value_threshold |
| Disk space | disk | - | &#10004;  | time_threshold, value_threshold |
| Queue | queue | &#10004;  | &#10004;  | time_threshold, value_threshold, queue_regex, vhost_regex, message_type |
| Connection | connection | &#10004; | &#10004; | time_threshold, value_threshold |
| Consumer | consumer | &#10004; | &#10004; | time_threshold, value_threshold, queue, vhost |
| Netsplit | netsplit | - | &#10004; | time_threshold |
| Server unreachable | server_unreachable  | - | &#10004;  | time_threshold |
| Notice | notice | &#10004; | &#10004; |

## Import

`cloudamqp_alarm` can be imported using CloudAMQP internal identifier of a recipient together (CSV separated) with the instance identifier. To see the recipient identifier, use [CloudAMQP API](https://docs.cloudamqp.com/cloudamqp_api.html#list-notification-recipients)

`terraform import cloudamqp_notificaion.recipient <recipient_id>,<instance_od>`

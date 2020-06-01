---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_alarm"
description: |-
  Creates and manages alarms to trigger notifications.
---

# cloudamqp_alarm

This resource allows you to create and manage alarms to trigger and send notifications to given recipients. There will always be default alarms (cpu, memory, disk and notice) created upon CloudAMQP instance creation. All default alarms use the default recipient for notifications. This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

Available for all subscription plans, but `lemur`and `tiger`are limited to fewer alarm types. The limited types supported can be seen in the table below in `Alarm Type Reference'.

## Example Usage

```hcl
resource "cloudamqp_alarm" "default_cpu_alarm" {
  instance_id       = cloudamqp_instance.instance.id
  type              = "cpu"
  enabled           = true
  value_threshol    = 90
  time_threshold    = 600
  recipient         = [1]
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

Required arguments for all alarms: instance_id, type and enabled
Optional argument for all alarms: tags, queue_regex, vhost_regex

Name | Type | Shared | Dedicated | Required arguments |
---- | ---- | ---- | ---- | ---- | ---- |
CPU | cpu | - | x | time_threshold, value_threshold
Memory | memory | - | x | time_threshold, value_threshold
Disk space | disk | - | x | time_threshold, value_threshold
Queue | queue | x | x | time_threshold, value_threshold, queue_regex, vhost_regex, message_type
Connection | connection | x | x | time_threshold, value_threshold
Consumer | consumer | x | x | time_threshold, value_threshold, queue, vhost
Netsplit | netsplit | - | x | time_threshold
Server unreachable | server_unreachable | - | x | time_threshold
Notice | notice | x | x |

## Import

`cloudamqp_alarm` can be imported using CloudAMQP internal identifier of a recipient together (CSV separated) with the instance identifier. To see the recipient identifier, use [CloudAMQP API](https://docs.cloudamqp.com/cloudamqp_api.html#list-notification-recipients)

`terraform import cloudamqp_notificaion.recipient <recipient_id>,<instance_od>`

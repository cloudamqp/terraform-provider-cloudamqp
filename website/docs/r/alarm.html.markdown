---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_alarm"
description: |-
  Creates and manages alarms to trigger notifications.
---

# cloudamqp_instance

This resource allows you to create and manage alarms to trigger and send notifications to give recipients. There will always be a default alarms (cpu, memory, disk and notice) created upon instance creation. All defualt alarms uses default recipient for notifications.

## Example Usage

```hcl
resource "cloudamqp_alarm" "defaul_cpu_alarm" {
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

* `instance_id` - (Required) The CloudAMQP instance ID.
* `type`        - (Required)
* `enabled`       - (Required)
* `value_threshold`        - (Optional)
* `time_threshold`-
* `queue_regex` -
* `vhost_regex`-
* `recipients`-


## Alarm Type reference

Valid options for notification type.

* cpu
* memory
* disk
* queue
* connection
* consumer
* netsplit
* server_unreachable
* notice

## Import

`cloudamqp_alarm` can be imported using CloudAMQP internal ID of a recipient together (CSV seperated) with the instance ID. To see the ID of a recipient, use [CloudAMQP API](https://docs.cloudamqp.com/cloudamqp_api.html#list-notification-recipients)

`terraform import cloudamqp_notificaion.recipient <recpient_ID>,<ID>`

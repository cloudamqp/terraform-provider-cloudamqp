---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_notification"
description: |-
  Creates and manages recipients to receive alarm notifications.
---

# cloudamqp_notification

This resource allows you to create and manage recipients to receive alarm notifications. There will always be a default recipient created upon instance creation. This recipient will use team email and receive notifications from default alarms. This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

Available for all subscription plans.

## Example Usage

```hcl
# New recipient to receieve notifications
resource "cloudamqp_notification" "recipient_01" {
  instance_id = cloudamqp_instance.instance.id
  type        = "email"
  value       = "alarm@example.com"
  name        = "alarm"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance ID.
* `type`        - (Required) Type of the notification. See valid options below.
* `value`       - (Required) Endpoint to send the notification.
* `name`        - (Optional) Display name of the recipient.

## Notification Type reference

Valid options for notification type.

* email
* webhook
* pagerduty
* victorops
* opsgenie
* opsgenie-eu
* slack

## Import

`cloudamqp_notification` can be imported using CloudAMQP internal identifier of a recipient together (CSV separated) with the instance identifier. To see the ID of a recipient, use [CloudAMQP API](https://docs.cloudamqp.com/cloudamqp_api.html#list-notification-recipients)

`terraform import cloudamqp_notificaion.recipient <recpient_id>,<indstance_id>`

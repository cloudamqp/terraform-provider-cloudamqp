---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_notification"
description: |-
  Creates and manages recipients to receive alarm notifications.
---

# cloudamqp_instance

This resource allows you to create and manage recipients to receive alarm notifications. There will always be a default recipient created upon instance creation. This recipient will use team email and receive notifications from default alarms.

## Example Usage

```hcl
resource "cloudamqp_notification" "default_recipient" {
  instance_id       = cloudamqp_instance.instance.id
  type              = "email"
  value             = "notification@example.com"
  name              = "default"
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

`cloudamqp_notification` can be imported using CloudAMQP internal ID of a recipient together (CSV seperated) with the instance ID. To see the ID of a recipient, use [CloudAMQP API](https://docs.cloudamqp.com/cloudamqp_api.html#list-notification-recipients)

`terraform import cloudamqp_notificaion.recipient <recpient_ID>,<ID>`

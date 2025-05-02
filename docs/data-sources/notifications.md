---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_notifications"
description: |-
  Get information on pre-defined or created recipients.
---

# cloudamqp_notifications

Use this data source to retrieve information about all notification recipients. Each recipient will
receive notifications assigned to an alarm that has triggered.

## Example Usage

```hcl
data "cloudamqp_notifications" "default_recipient" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attributes Reference

All attributes reference are computed

* `recipients` - List of alarms (see [below for nested schema](#nestedatt--recipients))

<a id="nestedatt--recipients"></a>

### Nested Schema for `recipients`

All attributes reference are computed

* `recipient_id`  - The identifier for the recipient.
* `name`          - The name of the recipient.
* `type`          - The type of the recipient.
* `value`         - The notification endpoint, where to send the notification.
* `options`       - Options argument (e.g. `rk` used for VictorOps routing key).

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

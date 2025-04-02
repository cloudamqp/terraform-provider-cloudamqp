---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_notification"
description: |-
  Get information on pre-defined or created recipients.
---

# cloudamqp_notification

Use this data source to retrieve information about default or created recipients. The recipient will
receive notifications assigned to an alarm that has triggered. To retrieve the recipient either use
`recipient_id` or `name`.

## Example Usage

```hcl
data "cloudamqp_notification" "default_recipient" {
  instance_id = cloudamqp_instance.instance.id
  name        = "default"
}
```

## Argument Reference

* `instance_id`   - (Required) The CloudAMQP instance identifier.
* `recipient_id`  - (Optional) The recipient identifier.
* `name`          - (Optional) The name set for the recipient.

## Attributes Reference

All attributes reference are computed

* `id`    - The identifier for this resource.
* `type`  - The type of the recipient.
* `value` - The notification endpoint, where to send the notification.
* `options`- Options argument (e.g. `rk` used for VictorOps routing key).

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

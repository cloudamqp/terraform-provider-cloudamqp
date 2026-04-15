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

* `id`          - The identifier for this resource.
* `type`        - The type of the recipient.
* `value`       - The notification endpoint, where to send the notification.
* `options`     - Options argument (e.g. `rk` used for VictorOps routing key).
* `responders`  - An array of reponders (only for OpsGenie). Each `responders` block
                  consists of the field documented below.

___

The options parameter:

* rk        - (Optional) Routing key to route alarm notification (can be used with Victorops).
* dedupkey  - (Optional) If multiple alarms are triggered using a recipient with this key, only the
              the first alarm will trigger a notification (can be used with PagerDuty). Leave blank
              to use the generated dedup key.

___

The `responders` block consists of:

* `type`      - (Required) Type of responder. [`team`, `user`, `escalation`, `schedule`]
* `id`        - (Optional) Identifier in UUID format
* `name`      - (Optional) Name of the responder
* `username`  - (Optional) Username of the responder

Responders of type `team`, `escalation` and `schedule` can use either id or name.
While `user` can use either id or username.

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

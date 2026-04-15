---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_notifications"
description: |-
  Get information on pre-defined or created recipients.
---

<!-- markdownlint-disable MD033 -->

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

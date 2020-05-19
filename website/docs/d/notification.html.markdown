---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_notification"
description: |-
  Get information about the recipients for the CloudAMQP instance.
---

# cloudamqp_notification

Use this data source to retrieve information about the recipients created by CloudAMQP instance. Require to know the identifier of the corresponding `cloudamqp_instance`resource or data source. Then either `recipient_id` or `name` to retrieve the recipient.

## Eample Usage

```hcl
data "cloudamqp_notification" "default_recipient" {
  instance_id = cloudamqp_instance.instance.id
  name = "defualt"
}
```

## Argument reference

* `instance_id`   - (Required) The CloudAMQP instance identifier.
* `recipient_id`  - (Optional) The recipient identifier.
* `name`          - (Optional) The name set for the recipient.

## Attribute reference

* `type`  - (Computed) The type of the recipient.
* `value` - (Computed) The notification endpoint, where to send the notification.

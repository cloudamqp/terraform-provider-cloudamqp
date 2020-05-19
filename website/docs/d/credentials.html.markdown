---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_credentials"
description: |-
  Get credentials information, extracted from `cloudamqp_instance.url`
---

# cloudamqp_alarm

Use this data source to retrieve information about the credentials configured in Rabbit MQ. Require to know the identifier of the corresponding `cloudamqp_instance`resource or data source.

## Eample Usage

```hcl
data "cloudamqp_credentials" "credentials" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attribute reference

* `username`    - (Computed/Sensitive) The username configured in Rabbit MQ.
* `password`    - (Computed/Sensitive) The password used by the `username`.

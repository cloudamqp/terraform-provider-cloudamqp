---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_credentials"
description: |-
  Get credentials information
---

# cloudamqp_credentials

Use this data source to retrieve information about the credentials of the configured user in Rabbit MQ. Information is extracted from `cloudamqp_instance.instance.url`.

## Example Usage

```hcl
data "cloudamqp_credentials" "credentials" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attributes reference

All attributes reference are computed.

* `id`          - The identifier for this data source.
* `username`    - (Sensitive) The username for the configured user in Rabbit MQ.
* `password`    - (Sensitive) The password used by the `username`.

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

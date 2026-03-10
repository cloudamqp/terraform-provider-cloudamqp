---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_credentials"
description: |-
  Get credentials information
---

# cloudamqp_credentials

~> **Deprecated** This data source is deprecated and will be removed in a future version. Use the `credentials` attribute from the `cloudamqp_instance` resource instead.

Use this data source to retrieve information about the credentials of the configured user in
RabbitMQ. Information is extracted from `cloudamqp_instance.instance.url`.

## Example Usage

```hcl
data "cloudamqp_credentials" "credentials" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attributes Reference

All attributes reference are computed.

* `id`          - The identifier for this data source.
* `username`    - (Sensitive) The username for the configured user in Rabbit MQ.
* `password`    - (Sensitive) The password used by the `username`.

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Known issues

The data source causes unnecessary provider reconfigurations when the associated `cloudamqp_instance` resource changes, leading to potential authentication failures during apply operations.

Migration example:

```hcl
# Old (deprecated)
data "cloudamqp_credentials" "credentials" {
  instance_id = cloudamqp_instance.instance.id
}

# New (recommended)
# Access credentials directly from the resource
resource "cloudamqp_instance" "instance" {
  # ...
}

# Use: cloudamqp_instance.instance.credentials.username
#      cloudamqp_instance.instance.credentials.password
```

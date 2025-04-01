---
layout: "cloudamqp"
page_title: "Message Broker HTTP API"
subcategory: "info"
description: |-
  Example of using Message Broker HTTP API to create user, vhost etc.
---

# Message Broker HTTP API

This provider doesn't support using the Message Broker HTTP API to create user, vhost etc. on the
underlying RabbitMQ or LavinMQ instance.

There are other providers here at the registry that are built for this. Example of the current
unofficial RabbitMQ provider with most downloads: [cyrilgdn]

## Example Usage

This can be used together with the CloudAMQP provider to access the Message Broker HTTP API to
create user, vhost etc.

```hcl
data "cloudamqp_credentials" "credentials" {
  instance_id = cloudamqp_instance.instance.id
}

provider "rabbitmq" {
  endpoint = "https://${cloudamqp_instance.instance.host}"
  username = data.cloudamqp_credentials.credentials.username
  password = data.cloudamqp_credentials.credentials.password
}
```

[cyrilgdn]: https://registry.terraform.io/providers/cyrilgdn/rabbitmq/latest

---
layout: "cloudamqp"
page_title: "HTTP API"
subcategory: "info"
description: |-
  Example of using HTTP API to create user, vhost etc.
---

# HTTP API

This provider doesn't support using the HTTP API to create user, vhost etc. on the underlying
RabbitMQ or LavinMQ instance.

There are other provider here at the registry that are built for this. Example of the current
unofficial RabbitMQ provider with most downloads:
https://registry.terraform.io/providers/cyrilgdn/rabbitmq/latest

## Example Usage

This can be used together with the CloudAMQP provider to access the HTTP API to create
user, vhost etc.

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

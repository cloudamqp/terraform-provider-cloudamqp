---
layout: "cloudamqp"
page_title: "Message Broker HTTP API"
subcategory: "info"
description: |-
  Example of using Message Broker HTTP API to create user, vhost etc.
---

# Message Broker HTTP API

This provider doesn't support using the Message Broker HTTP API to create user, vhost etc. on the underlying LavinMQ or RabbitMQ broker.

From CloudAMQP provider [v1.44.0] the `cloudamqp_instance` resource exposes the credentials, that can be used for provider to provider configuration.

There are other providers at the registry that are built for this.

- LavinMQ Broker API provider: [cloudamqp/lavinmq]
- Most downloaded unofficial RabbitMQ Broker API provider: [cyrilgdn]

## Example Usage

### LavinMQ

This can be used together with the CloudAMQP provider to access the Message Broker HTTP API to
create user, vhost etc.

```hcl
terraform {
  required_providers {
    cloudamqp = {
      source  = "cloudamqp/cloudamqp"
      version = "~> 1.44"
    }
    lavinmq = {
      source  = "cloudamqp/lavinmq"
      version = "~> 1.0"
    }
  }
}

resource "cloudamqp_instance" "lavinmq_instance" {
  name   = "provider-to-provider-configuration"
  plan   = "penguin-1"
  region = "amazon-web-services::us-east-1"
  tags   = ["terraform"]
}

provider "lavinmq" {
  baseurl  = format("https://%s", cloudamqp_instance.lavinmq_instance.host)
  username = cloudamqp_instance.lavinmq_instance.credentials.username
  password = cloudamqp_instance.lavinmq_instance.credentials.password
}

resource "lavinmq_vhost" "new_vhost" {
  name = "new_vhost"
}
```

### RabbitMQ

```hcl
terraform {
  required_providers {
    cloudamqp = {
      source  = "cloudamqp/cloudamqp"
      version = "~> 1.44"
    }
    rabbitmq = {
      source  = "cyrilgdn/rabbitmq"
      version = "~> 1.0"
    }
  }
}

resource "cloudamqp_instance" "rabbitmq_instance" {
  name   = "provider-to-provider-configuration"
  plan   = "bunny-1"
  region = "amazon-web-services::us-east-1"
  tags   = ["terraform"]
}

provider "rabbitmq" {
  endpoint = format("https://%s", cloudamqp_instance.rabbitmq_instance.host)
  username = cloudamqp_instance.rabbitmq_instance.credentials.username
  password = cloudamqp_instance.rabbitmq_instance.credentials.password
}

resource "rabbitmq_vhost" "new_vhost" {
  name = "new_vhost"
}
```

## Security Considerations

The credentials exposed by `cloudamqp_instance` are sensitive. Ensure:

- State files are stored securely (e.g., encrypted remote backend)
- CI/CD systems properly mask credentials in logs
- Use Terraform's `sensitive` attribute when referencing credentials

## Known issues

### Single provider configuration

If multiple `cloudamqp_instance` resources are created, each resource need its own broker provider configuration. Alias can be used for this, but will increase the complexity of the configuration.

```hcl
provider "lavinmq" {
  alias    = "test"
  baseurl  = format("https://%s", cloudamqp_instance.lavinmq_instance.host)
  username = cloudamqp_instance.lavinmq_instance.credentials.username
  password = cloudamqp_instance.lavinmq_instance.credentials.password
}

resource "lavinmq_vhost" "new_vhost" {
  provider = lavinmq.test
  name     = "new_vhost"
}
```

### cloudamqp_credentials data source

The `cloudamqp_credentials` data source is no longer recommended. Any changes to the associated `cloudamqp_instance` resource cause the data source to refresh, triggering unnecessary provider reconfiguration and authentication failures during apply.

[v1.44.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.44.0
[cloudamqp/lavinmq]: https://registry.terraform.io/providers/cloudamqp/lavinmq/latest
[cyrilgdn]: https://registry.terraform.io/providers/cyrilgdn/rabbitmq/latest

---
layout: "cloudamqp"
page_title: "Provider: CloudAMQP"
description: |-
  The CloudAMQP provider is used to interact with CloudAMQP organization resources.
---

# CloudAMQP Provider

The CloudAMQP provider is used to interact with CloudAMQP organization resources.

The provider allows you to manage your CloudAMQP instances and features. Create, configure and deploy Rabbit MQ to different cloud platforms. The provider needs to be configured with the proper API key before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the CloudAMQP Provider
provider "cloudamqp" {
  apikey        = var.cloudamqp_customer_api_key
}

# Create a new cloudamqp instance
resource "cloudamqp_instance" "instance" {
  name          = "terraform-cloudamqp-instance"
  plan          = "bunny"
  region        = "amazon-web-services::region=us-west-1"
  nodes         = 1
  tags          = [ "terraform" ]
  rmq_version   = "3.8.3"
}

# New recipient to receieve notifications
resource "cloudamqp_notification" "recipient_01" {
  instance_id = cloudamqp_instance.instance.id
  type        = "email"
  value       = "alarm@example.com"
  name        = "alarm"
}

# New cpu alarm
resource "cloudamqp_alarm" "cpu_alarm" {
  instance_id     = cloudamqp_instance.instance.id
  type            = "cpu"
  value_threshold = 90
  time_threshold  = 600
  enabled         = true
  recipients      = [cloudamqp_notification.recipient_01.id]
}

# Configure firewall
resource "cloudamqp_security_firewall" "firewall" {
  instance_id = cloudamqp_insntance.instance.id
  rules {
    ip = "10.54.72.0/0"
    ports = [4567]
    services = ["AMQPS"]
  }
}

# Cloudwatch logs integration
resource "cloudamqp_integration_log" "cloudwatchlog" {
  instance_id = cloudamqp_instance.instance.id
  name = "cloudwatchlog"
  access_key_id = var.aws_access_key
  secret_access_key = var.aws_secret_key
  region = var.aws_region
}

# Cloudwatch metrics integration
resource "cloudamqp_integration_metric" "cloudwatch" {
  instance_id = cloudamqp_instance.instance.id
  name = "cloudwatch"
  access_key_id = var.aws_access_key
  secret_access_key = var.aws_secret_key
  region = var.aws_region
}
```

## Argument Reference

The following arguments are supported in the `provider` block:

* `apikey` - (Required) This is the CloudAMQP Customer API key needed to make calls to the customer API.
             It can be sourced from login in to your CloudAMQP account and go to API access or go
             directly to [API Keys](https://customer.cloudamqp.com/apikeys).

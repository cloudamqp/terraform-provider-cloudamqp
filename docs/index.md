---
layout: "cloudamqp"
page_title: "Provider: CloudAMQP"
description: |-
  The CloudAMQP provider is used to interact with CloudAMQP organization resources.
---

# CloudAMQP Provider

The CloudAMQP provider is used to interact with CloudAMQP organization resources.

The provider allows you to manage your CloudAMQP instances and features. Create, configure and deploy [**RabbitMQ**](https://www.rabbitmq.com/) or [**LavinMQ**](https://lavinmq.com/) to different cloud platforms. The provider needs to be configured with the proper API key before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the CloudAMQP Provider
provider "cloudamqp" {
  apikey          = var.cloudamqp_customer_api_key
  enable_faster_instance_destroy = true // Optional configuration, can be left out.
}

# Create a new cloudamqp instance
resource "cloudamqp_instance" "instance" {
  name          = "terraform-cloudamqp-instance"
  plan          = "bunny-1"
  region        = "amazon-web-services::us-west-1"
  tags          = [ "terraform" ]
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
  instance_id = cloudamqp_instance.instance.id
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
             The API key can also be read from the environment variable `CLOUDAMQP_APIKEY`.

* `enable_faster_instance_destroy` - (Optional) This will speed up the destroy action for `cloudamqp_instance`
                                      when running `terraform destroy`. It's done by skipping delete behaviour
                                      for resources that don't need to be cleaned up when the servers are deleted.
                                      The argument can also be read from the environment variable
                                      `CLOUDAMQP_ENABLE_FASTER_INSTANCE_DESTROY`, default set to false.

___

***List of resources affected by `enable_faster_instance_destroy`:***

* cloudamqp_plugin
* cloudamqp_plugin_community
* cloudamqp_security_firewall

More information can be found under `Enable faster instance destroy` section on respective resource.

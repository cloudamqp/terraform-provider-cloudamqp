---
layout: "cloudamqp"
page_title: "Provider: CloudAMQP"
description: |-
  The CloudAMQP provider is used to interact with CloudAMQP organization resources.
---

# CloudAMQP Provider

The CloudAMQP provider is used to interact with CloudAMQP organization resources.

The provider allows you to manage your CloudAMQP instances and features. Create, configure and deploy Rabbit MQ to different cloud platforms. The provider need to be configured with the proper API key before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the CloudAMQP Provider
provider "cloudamqp" {
  apikey        = "${var.cloudamqp_apikey}"
}

# Create a new cloudamqp instance
resource "cludamqp_instance" "instance" {
  # ...
}
```

## Argument Reference

The following arguments are supported in the `provider` block:

* `apikey` - (Required) This is the CloudAMQP Customer API key needed to make calls to the customer API.
             It can be sourced from login in to your CloudAMQP account and go to API access or go
             directly to `https://customer.cloudamqp.com/apikeys`.

* `baseurl` - (Optional) This is the target CloudAMQP base API endpoint. Only used when running
              CloudAMQP locally during development. Default value set to
              `https://customer.cloudamqp.com`.

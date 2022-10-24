---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_privatelink_azure"
description: |-
  Enable PrivateLink for a CloudAMQP instance hosted in Azure.
---

# cloudamqp_privatelink_aws

Enable PrivateLink for a CloudAMQP instance hosted in Azure. If no existing VPC available when enable PrivateLink, a new VPC will be created with subnet `10.52.72.0/24`.

More information about [CloudAMQP Privatelink](https://www.cloudamqp.com/docs/cloudamqp-privatelink.html#azure-privatelink).

Only available for dedicated subscription plans.

## Example Usage

CloudAMQP instance without existing VPC

```hcl
resource "cloudamqp_instance" "instance" {
  name   = "Instance 01"
  plan   = "squirrel-1"
  region = "azure-arm::westus"
  tags   = ["test"]
  rmq_version = "3.10.8"
}

resource "cloudamqp_privatelink_azure" "privatelink" {
  instance_id = cloudamqp_instance.instance.id
  approved_subscriptions = [
    "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  ]
}
```

CloudAMQP instance already in an existing VPC.

```hcl
resource "cloudamqp_vpc" "vpc" {
  name = "Standalone VPC"
  region = "azure-arm::westus"
  subnet = "10.56.72.0/24"
  tags = ["test"]
}

resource "cloudamqp_instance" "instance" {
  name   = "Instance 01"
  plan   = "squirrel-1"
  region = "azure-arm::westus"
  tags   = ["test"]
  rmq_version = "3.10.8"
  vpc_id = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_privatelink_azure" "privatelink" {
  instance_id = cloudamqp_instance.instance.id
  approved_subscriptions = [
    "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  ]
}
```

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier.
* `approved_subscriptions` - (Required) Approved subscriptions to access the endpoint service. See format below.
* `sleep` - (Optional) Configurable sleep time (seconds) when enable PrivateLink. Default set to 60 seconds.
* `timeout` - (Optional) - Configurable timeout time (seconds) when enable PrivateLink. Default set to 3600 seconds.

Approved subscriptions format: <br>
`XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX`

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.
* `status`- PrivateLink status [enable, pending, disable]
* `service_name` - Service name (alias) of the PrivateLink, needed when creating the endpoint.
* `server_name` - Name of the server having the PrivateLink enabled.

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_privatelink_aws` can be imported using CloudAMQP internal identifier.

`terraform import cloudamqp_privatelink_aws.privatelink <id>`

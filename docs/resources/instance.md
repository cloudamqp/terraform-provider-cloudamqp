---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_instance"
description: |-
  Creates and manages a Rabbit MQ instance within CloudAMQP.
---

# cloudamqp_instance

This resource allows you to create and manage a CloudAMQP instance running Rabbit MQ and deploy to multiple cloud platforms provider and over multiple regions, see [Instance regions](../instance_region.html) for more information.

Once the instance is created it will be assigned a unique identifier. All other resource and data sources created for this instance needs to reference the instance identifier.

## Example Usage

```hcl
# Minimum free lemur instance
resource "cloudamqp_instance" "lemur_instance" {
  name = "terraform-free-instance"
  plan = "lemur"
  region = "amazon-web-services::us-west-1"
}

# New dedicated bunny instance
resource "cloudamqp_instance" "instance" {
  name              = "terraform-cloudamqp-instance"
  plan              = "bunny-1"
  region            = "amazon-web-services::us-west-1"
  nodes             = 1
  tags              = [ "terraform" ]
  rmq_version       = "3.8.3"
  no_default_alarms = true
}
```

## Argument Reference

The following arguments are supported:

* `name`        - (Required) Name of the CloudAMQP instance.
* `plan`        - (Required) The subscription plan. See available [plans](../guides/info_plan.md)
* `region`      - (Required) The region to host the instance in. See [Instance regions](../guides/info_region.md)
* `nodes`       - (Computed) Number of nodes, 1, 3 or 5. **Note: Changed from optional to computed. In order to change number of nodes, the subscription plan needs to be updated.**
* `tags`        - (Optional) One or more tags for the CloudAMQP instance, makes it possible to categories multiple instances in console view. Default there is no tags assigned.
* `rmq_version` - (Computed/Optional) The Rabbit MQ version. Can be left out, will then be set to default value used by CloudAMQP API. **Note: There is not yet any support in the provider to change the RMQ version. Once it's set in the initial creation, it will remain.**
* `vpc_subnet`  - (Optional) Creates a dedicated VPC subnet, shouldn't overlap with other VPC subnet, default subnet used 10.56.72.0/24. **NOTE: extra fee will be charged when using VPC, see [CloudAMQP](https://cloudamqp.com) for more information.**
* `no_default_alarms`- (Computed/Optional) Set to true to discard creating default alarms when the instance is created. Can be left out, will then use default value = false.



## Attributes Reference

* `url`     - (Computed) AMQP server endpoint. `amqps://{username}:{password}@{hostname}/{vhost}`
* `apikey`  - (Computed) API key needed to communicate to CloudAMQP's second API. The second API is used to manage alarms, integration and more, full description [CloudAMQP API](https://docs.cloudamqp.com/cloudamqp_api.html).
* `host`    - (Computed) The host name for the CloudAMQP instance.
* `vhost`   - (Computed) The virtual host used by Rabbit MQ.

## Import
`cloudamqp_instance`can be imported using CloudAMQP internal identifier. To retrieve the identifier for an instance, use [CloudAMQP customer API](https://docs.cloudamqp.com/#list-instances).

`terraform import cloudamqp_instance.instance <instance_id>`

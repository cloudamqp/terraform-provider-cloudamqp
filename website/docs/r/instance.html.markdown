---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_instance"
description: |-
  Creates and manages a Rabbit MQ instance within CloudAMQP.
---

# cloudamqp_instance

This resource allows you to create and manage Rabbit MQ instance through CloudAMQP and deploy to different cloud platforms. Minimum required arguments are name, plan and region. Once the instance is created it will be assigned a unique ID. This will be needed in all other data sources and resources.

## Example Usage

```hcl
resource "cloudamqp_instance" "instance" {
  name          = "terraform-cloudamqp-instance"
  plan          = "bunny"
  region        = "amazon-web-services::region = us-west-1"
  nodes         = 1
  tags          = [ "terraform" ]
  rmq_version   = "3.8.3"
  vpc_subnet    = "10.56.72.0/24"
}
```

## Argument Reference

The following arguments are supported:

* `name`        - (Required) Name of the CloudAMQP instance
* `plan`        - (Required) Subscription plan, see available [plans](`https://www.cloudamqp.com/plans.html`).
* `region`      - (Required) Name of the region to create the instance in. Combines cloud provider and region, {cloud_provider}::{region}
* `nodes`       - (Optional) Number of nodes, 1 to 3, in the CloudAMQP instance. The plan choosen must support number of nodes. Default set to 1.
* `tags`        - (Optional) One or mote tags for the CloudAMQP instance, makes it possible to categories multiple instances in console view. Default no tags assigned.
* `rmq_version` - (Optional) The Rabbit MQ version. Default set to current loaded default value in CloudAMQP API.
* `vpc_subnet`  - (Optional) Creates dedicated VPC subnet, shouldn't overlap with other VPC subnet. Extra fee ($99/month) will be charged when enable VPC. Default subnet used 10.56.72.0/24.

## Attributes Reference

* `url`     - (Computed) AMQP server endpoint. amqps://{username}:{password}@{hostname}/{vhost}
* `apikey`  - (Computed) API key needed to communicate to CloudAMQP secondary API. Used per instance to manage alarms, integration and more. `https://docs.cloudamqp.com/cloudamqp_api.html`.
* `host`    - (Computed) The host name for the CloudAMQP instance.
* `vhost`   - (Computed) The virtual host used by Rabbit MQ.

## Import
`cloudamqp_instance`can be imported using CloudAMQP internal ID of an instance. To see the ID of an instance, use [CloudAMQP customer API](https://docs.cloudamqp.com/#instances).

`terraform import cloudamqp_instance.instance <ID>`

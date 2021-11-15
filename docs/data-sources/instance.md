---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_instance"
description: |-
  Get information about an already created CloudAMQP instance
---

# cloudamqp_instance

Use this data source to retrieve information about an already created CloudAMQP instance. In order to retrieve the correct information, the CoudAMQP instance identifier is needed.

## Example Usage

```hcl
data "cloudamqp_instance" "instance" {
  instance_id = <id>
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attributes reference

All attributes reference are computed

* `id`          - The identifier for this resource.
* `name`        - The name of the CloudAMQP instance.
* `plan`        - The subscription plan for the CloudAMQP instance.
* `region`      - The cloud platform and region that host the CloudAMQP instance, `{platform}::{region}`.
* `vpc_id`      - ID of the VPC configured for the CloudAMQP instance.
* `vpc_subnet`  - Dedicated VPC subnet configured for the CloudAMQP instance.
* `nodes`       - Number of nodes in the cluster of the CloudAMQP instance.
* `rmq_version` - The version of installed Rabbit MQ.
* `url`         - (Sensitive) The AMQP URL (uses the internal hostname if the instance was created with VPC), used by clients to connect for pub/sub.
* `apikey`      - (Sensitive) The API key to secondary API handing alarms, integration etc.
* `tags`        - Tags the CloudAMQP instance with categories.
* `host`        - The external hostname for the CloudAMQP instance.
* `host_internal` - The internal hostname for the CloudAMQP instance.
* `vhost`       - The virtual host configured in Rabbit MQ.

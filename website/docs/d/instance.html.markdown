---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_instance"
description: |-
  Get information about the already created CloudAMQP instance
---

# cloudamqp_instance

Use this data source to retrieve information about the CloudAMQP instance. Require to know the CloudAMQP instance identifier in order to retrieve the correct information.

## Eample Usage

```hcl
data "cloudamqp_instance" "instance" {
  instance_id = <id>
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attribute reference

* `name`        - (Computed) The name of the CloudAMQP instance.
* `plan`        - (Computed) The subscription plan for the CloudAMQP instance.
* `region`      - (Computed) The cloud platform and region that host the CloudAMQP instance, `{platform}::{region}`.
* `vpc_subnet`  - (Computed) Dedicated VPC subnet configured for the CloudAMQP instance.
* `nodes`       - (Computed) Number of nodes in the cluster of the CloudAMQP instance.
* `rmq_version` - (Computed) The version of installed Rabbit MQ.
* `url`         - (Computed/Sensitive) The AMQP url, used by clients to connect for pub/sub. Contains configured Rabbit MQ credentials.
* `apikey`      - (Computed/Sensitive) The API key to secondary API handing alarms, integration etc.
* `tags`        - (Computed) Tags the CloudAMQP instance with categories.
* `host`        - (Computed) The hostname configred for the CloudAMQP instance.
* `vhost`       - (Computed) The virtaul host configured in Rabbit MQ.

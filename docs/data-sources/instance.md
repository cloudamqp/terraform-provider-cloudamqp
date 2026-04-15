---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_instance"
description: |-
  Get information about an already created CloudAMQP instance
---

<!-- markdownlint-disable MD033 -->

# cloudamqp_instance

Use this data source to retrieve information about an already created CloudAMQP instance. In order
to retrieve the correct information, the CoudAMQP instance identifier is needed.

## Example Usage

```hcl
data "cloudamqp_instance" "instance" {
  instance_id = <id>
}
```

<details>
  <summary>
    <b>
      <i>Provider-to-provider configuration, from </i>
      <a href="https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.44.1">v1.44.1</a>
    </b>
  </summary>

```hcl
data "cloudamqp_instance" "instance" {
  instance_id = <id>
}

provider "lavinmq" {
  baseurl  = format("https://%s", data.cloudamqp_instance.instance.host)
  username = data.cloudamqp_instance.instance.credentials.username
  password = data.cloudamqp_instance.instance.credentials.password
}

resource "lavinmq_vhost" "new_vhost" {
  name = "new_vhost"
}
```

</details>

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attributes Reference

All attributes reference are computed

* `id`          - The identifier for this resource.
* `name`        - The name of the CloudAMQP instance.
* `plan`        - The subscription plan for the CloudAMQP instance.
* `region`      - The cloud platform and region that host the CloudAMQP instance,
                  `{platform}::{region}`.
* `vpc_id`      - ID of the VPC configured for the CloudAMQP instance.
* `vpc_subnet`  - Dedicated VPC subnet configured for the CloudAMQP instance.
* `nodes`       - Number of nodes in the cluster of the CloudAMQP instance.
* `rmq_version` - The version of installed Rabbit MQ.
* `url`         - (Sensitive) The AMQP URL (uses the internal hostname if the instance was created
                  with VPC), used by clients to connect for pub/sub.
* `apikey`      - (Sensitive) The API key to secondary API handing alarms, integration etc.
* `tags`        - Tags the CloudAMQP instance with categories.
* `host`        - The external hostname for the CloudAMQP instance.
* `host_internal` - The internal hostname for the CloudAMQP instance.
* `vhost`       - The virtual host configured in Rabbit MQ.
* `dedicated`   - Information if the CloudAMQP instance is shared or dedicated.
* `backend`     - Information if the CloudAMQP instance runs either RabbitMQ or LavinMQ.
* `credentials`   - (Sensitive) Broker credentials block with information extracted from URL.

___

The `credentials` block consists of:

* `username` - The username to access the broker.
* `password` - The password for the user to access the broker.

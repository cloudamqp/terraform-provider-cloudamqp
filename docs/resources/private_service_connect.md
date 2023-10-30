---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_private_service_connect"
description: |-
  Enable Private Service Connect (Privatelink) for a CloudAMQP instance hosted in GCP.
---

# cloudamqp_private_service_connect

Enable Private Service Connect (Privatelink) for a CloudAMQP instance hosted in GCP. If no existing
VPC available when enable Private Service Connect, a new VPC will be created with subnet `10.52.72.0/24`.

~> **Note:** Enabling Private Service Connect will automatically add firewall rules for the peered subnet.
<details>
 <summary>
    <i>Default Private Service Connect firewall rule</i>
  </summary>
```hcl
rules {
  Description = "Private Service Connect"
  ip          = "10.0.0.0/24"
  ports       = []
  services    = ["AMQP", "AMQPS", "HTTPS", "STREAM", "STREAM_SSL", "STOMP", "STOMPS", "MQTT", "MQTTS"]
}
```
</details>


Only available for dedicated subscription plans.

## Example Usage

<details>
  <summary>
    <b>
      <i>CloudAMQP instance without existing VPC</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_instance" "instance" {
  name   = "Instance 01"
  plan   = "bunny-1"
  region = "amazon-web-services::us-west-1"
  tags   = []
}

resource "cloudamqp_private_service_connect" "private_service_connect" {
  instance_id = cloudamqp_instance.instance.id
  allowed_project = [
    "<gcp-project-id>"
  ]
}
```
</details>

<details>
  <summary>
    <b>
      <i>CloudAMQP instance in an existing VPC</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_vpc" "vpc" {
  name = "Standalone VPC"
  region = "amazon-web-services::us-west-1"
  subnet = "10.56.72.0/24"
  tags = []
}

resource "cloudamqp_instance" "instance" {
  name   = "Instance 01"
  plan   = "bunny-1"
  region = "amazon-web-services::us-west-1"
  tags   = []
  vpc_id = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_private_service_connect" "private_service_connect" {
  instance_id = cloudamqp_instance.instance.id
  allowed_projects = [
    "<gcp-project-id>"
  ]
}
```
</details>

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier.
* `allowed_projects` - (Required) Only allowed GCP project will have access.
* `sleep` - (Optional) Configurable sleep time (seconds) when enable Private Service Connect. Default set to 60 seconds.
* `timeout` - (Optional) - Configurable timeout time (seconds) when enable Private Service Connect. Default set to 3600 seconds.

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.
* `status`- Private Service Connect status [enable, pending, disable]

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_private_service_connect` can be imported using CloudAMQP internal identifier.

`terraform import cloudamqp_private_service_connect.private_serivce_connect <id>`

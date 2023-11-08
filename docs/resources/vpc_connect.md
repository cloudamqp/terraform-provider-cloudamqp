---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_vpc_connect"
description: |-
  Enable VPC connect (Privatelink or Private Service Connect) for a CloudAMQP instance hosted in
  AWS, Azure or GCP.
---

# cloudamqp_vpc_connect

This resource is a generic way to handle PrivateLink (AWS and Azure) and Private Service Connect (GCP).
Communication between resources can be done just as they were living inside a VPC. CloudAMQP creates an Endpoint
Service to connect the VPC and creating a new network interface to handle the communicate.

If no existing VPC available when enable VPC connect, a new VPC will be created with subnet `10.52.72.0/24`.

More information can be found at: [CloudAMQP privatelink](https://www.cloudamqp.com/docs/cloudamqp-privatelink.html)

-> **Note:** Enabling VPC Connect will automatically add a firewall rule.

<details>
 <summary>
    <b>
      <i>Default PrivateLink firewall rule [AWS, Azure]</i>
    </b>
  </summary>
```hcl
rules {
  Description = "PrivateLink setup"
  ip          = "<VPC Subnet>"
  ports       = []
  services    = ["AMQP", "AMQPS", "HTTPS", "STREAM", "STREAM_SSL", "STOMP", "STOMPS", "MQTT", "MQTTS"]
}
```
</details>

<details>
 <summary>
    <b>
      <i>Default Private Service Connect firewall rule [GCP]</i>
    </b>
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

If now managadable standalone VPC exists for the CloudAMQP instance. One will be avai

<details>
  <summary>
    <b>
      <i>Enable VPC Connect (PrivateLink) in AWS</i>
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

resource "cloudamqp_vpc_connect" "vpc_connect" {
  instance_id = cloudamqp_instance.instance.id
  region = cloudamqp_instance.instance.region
  allowlist = [
    "arn:aws:iam::aws-account-id:user/user-name"
  ]
}
```
</details>

<details>
  <summary>
    <b>
      <i>Enable VPC Connect (PrivateLink) in Azure</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_vpc" "vpc" {
  name = "Standalone VPC"
  region = "azure-arm::westus"
  subnet = "10.56.72.0/24"
  tags = []
}

resource "cloudamqp_instance" "instance" {
  name   = "Instance 01"
  plan   = "bunny-1"
  region = "azure-arm::westus"
  tags   = []
  vpc_id = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_vpc_connect" "vpc_connect" {
  instance_id = cloudamqp_instance.instance.id
  region = cloudamqp_instance.instance.region
  allowlist = [
    "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  ]
}
```
</details>

<details>
  <summary>
    <b>
      <i>Enable VPC Connect (Private Service Connect) in GCP</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_vpc" "vpc" {
  name = "Standalone VPC"
  region = "google-compute-engine::us-west1"
  subnet = "10.56.72.0/24"
  tags = []
}

resource "cloudamqp_instance" "instance" {
  name   = "Instance 01"
  plan   = "bunny-1"
  region = "google-compute-engine::us-west1"
  tags   = []
  vpc_id = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_vpc_connect" "vpc_connect" {
  instance_id = cloudamqp_instance.instance.id
  region = cloudamqp_instance.instance.region
  allowlist = [
    "some-project-123456"
  ]
}
```
</details>

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier.
* `region` - (Required) The region where the CloudAMQP instance is hosted.
* `allowlist` - (Required) List of allowed principals, projects or subscriptions depending
  on hosting platform provider. See below table.
* `sleep` - (Optional) Configurable sleep time (seconds) when enable Private Service Connect.
  Default set to 60 seconds.
* `timeout` - (Optional) Configurable timeout time (seconds) when enable Private Service Connect.
  Default set to 3600 seconds.

___

The `allowlist` data depending on platform:

| Platform | Description         | Format                                                                                                                             |
|----------|---------------------|------------------------------------------------------------------------------------------------------------------------------------|
| AWS      | IAM ARN principals  | arn:aws:iam::aws-account-id:root<br /> arn:aws:iam::aws-account-id:user/user-name<br /> arn:aws:iam::aws-account-id:role/role-name |
| Azure    | Subscription (GUID) | XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX                                                                                               |
| GCP      | Project IDs*        | 6 to 30 lowercase letters, digits, or hyphens                                                                                      |

*https://cloud.google.com/resource-manager/reference/rest/v1/projects

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.
* `status`- Private Service Connect status [enable, pending, disable]
* `service_name` - Service name (alias for Azure) of the PrivateLink.
* `active_zones` - Covering availability zones used when creating an endpoint from other VPC. (AWS)

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

Since `region` also is required, suggest to reuse the argument from CloudAMQP instance,
`cloudamqp_instance.instance.region`.

## Import

`cloudamqp_vpc_connect` can be imported using CloudAMQP internal identifier.

`terraform import cloudamqp_vpc_connect.vpc_connect <id>`

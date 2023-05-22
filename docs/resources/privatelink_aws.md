---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_privatelink_aws"
description: |-
  Enable PrivateLink for a CloudAMQP instance hosted in AWS.
---

# cloudamqp_privatelink_aws

Enable PrivateLink for a CloudAMQP instance hosted in AWS. If no existing VPC available when enable PrivateLink, a new VPC will be created with subnet `10.52.72.0/24`.

~> **NOTE:** Once the PrivateLink is enabled our backend will automatically create a firewall rule for this.
<details>
 <summary>
    <i>Default PrivateLink firewall rule</i>
  </summary>
```hcl
rules {
  Description = "PrivateLink setup"
  ip          = "10.56.72.0/24"
  ports       = []
  services    = ["AMQP", "AMQPS", "HTTPS", "STREAM", "STREAM_SSL", "STOMP", "STOMPS", "MQTT", "MQTTS"]
}
```
</details>

Pricing is available at [cloudamqp.com](https://www.cloudamqp.com/plans.html) and more information about [CloudAMQP Privatelink](https://www.cloudamqp.com/docs/cloudamqp-privatelink.html#aws-privatelink).

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

resource "cloudamqp_privatelink_aws" "privatelink" {
  instance_id = cloudamqp_instance.instance.id
  allowed_principals = [
    "arn:aws:iam::aws-account-id:user/user-name"
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

resource "cloudamqp_privatelink_aws" "privatelink" {
  instance_id = cloudamqp_instance.instance.id
  allowed_principals = [
    "arn:aws:iam::aws-account-id:user/user-name"
  ]
}
```
</details>

<details>
  <summary>
    <b>
      <i>CloudAMQP instance in an existing VPC with managed firewall rules</i>
    </b>
  </summary>

Chain firewall to privatelink resources to make sure the firewall rules are applied after privatelink have been enabled.

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

resource "cloudamqp_privatelink_aws" "privatelink" {
  instance_id = cloudamqp_instance.instance.id
  allowed_principals = [
    "arn:aws:iam::aws-account-id:user/user-name"
  ]
}

resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  rules {
    Description = "PrivateLink setup"
    ip          = cloudamqp_vpc.vpc.subnet
    ports       = []
    services    = ["AMQP", "AMQPS", "HTTPS", "STREAM", "STREAM_SSL"]
  }

  rules {
    description = "MGMT interface"
    ip = "0.0.0.0/0"
    ports = []
    services = ["HTTPS"]
  }

  depends_on = [
    cloudamqp_privatelink_aws.privatelink
   ]
}
```
</details>

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier.
* `allowed_principals` - (Required) Allowed principals to access the endpoint service.
* `sleep` - (Optional) Configurable sleep time (seconds) when enable PrivateLink. Default set to 60 seconds.
* `timeout` - (Optional) - Configurable timeout time (seconds) when enable PrivateLink. Default set to 3600 seconds.

Allowed principals format: <br>
`arn:aws:iam::aws-account-id:root` <br>
`arn:aws:iam::aws-account-id:user/user-name` <br>
`arn:aws:iam::aws-account-id:role/role-name`

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.
* `status`- PrivateLink status [enable, pending, disable]
* `service_name` - Service name of the PrivateLink used when creating the endpoint from other VPC.
* `active_zones` - Covering availability zones used when creating an Endpoint from other VPC.

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_privatelink_aws` can be imported using CloudAMQP internal identifier.

`terraform import cloudamqp_privatelink_aws.privatelink <id>`

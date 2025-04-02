---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_privatelink_aws"
description: |-
  Enable PrivateLink for a CloudAMQP instance hosted in AWS.
---

# cloudamqp_privatelink_aws

Enable PrivateLink for a CloudAMQP instance hosted in AWS. If no existing VPC available when enable
PrivateLink, a new VPC will be created with subnet `10.52.72.0/24`.

-> **Note:** Enabling PrivateLink will automatically add firewall rules for the peered subnet.

<details>
 <summary>
    <i>Default PrivateLink firewall rule</i>
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

Pricing is available at [CloudAMQP plans] where you can also find more information about
[CloudAMQP PrivateLink].

Only available for dedicated subscription plans.

~> **Warning:** This resource considered deprecated and will be removed in next major version (v2.0).
Recommended to start using the new resource [cloudamqp_vpc_connect].

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
  name    = "Standalone VPC"
  region  = "amazon-web-services::us-west-1"
  subnet  = "10.56.72.0/24"
  tags    = []
}

resource "cloudamqp_instance" "instance" {
  name                = "Instance 01"
  plan                = "bunny-1"
  region              = "amazon-web-services::us-west-1"
  tags                = []
  vpc_id              = cloudamqp_vpc.vpc.id
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

## Argument Reference

* `instance_id`         - (Required) The CloudAMQP instance identifier.
* `allowed_principals`  - (Required) Allowed principals to access the endpoint service.
* `sleep`               - (Optional) Configurable sleep time (seconds) when enable PrivateLink.
                          Default set to 10 seconds.

  ***Note:*** Available from [v1.29.0]

* `timeout`             - (Optional) Configurable timeout time (seconds) when enable PrivateLink.
                          Default set to 1800 seconds.

  ***Note:*** Available from [v1.29.0]

Allowed principals format: <br>
`arn:aws:iam::aws-account-id:root` <br>
`arn:aws:iam::aws-account-id:user/user-name` <br>
`arn:aws:iam::aws-account-id:role/role-name`

## Attributes Reference

All attributes reference are computed

* `id`            - The identifier for this resource. Will be same as `instance_id`
* `status`        - PrivateLink status [enable, pending, disable]
* `service_name`  - Service name of the PrivateLink used when creating the endpoint from other VPC.
* `active_zones`  - Covering availability zones used when creating an Endpoint from other VPC.

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_privatelink_aws` can be imported using CloudAMQP instance identifier. To retrieve the
identifier, use [CloudAMQP API list intances].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_privatelink_aws.privatelink
  id = cloudamqp_instance.instance.id
}
```

Or use Terraform CLI:

`terraform import cloudamqp_privatelink_aws.privatelink <id>`

## Create PrivateLink with additional firewall rules

To create a PrivateLink configuration with additional firewall rules, it's required to chain the
[cloudamqp_security_firewall] resource to avoid parallel conflicting resource calls. You can do this
by making the firewall resource depend on the PrivateLink resource,
`cloudamqp_privatelink_aws.privatelink`.

Furthermore, since all firewall rules are overwritten, the otherwise automatically added rules for
the PrivateLink also needs to be added.

## Example usage with additional firewall rules

<details>
  <summary>
    <b>
      <i>CloudAMQP instance in an existing VPC with managed firewall rules</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_vpc" "vpc" {
  name    = "Standalone VPC"
  region  = "amazon-web-services::us-west-1"
  subnet  = "10.56.72.0/24"
  tags    = []
}

resource "cloudamqp_instance" "instance" {
  name                = "Instance 01"
  plan                = "bunny-1"
  region              = "amazon-web-services::us-west-1"
  tags                = []
  vpc_id              = cloudamqp_vpc.vpc.id
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
    description = "Custom PrivateLink setup"
    ip          = cloudamqp_vpc.vpc.subnet
    ports       = []
    services    = ["AMQP", "AMQPS", "HTTPS", "STREAM", "STREAM_SSL"]
  }

  rules {
    description = "MGMT interface"
    ip          = "0.0.0.0/0"
    ports       = []
    services    = ["HTTPS"]
  }

  depends_on = [
    cloudamqp_privatelink_aws.privatelink
   ]
}
```

</details>

[CloudAMQP API list intances]: https://docs.cloudamqp.com/#list-instances
[CloudAMQP plans]: https://www.cloudamqp.com/plans.html
[CloudAMQP PrivateLink]: https://www.cloudamqp.com/docs/cloudamqp-privatelink.html#aws-privatelink
[cloudamqp_security_firewall]: ./security_firewall.md
[cloudamqp_vpc_connect]: ./vpc_connect.md
[v1.29.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.29.0

---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_vpc_connect"
description: |-
  Enable VPC connect (Privatelink or Private Service Connect) for a CloudAMQP instance hosted in
  AWS, Azure or GCP.
---

# cloudamqp_vpc_connect

This resource is a generic way to handle PrivateLink (AWS and Azure) and Private Service Connect
(GCP). Communication between resources can be done just as they were living inside a VPC. CloudAMQP
creates an Endpoint Service to connect the VPC and creating a new network interface to handle the
communicate.

If no existing VPC available when enable VPC connect, a new VPC will be created with subnet
`10.52.72.0/24`.

More information can be found at: [CloudAMQP VPC Connect]

-> **Note:** Enabling VPC Connect will automatically add a firewall rule.

<details>
 <summary>
    <b>
      <i>Default PrivateLink firewall rule [AWS, Azure]</i>
    </b>
  </summary>

For LavinMQ:

```hcl
rules {
  Description = "PrivateLink setup"
  ip          = "<VPC Subnet>"
  ports       = [5552, 5551, 61613, 61614]
  services    = ["AMQP", "AMQPS", "HTTPS", "MQTT", "MQTTS"]
}
```

For RabbitMQ:

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

For LavinMQ:

```hcl
rules {
  Description = "Private Service Connect"
  ip          = "10.0.0.0/24"
  ports       = [5552, 5551, 61613, 61614]
  services    = ["AMQP", "AMQPS", "HTTPS", "MQTT", "MQTTS"]
}
```

For RabbitMQ:

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
      <i>Enable VPC Connect (PrivateLink) in AWS</i>
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
  plan                = "penguin-1"
  region              = "amazon-web-services::us-west-1"
  tags                = []
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_vpc_connect" "vpc_connect" {
  instance_id = cloudamqp_instance.instance.id
  region      = cloudamqp_instance.instance.region
  allowed_principals = [
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
  name    = "Standalone VPC"
  region  = "azure-arm::westus"
  subnet  = "10.56.72.0/24"
  tags    = []
}

resource "cloudamqp_instance" "instance" {
  name                = "Instance 01"
  plan                = "penguin-1"
  region              = "azure-arm::westus"
  tags                = []
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_vpc_connect" "vpc_connect" {
  instance_id = cloudamqp_instance.instance.id
  region      = cloudamqp_instance.instance.region
  approved_subscriptions = [
    "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  ]
}
```

The attribute `service_name` found in resource `cloudamqp_vpc_connect` corresponds to the alias in
the resource `azurerm_private_endpoint` of the Azure provider. This can be used when creating the
private endpoint.

```hcl
resource "azurerm_private_endpoint" "example" {
  name                = "example-endpoint"
  location            = data.azurerm_resource_group.example.location
  resource_group_name = data.azurerm_resource_group.example.name
  subnet_id           = data.azurerm_subnet.subnet.id

  private_service_connection {
    name                              = "example-privateserviceconnection"
    private_connection_resource_alias = cloudamqp_vpc_connect.vpc_connect.service_name
    is_manual_connection              = true
    request_message                   = "PL"
  }
}
```

More information about the resource and argument can be found here:
[private_connection_resource_alias]. Or check their example "Using a Private Link Service Alias with
existing resources".

</details>

<details>
  <summary>
    <b>
      <i>Enable VPC Connect (Private Service Connect) in GCP</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_vpc" "vpc" {
  name    = "Standalone VPC"
  region  = "google-compute-engine::us-west1"
  subnet  = "10.56.72.0/24"
  tags    = []
}

resource "cloudamqp_instance" "instance" {
  name                = "Instance 01"
  plan                = "penguin-1"
  region              = "google-compute-engine::us-west1"
  tags                = []
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_vpc_connect" "vpc_connect" {
  instance_id = cloudamqp_instance.instance.id
  region      = cloudamqp_instance.instance.region
  allowed_projects = [
    "some-project-123456"
  ]
}
```

</details>

## Argument Reference

* `instance_id`             - (Required) The CloudAMQP instance identifier.
* `region`                  - (Required) The region where the CloudAMQP instance is hosted.
* `allowed_principals`      - (Optional) List of allowed prinicpals used by AWS, see below table.
* `approved_subscriptions`  - (Optional) List of approved subscriptions used by Azure, see below
                              table.
* `allowed_projects`        - (Optional) List of allowed projects used by GCP, see below table.
* `sleep`                   - (Optional) Configurable sleep time (seconds) when enable Private
                              Service Connect. Default set to 10 seconds.
* `timeout`                 - (Optional) Configurable timeout time (seconds) when enable Private
                              Service Connect. Default set to 1800 seconds.

___

The `allowed_principals`, `approved_subscriptions` or `allowed_projects` data depends on the
provider platform:

| Platform | Description | Format |
|---|---|---|
| AWS | IAM ARN principals | arn:aws:iam::aws-account-id:root<br>arn:aws:iam::aws-account-id:user/user-name<br> arn:aws:iam::aws-account-id:role/role-name |
| Azure | Subscription (GUID) | XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX |
| GCP | Project IDs [Google docs] | 6 to 30 lowercase letters, digits, or hyphens |

## Attributes Reference

All attributes reference are computed

* `id`            - The identifier for this resource. Will be same as `instance_id`
* `status`        - Private Service Connect status [enable, pending, disable]
* `service_name`  - Service name (alias for Azure, see example above) of the PrivateLink.
* `active_zones`  - Covering availability zones used when creating an endpoint from other VPC. (AWS)

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

Since `region` also is required, suggest to reuse the argument from CloudAMQP instance,
`cloudamqp_instance.instance.region`.

## Import

`cloudamqp_vpc_connect` can be imported using CloudAMQP instance identifier. To
retrieve the identifier, use [CloudAMQP API list intances].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_vpc_connect.this
  id = cloudamqp_instance.instance.id
}
```

Or use Terraform CLI:

`terraform import cloudamqp_vpc_connect.vpc_connect <id>`

## Create VPC Connect with additional firewall rules

To create a PrivateLink/Private Service Connect configuration with additional firewall rules, it's
required to chain the [cloudamqp_security_firewall] resource to avoid parallel conflicting resource
calls. You can do this by making the firewall resource depend on the VPC Connect resource
`cloudamqp_vpc_connect.vpc_connect`.

Furthermore, since all firewall rules are overwritten, the otherwise automatically added rules for
the VPC Connect also needs to be added.

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
  plan                = "penguin-1"
  region              = "amazon-web-services::us-west-1"
  tags                = []
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_vpc_connect" "vpc_connect" {
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
    services    = ["AMQP", "AMQPS", "HTTPS"]
  }

  rules {
    description = "MGMT interface"
    ip          = "0.0.0.0/0"
    ports       = []
    services    = ["HTTPS"]
  }

  depends_on = [
    cloudamqp_vpc_connect.vpc_connect
   ]
}
```

</details>

[CloudAMQP API list intances]: https://docs.cloudamqp.com/index.html#tag/instances/get/instances
[CloudAMQP VPC Connect]: https://www.cloudamqp.com/docs/cloudamqp-vpc-connect.html
[cloudamqp_security_firewall]: https://registry.terraform.io/providers/cloudamqp/cloudamqp/latest/docs/resources/security_firewall
[Google docs]: https://cloud.google.com/resource-manager/reference/rest/v1/projects
[private_connection_resource_alias]: ./private_endpoint#private_connection_resource_alias-1

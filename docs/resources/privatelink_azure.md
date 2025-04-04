---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_privatelink_azure"
description: |-
  Enable PrivateLink for a CloudAMQP instance hosted in Azure.
---

# cloudamqp_privatelink_azure

Enable PrivateLink for a CloudAMQP instance hosted in Azure. If no existing VPC available when
enable PrivateLink, a new VPC will be created with subnet `10.52.72.0/24`.

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
  region = "azure-arm::westus"
  tags   = []
}

resource "cloudamqp_privatelink_azure" "privatelink" {
  instance_id = cloudamqp_instance.instance.id
  approved_subscriptions = [
    "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
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
  region  = "azure-arm::westus"
  subnet  = "10.56.72.0/24"
  tags    = []
}

resource "cloudamqp_instance" "instance" {
  name                = "Instance 01"
  plan                = "bunny-1"
  region              = "azure-arm::westus"
  tags                = []
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_privatelink_azure" "privatelink" {
  instance_id = cloudamqp_instance.instance.id
  approved_subscriptions = [
    "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  ]
}
```

</details>

## Argument Reference

* `instance_id`             - (Required) The CloudAMQP instance identifier.
* `approved_subscriptions`  - (Required) Approved subscriptions to access the endpoint service.
                              See format below.
* `sleep`                   - (Optional) Configurable sleep time (seconds) when enable PrivateLink.
                              Default set to 10 seconds.

  ***Note:*** Available from [v1.29.0]

* `timeout`                 - (Optional) Configurable timeout time (seconds) when enable PrivateLink.
                              Default set to 1800 seconds.

  ***Note:*** Available from [v1.29.0]

Approved subscriptions format (GUID): <br>
`XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX`

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource. Will be same as `instance_id`
* `status`- PrivateLink status [enable, pending, disable]
* `service_name` - Service name (alias) of the PrivateLink, needed when creating the endpoint.
* `server_name` - Name of the server having the PrivateLink enabled.

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_privatelink_azure` can be imported using CloudAMQP instance identifier. To retrieve the
identifier, use [CloudAMQP API list intances].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_privatelink_azure.privatelink
  id = cloudamqp_instance.instance.id
}
```

Or use Terraform CLI:

`terraform import cloudamqp_privatelink_azure.privatelink <id>`

`cloudamqp_privatelink_aws` can be imported using CloudAMQP instance identifier.

## Create PrivateLink with additional firewall rules

To create a PrivateLink configuration with additional firewall rules, it's required to chain the
[cloudamqp_security_firewall] resource to avoid parallel conflicting resource calls. You can do this
by making the firewall resource depend on the PrivateLink resource
`cloudamqp_privatelink_azure.privatelink`.

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
  region  = "azure-arm::westus"
  subnet  = "10.56.72.0/24"
  tags    = []
}

resource "cloudamqp_instance" "instance" {
  name                = "Instance 01"
  plan                = "bunny-1"
  region              = "azure-arm::westus"
  tags                = []
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_privatelink_azure" "privatelink" {
  instance_id = cloudamqp_instance.instance.id
  approved_subscriptions = [
    "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
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
    cloudamqp_privatelink_azure.privatelink
   ]
}
```

</details>

[CloudAMQP API list intances]: https://docs.cloudamqp.com/#list-instances
[CloudAMQP plans]: https://www.cloudamqp.com/plans.html
[CloudAMQP PrivateLink]: https://www.cloudamqp.com/docs/cloudamqp-privatelink.html#azure-privatelink
[cloudamqp_security_firewall]: ./security_firewall.md
[cloudamqp_vpc_connect]: ./vpc_connect.md
[v1.29.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.29.0

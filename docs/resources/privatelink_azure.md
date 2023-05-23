---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_privatelink_azure"
description: |-
  Enable PrivateLink for a CloudAMQP instance hosted in Azure.
---

# cloudamqp_privatelink_azure

Enable PrivateLink for a CloudAMQP instance hosted in Azure. If no existing VPC available when enable PrivateLink, a new VPC will be created with subnet `10.52.72.0/24`.

~> **Note:** Enabling PrivateLink will automatically add firewall rules for the peered subnet.
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

Pricing is available at [cloudamqp.com](https://www.cloudamqp.com/plans.html) where you can also find more information about [CloudAMQP PrivateLink](https://www.cloudamqp.com/docs/cloudamqp-privatelink.html#azure-privatelink).

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

resource "cloudamqp_privatelink_azure" "privatelink" {
  instance_id = cloudamqp_instance.instance.id
  approved_subscriptions = [
    "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  ]
}
```
</details>

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier.
* `approved_subscriptions` - (Required) Approved subscriptions to access the endpoint service. See format below.
* `sleep` - (Optional) Configurable sleep time (seconds) when enable PrivateLink. Default set to 60 seconds.
* `timeout` - (Optional) - Configurable timeout time (seconds) when enable PrivateLink. Default set to 3600 seconds.

Approved subscriptions format: <br>
`XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX`

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.
* `status`- PrivateLink status [enable, pending, disable]
* `service_name` - Service name (alias) of the PrivateLink, needed when creating the endpoint.
* `server_name` - Name of the server having the PrivateLink enabled.

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_privatelink_aws` can be imported using CloudAMQP internal identifier.

`terraform import cloudamqp_privatelink_aws.privatelink <id>`

## Create PrivateLink with additional firewall rules

To create a PrivateLink configuration with additional firewall rules, it's required to chain the [cloudamqp_security_firewall](https://registry.terraform.io/providers/cloudamqp/cloudamqp/latest/docs/resources/security_firewall)
resource to avoid parallel conflicting resource calls. You can do this by making the firewall resource depend on the PrivateLink resource, `cloudamqp_privatelink_azure.privatelink`.

Furthermore, since all firewall rules are overwritten, the otherwise automatically added rules for the PrivateLink also needs to be added.

## Example usage with additional firewall rules

<details>
  <summary>
    <b>
      <i>CloudAMQP instance in an existing VPC with managed firewall rules</i>
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

resource "cloudamqp_privatelink_azure" "privatelink" {
  instance_id = cloudamqp_instance.instance.id
  approved_subscriptions = [
    "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  ]
}

resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  rules {
    Description = "Custom PrivateLink setup"
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
    cloudamqp_privatelink_azure.privatelink
   ]
}
```
</details>

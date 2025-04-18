---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_security_firewall"
description: |-
  Configure and manage firewall rules
---

# cloudamqp_security_firewall

This resource allows you to configure and manage firewall rules for the CloudAMQP instance.

~> **WARNING:** Firewall rules applied with this resource will replace any existing firewall rules.
Make sure all wanted rules are present to not lose them.

-> **NOTE:** From [v1.33.0] when destroying this resource the firewall on the servers will also be
removed. I.e. the firewall will be completely closed.

Only available for dedicated subscription plans.

## Example Usage

```hcl
resource "cloudamqp_security_firewall" "this" {
  instance_id = cloudamqp_instance.instance.id

  rules {
    ip        = "192.168.0.0/24"
    ports     = [4567, 4568]
    services  = ["AMQP","AMQPS", "HTTPS"]
  }

  rules {
    ip        = "10.56.72.0/24"
    ports     = []
    services  = ["AMQP","AMQPS", "HTTPS"]
  }

  // Single IP address
  rules {
    ip        = "192.168.1.10/32"
    ports     = []
    services  = ["AMQP","AMQPS", "HTTPS"]
  }
}
```

<details>
  <summary>
    <b>
      <i>Faster instance destroy when running `terraform destroy` from </i>
      <a href="https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.27.0">v1.27.0</a>
    </b>
  </summary>

CloudAMQP Terraform provider [v1.27.0] enables faster `cloudamqp_instance` destroy when running
`terraform destroy`.

```hcl
# Configure the CloudAMQP Provider
provider "cloudamqp" {
  apikey                          = var.cloudamqp_customer_api_key
  enable_faster_instance_destroy  = true
}

resource "cloudamqp_instance" "instance" {
  name    = "terraform-cloudamqp-instance"
  plan    = "penguin-1"
  region  = "amazon-web-services::us-west-1"
  tags    = ["terraform"]
}

resource "cloudamqp_security_firewall" "this" {
  instance_id = cloudamqp_instance.instance.id

  rules {
    ip          = "0.0.0.0/0"
    ports       = []
    services    = ["HTTPS"]
    description = "MGMT interface"
  }

  rules {
    ip          = "10.56.72.0/24"
    ports       = []
    services    = ["AMQP","AMQPS", "HTTPS"]
    description = "VPC subnet"
  }
}
```

</details>

## Argument Reference

Top level argument reference

* `instance_id` - (Required) The CloudAMQP instance ID.
* `rules`       - (Required) An array of rules, minimum of 1 needs to be configured. Each `rules`
                  block consists of the field documented below.
* `sleep`       - (Optional) Configurable sleep time in seconds between retries for firewall
                  configuration. Default set to 30 seconds.
* `timeout`     - (Optional) Configurable timeout time in seconds for firewall configuration.
                  Default set to 1800 seconds.

___

The `rules` block consists of:

* `ip`          - (Required) CIDR address: IP address with CIDR notation (e.g. 10.56.72.0/24)
* `ports`       - (Optional) Custom ports to be opened
* `services`    - (Required) Pre-defined service ports, see table below
* `description` - (Optional) Description name of the rule. e.g. Default.

Pre-defined services for RabbitMQ:

| Service name | Port  |
|--------------|-------|
| AMQP         | 5672  |
| AMQPS        | 5671  |
| HTTPS        | 443   |
| MQTT         | 1883  |
| MQTTS        | 8883  |
| STOMP        | 61613 |
| STOMPS       | 61614 |
| STREAM       | 5552  |
| STREAM_SSL   | 5551  |

Pre-defined services for LavinMQ:

| Service name | Port  |
|--------------|-------|
| AMQP         | 5672  |
| AMQPS        | 5671  |
| HTTPS        | 443   |
| MQTT         | 1883  |
| MQTTS        | 8883  |

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

If used together with [VPC GPC peering], see additional information.

## Import

`cloudamqp_security_firewall` can be imported using CloudAMQP instance identifier. To
retrieve the identifier, use [CloudAMQP API list intances].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_security_firewall.firewall
  id = cloudamqp_instance.instance.id
}
```

Or use Terraform CLI:

`terraform import cloudamqp_security_firewall.firewall <instance_id>`

## Destroy the resource

From [v1.33.0] when destroying this resource the firewall on the servers will be removed. I.e. the
firewall will be completly closed.

Older version will instead update the firewall with a default rule.

```hcl
rules {
  ip          = "0.0.0.0/0"
  ports       = []
  services    = ["AMQP", "AMQPS", "STOMP", "STOMPS", "MQTT", "MQTTS", "HTTPS", "STREAM", "STREAM_SSL"]
  description = "Default"
}
```

## Enable faster instance destroy

When running `terraform destroy` this resource will try configure the firewall with default rules
before deleting `cloudamqp_instance`. This is not necessary since the servers will be deleted.

Set `enable_faster_instance_destroy` to ***true*** in the provider configuration to skip this.

## Known issues

<details>
  <summary>Custom ports trigger new update every time</summary>

  Before release v1.15.1 using the custom ports can cause a missmatch upon reading data and
  trigger a new update every time.

  Reason is that there is a bug in validating the response from the underlying API.

  Update the provider to at least [v1.15.1] to fix the issue.
 </details>

<details>
  <summary>Using pre-defined service port in ports</summary>

Using one of the port from the pre-defined services in ports argument, see example of using port
5671 instead of the service *AMQPS*.

```hcl
resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  rules {
    ip        = "192.168.0.0/24"
    ports     = [5671]
    services  = []
  }
}
```

Will still create the firewall rule for the instance, but will trigger a new update each `plan` or
`apply`. Due to a missmatch between state file and underlying API response.

To solve this, edit the configuration file and change port 5671 to service *AMQPS* and run
`terraform apply -refresh-only` to only update the state file and remove the missmatch.

```hcl
resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  rules {
    ip        = "192.168.0.0/24"
    ports     = []
    services  = ["AMQPS"]
  }
}
```

The provider from [v1.15.2] will start to warn about using this.

 </details>

[CloudAMQP API list intances]: https://docs.cloudamqp.com/#list-instances
[v1.15.1]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.15.1
[v1.15.2]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.15.2
[v1.27.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.27.0
[v1.33.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.33.0
[VPC GPC peering]: ./vpc_gcp_peering#create-vpc-peering-with-additional-firewall-rules

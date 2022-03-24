---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_security_firewall"
description: |-
  Configure and manage firewall rules
---

# cloudamqp_security_firewall

This resource allows you to configure and manage firewall rules for the CloudAMQP instance. Beware that all rules need to be present, since all older configurations will be overwritten.

Only available for dedicated subscription plans.

## Example Usage

```hcl
resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  rules {
    ip          = "192.168.0.0/24"
    ports       = [4567, 4568]
    services    = ["AMQP","AMQPS", "HTTPS"]
  }

  rules {
    ip          = "10.56.72.0/24"
    ports       = []
    services    = [AMQP","AMQPS", "HTTPS"]
  }
}
```

## Argument Reference

Top level argument reference

* `instance_id` - (Required) The CloudAMQP instance ID.
* `rules`       - (Required) An array of rules, minimum of 1 needs to be configured. Each `rules` block consists of the field documented below.

___

The `rules` block consists of:

* `ip`          - (Required) Source ip and netmask for the rule. (e.g. 10.56.72.0/24)
* `ports`       - (Optional) Custom ports to be opened
* `services`    - (Required) Pre-defined service ports, see table below
* `description` - (Optional) Description name of the rule. e.g. Default.

Pre-defined services:

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

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

If used together with [VPC GPC peering](https://registry.terraform.io/providers/cloudamqp/cloudamqp/latest/docs/resources/vpc_gcp_peering#create-vpc-peering-with-additional-firewall-rules), see additional information.

## Import

`cloudamqp_security_firewall` can be imported using CloudAMQP instance identifier.

`terraform import cloudamqp_security_firewall.firewall <instance_id>`

## Known issues

<details>
  <summary>Custom ports trigger new update every time</summary>

  Before release [v1.15.1](https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.15.1) using the custom ports can cause a missmatch upon reading data and trigger a new update every time.

  Reason is that there is a bug in validating the response from the underlying API.

  Update the provider to at least v1.15.1 to fix the issue.
 </details>

<details>
  <summary>Using pre-defined service port in ports</summary>

Using one of the port from the pre-defined services in ports argument, see example of using port 5671 instead of the service *AMQPS*.

```hcl
resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  rules {
    ip          = "192.168.0.0/24"
    ports       = [5671]
    services    = []
  }
}
```

Will still create the firewall rule for the instance, but will trigger a new update each `plan` or `apply`. Due to a missmatch between state file and underlying API response.

To solve this, edit the configuration file and change port 5671 to service *AMQPS* and run `terraform apply -refresh-only` to only update the state file and remove the missmatch.

```hcl
resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  rules {
    ip          = "192.168.0.0/24"
    ports       = []
    services    = ["AMQPS"]
  }
}
```

The provider from [v1.15.2](https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.16.0) will start to warn about using this.
 </details>

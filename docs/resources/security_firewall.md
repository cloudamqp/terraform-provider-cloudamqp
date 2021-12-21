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
* `services`    - (Required) Pre-defined service ports
* `description` - (Optional) Description name of the rule. e.g. Default.

Supported services: *AMQP*, *AMQPS*, *HTTPS*, *MQTT*, *MQTTS*, *STOMP*, *STOMPS*, *STREAM*, *STREAM\_SSL*

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

If used together with [VPC GPC peering](https://registry.terraform.io/providers/cloudamqp/cloudamqp/latest/docs/resources/vpc_gcp_peering#create-vpc-peering-with-additional-firewall-rules), see additional information.

## Import

`cloudamqp_security_firewall` can be imported using CloudAMQP instance identifier.

`terraform import cloudamqp_security_firewall.firewall <instance_id>`

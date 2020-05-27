---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_security_firewall"
description: |-
  Configure and manage firewall rules
---

# cloudamqp_security_firewall

This resource allows you to configure and manage firewall rules for the CloudAMQP instance. Beware that all rules needs to present, since all older configurations will be overwritten. Depends on `cloudamqp_instance`resource and the instance identifier.

Only available for dedicated subscription plans.

## Example Usage

```hcl
resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  rules {
    ip          = "192.168.0.0/24"
    ports       = [4567, 4568]
    services    = ["AMQP","AMQPS"]
  }

  rules {
    ip          = "10.56.72.0/24"
    ports       = []
    services    = ["AMQP","AMQPS"]
  }
}
```

## Argument Reference

Top level argument reference

* `instance_id` - (Required) The CloudAMQP instance ID.
* `rules`       - (Required) An array of rules, minimum of 1 needs to be configured. Each `rules` block consists of the filed documented below.

___

The `rules` block consists of:

* `ip`        - (Required) Source ip and netmask for the rule. (e.g. 10.56.72.0/24)
* `ports`     - (Optional) Custom ports to be opened
* `services`  - (Required) Pre-defined service ports

Supported services: *AMQP*, *AMQPS*, *MQTT*, *MQTTS*, *STOMP*, *STOMPS*

## Import

`cloudamqp_security_firewall` can be imported using CloudAMQP instance identifier.

`terraform import cloudamqp_security_firewall.firewall <instance_id>`

---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_security_firewall"
description: |-
  Get information of firewall rules
---

# cloudamqp_security_firewall

Use this data source to retrieve information about the firewall rules that are open.

## Example Usage

```hcl
data "cloudamqp_security_firewall" "firewall" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attributes reference

All attributes reference are computed

* `rules`       - An array of firewall rules, each `rules` block consists of the field documented below.

___

The `rules` block consists of:

* `ip`          - CIDR address: IP address with CIDR notation (e.g. 10.56.72.0/24)
* `ports`       - Custom ports opened
* `services`    - Pre-defined service ports, see table below
* `description` - Description name of the rule. e.g. Default.

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

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

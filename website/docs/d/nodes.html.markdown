---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_nodes"
description: |-
  Get information about the node(s) in the CloudAMQP instance.
---

# cloudamqp_alarm

Use this data source to retrieve information about the node(s) created by CloudAMQP instance. Require to know the identifier of the corresponding `cloudamqp_instance`resource or data source.

## Eample Usage

```hcl
data "cloudamqp_nodes" "nodes" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attribute reference

* `nodes`               - (Computed) Elem array of nodes
  * `hostname`          - (Computed) Hostname assigned to the node.
  * `name`              - (Computed) Name of the node.
  * `running`           - (Computed) Is the node running?
  * `rabbitmq_version`  - (Computed) Currently configured Rabbit MQ version on the node.
  * `erlang_version`    - (Computed) Currently used Erlanbg version on the node.
  * `hipe`              - (Computed) Enable or disable High-performance Erlang.

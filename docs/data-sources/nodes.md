---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_nodes"
description: |-
  Get information about the node(s) in the CloudAMQP instance.
---

# cloudamqp_nodes

Use this data source to retrieve information about the node(s) created by CloudAMQP instance.

## Example Usage

```hcl
data "cloudamqp_nodes" "nodes" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attributes Reference

All attributes reference are computed

* `id`    - The identifier for this resource.
* `nodes` - An array of node information. Each `nodes` block consists of the fields documented below.

___

The `nodes` block consist of

* `hostname`              - External hostname assigned to the node.
* `hostname_internal`     - Internal hostname assigned to the node.
* `name`                  - Name of the node.
* `running`               - Is the node running?
* `rabbitmq_version`      - Currently configured Rabbit MQ version on the node.
* `erlang_version`        - Currently used Erlang version on the node.
* `hipe`                  - Enable or disable High-performance Erlang.
* `configured`            - Is the node configured?
* `disk_size`             - Subscription plan disk size
* `additional_disk_size`  - Additional added disk size
* `availability_zone`     - Availability zone the node is hosted in.

  ***Note:*** Total disk size = disk_size + additional_disk_size

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

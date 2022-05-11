---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_node_actions"
description: |-
  Invoke actions on a specific node (e.g. restart RabbitMQ).
---

# cloudamqp_node_actions

This resource allows you to invoke actions on a specific node.

Only available for dedicated subscription plans.

## Example Usage

Already know the node identifier (e.g. from state file)

```hcl
# New recipient to receieve notifications
resource "cloudamqp_node_actions" "node_action" {
  instance_id = cloudamqp_instance.instance.id
  node_id = <node_id>
  action = "restart"
}
```

Using data source `cloudamqp_nodes`

```hcl
data "cloudamqp_nodes" "list_nodes" {
  instance_id = cloudamqp_instance.instance.id
}

resource "cloudamqp_node_actions" "node_action" {
  instance_id = cloudamqp_instance.instance.id
  node_id = data.cloudamqp_nodes.list_nodes.nodes[0].node_id
  action = "restart"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id`   - (Required) The CloudAMQP instance ID.
* `node_id`       - (Required) The node ID.
* `action`        - (Required) Endpoint to send the notification.

## Attributes Reference

All attributes reference are computed

* `id`      - The identifier for this resource.
* `running` - If the node is running.

## Action reference

Valid options for action.

| Action       | Info                               |
|--------------|------------------------------------|
| start        | Start RabbitMQ                     |
| stop         | Stop RabbitMQ                      |
| restart      | Restart RabbitMQ                   |
| reboot       | Reboot the node                    |
| mgmt.restart | Restart the RabbitMQ mgmt interace |

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id` and node identifier.

## Import

This resource cannot be imported.

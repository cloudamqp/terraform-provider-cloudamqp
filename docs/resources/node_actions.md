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

<details>
  <summary>
    <b>
      <i>Already know the node identifier (e.g. from state file)</i>
    </b>
  </summary>

```hcl
# New recipient to receieve notifications
resource "cloudamqp_node_actions" "node_action" {
  instance_id = cloudamqp_instance.instance.id
  node_name   = "<node name>"
  action      = "restart"
}
```

</details>

<details>
  <summary>
    <b>
      <i>Multi node RabbitMQ restart</i>
    </b>
  </summary>

Using data source `cloudamqp_nodes` to restart RabbitMQ on all nodes.

-> **Note:** RabbitMQ restart on multiple nodes need to be chained, let one node restart at the time.

```hcl
data "cloudamqp_nodes" "list_nodes" {
  instance_id = cloudamqp_instance.instance.id
}

resource "cloudamqp_node_actions" "restart_01" {
  instance_id = cloudamqp_instance.instance.id
  action      = "restart"
  node_name   = data.cloudamqp_nodes.list_nodes.nodes[0].name
}

resource "cloudamqp_node_actions" "restart_02" {
  instance_id = cloudamqp_instance.instance.id
  action      = "restart"
  node_name   = data.cloudamqp_nodes.list_nodes.nodes[1].name

  depends_on = [
    cloudamqp_node_actions.restart_01,
  ]
}

resource "cloudamqp_node_actions" "restart_03" {
  instance_id = cloudamqp_instance.instance.id
  action      = "restart"
  node_name   = data.cloudamqp_nodes.list_nodes.nodes[2].name

  depends_on = [
    cloudamqp_node_actions.restart_01,
    cloudamqp_node_actions.restart_02,
  ]
}

```

</details>

<details>
  <summary>
    <b>
      <i>Combine log level configuration change with multi node RabbitMQ restart</i>
    </b>
  </summary>

```hcl
data "cloudamqp_nodes" "list_nodes" {
  instance_id = cloudamqp_instance.instance.id
}

resource "cloudamqp_rabbitmq_configuration" "rabbitmq_config" {
  instance_id         = cloudamqp_instance.instance.id
  log_exchange_level  = "info"
}

resource "cloudamqp_node_actions" "restart_01" {
  instance_id = cloudamqp_instance.instance.id
  action      = "restart"
  node_name   = data.cloudamqp_nodes.list_nodes.nodes[0].name

  depends_on = [
    cloudamqp_rabbitmq_configuration.rabbitmq_config,
  ]
}

resource "cloudamqp_node_actions" "restart_02" {
  instance_id = cloudamqp_instance.instance.id
  action      = "restart"
  node_name   = data.cloudamqp_nodes.list_nodes.nodes[1].name

  depends_on = [
    cloudamqp_rabbitmq_configuration.rabbitmq_config,
    cloudamqp_node_actions.restart_01,
  ]
}

resource "cloudamqp_node_actions" "restart_03" {
  instance_id = cloudamqp_instance.instance.id
  action      = "restart"
  node_name   = data.cloudamqp_nodes.list_nodes.nodes[2].name
  
  depends_on = [
    cloudamqp_rabbitmq_configuration.rabbitmq_config,
    cloudamqp_node_actions.restart_01,
    cloudamqp_node_actions.restart_02,
  ]
}

```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id`   - (Required) The CloudAMQP instance ID.
* `node_name`     - (Required) The node name, e.g `green-guinea-pig-01`.
* `action`        - (Required) The action to invoke on the node.

## Attributes Reference

All attributes reference are computed

* `id`      - The identifier for this resource.
* `running` - If the node is running.

## Action reference

Valid actions for ***LavinMQ***.

| Action       | Info                               |
|--------------|------------------------------------|
| start        | Start LavinMQ                      |
| stop         | Stop LavinMQ                       |
| restart      | Restart LavinMQ                    |
| reboot       | Reboot the node                    |

Valid actions for ***RabbitMQ***.

| Action       | Info                               |
|--------------|------------------------------------|
| start        | Start RabbitMQ                     |
| stop         | Stop RabbitMQ                      |
| restart      | Restart RabbitMQ                   |
| reboot       | Reboot the node                    |
| mgmt.restart | Restart the RabbitMQ mgmt interace |

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id` and node
name.

## Import

This resource cannot be imported.

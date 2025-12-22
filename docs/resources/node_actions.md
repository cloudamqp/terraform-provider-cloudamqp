---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_node_actions"
description: |-
  Invoke actions on specific nodes or entire cluster (e.g. restart RabbitMQ).
---

# cloudamqp_node_actions

This resource allows you to invoke actions on specific nodes or the entire cluster. Actions can target individual nodes, multiple nodes, or all nodes in the cluster at once.

Only available for dedicated subscription plans.

-> **Note:** From version 1.41.0, this resource supports cluster-level actions (`cluster.start`, `cluster.stop`, `cluster.restart`) and the `node_names` list attribute for targeting multiple nodes. The `node_name` attribute is deprecated in favor of `node_names`.

## Example Usage

<details>
  <summary>
    <b>
      <i>Cluster-wide broker restart (recommended for v1.41.0+)</i>
    </b>
  </summary>

Restart the broker on all nodes of the cluster at once. Making sure the broker is stopped and started in correct order. This is the simplest approach for cluster-wide operations.

```hcl
resource "cloudamqp_node_actions" "cluster_restart" {
  instance_id = cloudamqp_instance.instance.id
  action      = "cluster.restart"
}
```

</details>

<details>
  <summary>
    <b>
      <i>Restart broker on specific nodes using node_names</i>
    </b>
  </summary>

Target specific nodes using the `node_names` list attribute.

```hcl
data "cloudamqp_nodes" "nodes" {
  instance_id = cloudamqp_instance.instance.id
}

resource "cloudamqp_node_actions" "restart_subset" {
  instance_id = cloudamqp_instance.instance.id
  action      = "restart"
  node_names  = [
    data.cloudamqp_nodes.nodes.nodes[0].name,
    data.cloudamqp_nodes.nodes.nodes[1].name
  ]
}
```

</details>

<details>
  <summary>
    <b>
      <i>Reboot a single node</i>
    </b>
  </summary>

Reboot the entire node (VM) rather than just the broker.

```hcl
data "cloudamqp_nodes" "nodes" {
  instance_id = cloudamqp_instance.instance.id
}

resource "cloudamqp_node_actions" "reboot_node" {
  instance_id = cloudamqp_instance.instance.id
  action      = "reboot"
  node_names  = [data.cloudamqp_nodes.nodes.nodes[0].name]
}
```

</details>

<details>
  <summary>
    <b>
      <i>Restart RabbitMQ management interface</i>
    </b>
  </summary>

Only restart the management interface without affecting the broker.

```hcl
data "cloudamqp_nodes" "nodes" {
  instance_id = cloudamqp_instance.instance.id
}

resource "cloudamqp_node_actions" "mgmt_restart" {
  instance_id = cloudamqp_instance.instance.id
  action      = "mgmt.restart"
  node_names  = [data.cloudamqp_nodes.nodes.nodes[0].name]
}
```

</details>

<details>
  <summary>
    <b>
      <i>Combine with configuration changes</i>
    </b>
  </summary>

Apply configuration changes and restart the cluster.

```hcl
resource "cloudamqp_rabbitmq_configuration" "rabbitmq_config" {
  instance_id        = cloudamqp_instance.instance.id
  log_exchange_level = "info"
}

resource "cloudamqp_node_actions" "cluster_restart" {
  instance_id = cloudamqp_instance.instance.id
  action      = "cluster.restart"

  depends_on = [
    cloudamqp_rabbitmq_configuration.rabbitmq_config,
  ]
}
```

</details>

<details>
  <summary>
    <b>
      <i>Legacy Usage (pre-1.41.0)</i>
    </b>
  </summary>

These examples show the older approach using `node_name` (singular) and chained restarts. While still supported, the cluster-level actions above are recommended for new configurations.

**Single node restart:**

```hcl
resource "cloudamqp_node_actions" "node_action" {
  instance_id = cloudamqp_instance.instance.id
  node_name   = "<node name>"
  action      = "restart"
}
```

**Chained multi-node restart:**

-> **Note:** This approach restarts nodes sequentially to minimize cluster disruption. Consider using `cluster.restart` for simpler configuration.

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

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance ID.
* `node_names`  - (Optional) List of node names to perform the action on, e.g. `["green-guinea-pig-01", "green-guinea-pig-02"]`. For cluster-level actions (`cluster.start`, `cluster.stop`, `cluster.restart`), this can be omitted and the action will automatically apply to all nodes.
* `node_name`   - (Optional, Deprecated) The node name, e.g. `green-guinea-pig-01`. Use `node_names` instead. This attribute will be removed in a future version.
* `action`      - (Required) The action to invoke. See [Action reference](#action-reference) below for valid values.
* `sleep`       - (Optional) Sleep interval in seconds between polling for node status. Default: `10`.
* `timeout`     - (Optional) Timeout in seconds for the action to complete. Default: `1800` (30 minutes).

-> **Note:** Either `node_name` or `node_names` must be specified for non-cluster actions. Cluster actions (`cluster.start`, `cluster.stop`, `cluster.restart`) can omit both and will automatically target all nodes.

## Attributes Reference

All attributes reference are computed

* `id`      - The identifier for this resource.
* `running` - If the node is running.

## Action reference

Actions are categorized by what they affect:

### Broker Actions

These actions control the message broker software (RabbitMQ or LavinMQ) on the specified nodes.

| Action  | Info                                      | Applies to        |
|---------|-------------------------------------------|-------------------|
| start   | Start the message broker                  | RabbitMQ, LavinMQ |
| stop    | Stop the message broker                   | RabbitMQ, LavinMQ |
| restart | Restart the message broker                | RabbitMQ, LavinMQ |

### Management Interface Actions

These actions control the management interface without affecting the broker itself.

| Action       | Info                                      | Applies to |
|--------------|-------------------------------------------|------------|
| mgmt.restart | Restart the RabbitMQ management interface | RabbitMQ   |

### Node Actions

These actions affect the entire node (VM), not just the broker software.

| Action | Info                                          | Applies to        |
|--------|-----------------------------------------------|-------------------|
| reboot | Reboot the entire node (VM)                   | RabbitMQ, LavinMQ |

### Cluster Actions

-> **Available from version 1.41.0**

These actions operate on all nodes in the cluster simultaneously. The `node_names` attribute can be omitted for these actions.

| Action          | Info                                            | Applies to        |
|-----------------|-------------------------------------------------|-------------------|
| cluster.start   | Start the message broker on all cluster nodes   | RabbitMQ, LavinMQ |
| cluster.stop    | Stop the message broker on all cluster nodes    | RabbitMQ, LavinMQ |
| cluster.restart | Restart the message broker on all cluster nodes | RabbitMQ, LavinMQ |

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`. For non-cluster actions, it also requires either `node_name` or `node_names` to specify which nodes to act upon. Cluster-level actions automatically apply to all nodes in the cluster.

## Import

This resource cannot be imported.

---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_rabbitmq_configuration"
description: |-
  Update Rabbit MQ config
---

# cloudamqp_rabbitmq_configuration

This resource allows you update RabbitMQ config.

Only available for dedicated subscription plans running ***RabbitMQ***.



## Example Usage

<details>
  <summary>
    <b>
      <i>RabbitMQ configuration and using 0 values</i>
    </b>
  </summary>

From [v1.35.0] and migrating this resource to Terraform plugin Framework.
It's now possible to use 0 values in the configuration.

```hcl
resource "cloudamqp_rabbitmq_configuration" "rabbitmq_config" {
  instance_id = cloudamqp_instance.instance.id
  heartbeat   = 0
}
```

</details>

<details>
  <summary>
    <b>
      <i>RabbitMQ configuration with default values</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_rabbitmq_configuration" "rabbitmq_config" {
  instance_id = cloudamqp_instance.instance.id
  channel_max                   = 0
  connection_max                = -1
  consumer_timeout              = 7200000
  heartbeat                     = 120
  log_exchange_level            = "error"
  max_message_size              = 134217728
  queue_index_embed_msgs_below  = 4096
  vm_memory_high_watermark      = 0.81
  cluster_partition_handling    = "autoheal"
}
```

</details>

<details>
  <summary>
    <b>
      <i>Change log level and combine `cloudamqp_node_actions` for RabbitMQ restart</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_rabbitmq_configuration" "rabbitmq_config" {
  instance_id = cloudamqp_instance.instance.id
  channel_max                   = 0
  connection_max                = -1
  consumer_timeout              = 7200000
  heartbeat                     = 120
  log_exchange_level            = "info"
  max_message_size              = 134217728
  queue_index_embed_msgs_below  = 4096
  vm_memory_high_watermark      = 0.81
  cluster_partition_handling    = "autoheal"
}

data "cloudamqp_nodes" "list_nodes" {
  instance_id = cloudamqp_instance.instance.id
}

resource "cloudamqp_node_actions" "node_action" {
  instance_id = cloudamqp_instance.instance.id
  node_name   = data.cloudamqp_nodes.list_nodes.nodes[0].name
  action      = "restart"

  depends_on = [
    cloudamqp_rabbitmq_configuration.rabbitmq_config,
  ]
}
```

</details>

<details>
  <summary>
    <b>
      <i>
        Only change log level for exchange. All other values will be read from the RabbitMQ
        configuration.
      </i>
    </b>
  </summary>

```hcl
resource "cloudamqp_rabbitmq_configuration" "rabbit_config" {
  instance_id         = cloudamqp_instance.instance.id
  log_exchange_level  = "info"
}
```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id`                   - (Required) The CloudAMQP instance ID.
* `heartbeat`                     - (Optional/Computed) Set the server AMQP 0-9-1 heartbeat timeout
                                    in seconds.
* `connection_max`                - (Optional/Computed) Set the maximum permissible number of
                                    connection.
* `channel_max`                   - (Optional/Computed) Set the maximum permissible number of
                                    channels per connection.
* `consumer_timeout`              - (Optional/Computed) A consumer that has recevied a message and
                                    does not acknowledge that message within the timeout in
                                    milliseconds
* `vm_memory_high_watermark`      - (Optional/Computed) When the server will enter memory based
                                    flow-control as relative to the maximum available memory.
* `queue_index_embed_msgs_below`  - (Optional/Computed) Size in bytes below which to embed messages
                                    in the queue index. 0 will turn off payload embedding in the
                                    queue index.
* `max_message_size`              - (Optional/Computed) The largest allowed message payload size in
                                    bytes.
* `log_exchange_level`            - (Optional/Computed) Log level for the logger used for log
                                    integrations and the CloudAMQP Console log view.
* `cluster_partition_handling`    - (Optional/Computed) Set how the cluster should handle network
                                    partition.
* `sleep`                         - (Optional) Configurable sleep time in seconds between retries
                                    for RabbitMQ configuration. Default set to 60 seconds.
* `timeout`                       - (Optional) - Configurable timeout time in seconds for RabbitMQ
                                    configuration. Default set to 3600 seconds.

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Argument threshold values

### heartbeat

| Type | Default | Min  | Affect |
|---|---|---|---|
| int | 120 | 0 | Only effects new connection |

### connection_max

| Type | Default | Min  | Affect |
|---|---|---|---|
| int | -1 | 1 | Applied immediately (RabbitMQ restart required before 3.11.13) |

Note: -1 in the provider corresponds to INFINITY in the RabbitMQ config

### channel_max

| Type | Default | Min | Affect |
|---|---|---|---|
| int | 128 | 0 | Only effects new connections |

Note: 0 means "no limit"

### consumer_timeout

| Type | Default | Min | Max | Unit | Affect |
|---|---|---|---|---|---|
| int | 7200000 | 10000 | 86400000 | milliseconds | Only effects new channels |

Note: -1 in the provider corresponds to false (disable) in the RabbitMQ config

### vm_memory_high_watermark

| Type | Default | Min | Max | Affect |
|---|---|---|---|---|
 | float | 0.81 | 0.4 | 0.9 | Applied immediately |

### queue_index_embed_msgs_below

| Type | Default | Min | Max | Unit | Affect |
|---|---|---|---|---|---|
| int | 4096 | 0 | 10485760 | bytes | Applied immediately for new queues |

Note: Existing queues requires restart

### max_message_size

| Type | Default | Min | Max | Unit | Affect |
|---|---|---|---|---|---|
| int | 134217728 | 1 | 536870912 | bytes | Only effects new channels |

### log_exchange_level

| Type | Default | Affect |
|---|---|---|
| string | error | RabbitMQ restart required |

Note: `debug, info, warning, error, critical, none`

### cluster_partition_handling

| Type  | Affect | Note |
|---|---|---|
| string | Applied immediately | `autoheal, pause_minority, ignore` |

Recommended setting for cluster_partition_handling: `autoheal` for cluster with 1-2
nodes, `pause_minority` for cluster with 3 or more nodes. While `ignore` setting is not recommended.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_rabbitmq_configuration` can be imported using the CloudAMQP instance identifier.  To
retrieve the identifier, use [CloudAMQP API list intances].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_rabbitmq_configuration.config
  id = cloudamqp_instance.instance.id
}
```

Or use Terraform CLI:

`terraform import cloudamqp_rabbitmq_configuration.config <instance_id>`

## Known issues

<details>
  <summary>Cannot set heartbeat=0 when creating this resource</summary>

-> **Note:** This is no longer the case from [v1.35.0].

The provider is built by older `Terraform Plugin SDK` which doesn't support nullable configuration
values. Instead the values will be set to it's default value based on it's schema primitive type.

* schema.TypeString = ""
* schema.TypeInt = 0
* schema.TypeFloat = 0.0
* schema.TypeBool = false

During initial create of this resource, we need to exclude all arguments that can take these default
values. Argument such as `hearbeat`, `channel_max`, etc. cannot be set to its default value, 0 in
these cases. Current workaround is to use the default value in the initial create run, then change
to the wanted value in the re-run.

Will be solved once we migrate the current provider to `Terraform Plugin Framework`.

</details>

[CloudAMQP API list intances]: https://docs.cloudamqp.com/index.html#tag/instances/get/instances
[v1.35.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.35.0

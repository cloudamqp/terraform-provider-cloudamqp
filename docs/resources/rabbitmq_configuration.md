---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_rabbitmq_configuration"
description: |-
  Update Rabbit MQ config
---

# cloudamqp_rabbitmq_configuration

This resource allows you update RabbitMQ config.

Only available for dedicated subscription plans.

## Example Usage

<details>
  <summary>
    <b>
      <i>RabbitMQ configuration with default values</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_rabbitmq_configuration" "rabbitmq_config" {
  instance_id = cloudamqp_instance.instance.id
  channel_max = 0
  connection_max = -1
  consumer_timeout = 7200000
  heartbeat = 120
  log_exchange_level = "error"
  max_message_size = 134217728
  queue_index_embed_msgs_below = 4096
  vm_memory_high_watermark = 0.81
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
  channel_max = 0
  connection_max = -1
  consumer_timeout = 7200000
  heartbeat = 120
  log_exchange_level = "info"
  max_message_size = 134217728
  queue_index_embed_msgs_below = 4096
  vm_memory_high_watermark = 0.81
}

data "cloudamqp_nodes" "list_nodes" {
  instance_id = cloudamqp_instance.instance.id
}

resource "cloudamqp_node_actions" "node_action" {
  instance_id = cloudamqp_instance.instance.id
  node_id = data.cloudamqp_nodes.list_nodes.nodes[0].node_id
  action = "restart"

  depends_on = [
    cloudamqp_rabbitmq_configuration.rabbitmq_config,
  ]
}
```
</details>

<details>
  <summary>
    <b>
      <i>Only change log level for exchange. All other values will be set to default, see table below.</i>
    </b>
  </summary>


```hcl
resource "cloudamqp_rabbitmq_configuration" "rabbit_config" {
  instance_id = cloudamqp_instance.instance.id
  log_exchange_level = "info"
}
```
</details>

## Argument Reference

The following arguments are supported:

* `instance_id`                   - (Required) The CloudAMQP instance ID.
* `heartbeat`                     - (Optional) Set the server AMQP 0-9-1 heartbeat timeout in seconds.
* `connection_max`                - (Optional) Set the maximum permissible number of connection.
* `channel_max`                   - (Optional) Set the maximum permissible number of channels per connection.
* `consumer_timeout`              - (Optional) A consumer that has recevied a message and does not acknowledge that message within the timeout in milliseconds
* `vm_memory_high_watermark`      - (Optional) When the server will enter memory based flow-control as relative to the maximum available memory.
* `queue_index_embed_msgs_below`  - (Optional) Size in bytes below which to embed messages in the queue index.
* `max_message_size`              - (Optional) The largest allowed message payload size in bytes.
* `log_exchange_level`            - (Optional) Log level for the logger used for log integrations and the CloudAMQP Console log view.

  ***Note: Requires a restart of RabbitMQ to be applied.***

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Argument threshold values

| Argument                     | Type   | Default   | Min   | Max       | Note                                                              |
|------------------------------|--------|-----------|-------|-----------|-------------------------------------------------------------------|
| heartbeat                    | int    | 120       | 1     | -         |                                                                   |
| connection_max               | int    | -1        | 1     | -         | -1 in the provider corresponds to INFINITY in the RabbitMQ config |
| channel_max                  | int    | 0         | 0     | -         | 0 means "no limit"                                                |
| consumer_timeout             | int    | 7200000   | 10000 | 25000000  | Timeout in milliseconds                                           |
| vm_memory_high_watermark     | float  | 0.81      | 0.4   | 0.9       |                                                                   |
| queue_index_embed_msgs_below | int    | 4096      | 1     | 10485760  |                                                                   |
| max_message_size             | int    | 134217728 | 1     | 536870912 | Size in bytes                                                     |
| log_exchange_level           | string | error     | -     | -         | debug, info, warning, error, critical                             |

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_rabbitmq_configuration` can be imported using the CloudAMQP instance identifier.

`terraform import cloudamqp_rabbitmq_configuration.config <instance_id>`

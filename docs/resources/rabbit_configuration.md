---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_rabbit_configuration"
description: |-
  Enable and disable Rabbit MQ plugin.
---

# cloudamqp_rabbit_configuration

This resource allows you update RabbitMQ config. The resource needs to first be imported once the CloudAMQP instance have been created and will present the current RabbitMQ config.

Only available for dedicated subscription plans.

## Example Usage

```hcl
resource "cloudamqp_rabbit_configuration" "rabbit_config" {
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

## Argument Reference

The following arguments are supported:

* `instance_id`                   - (Required) The CloudAMQP instance ID.
* `heartbeat`                     - (Computed/Optional) Set the server AMQP 0-9-1 heartbeat timeout in seconds.
* `connection_max`                - (Computed/Optional) Set the maximum permissible number of connection.
* `channel_max`                   - (Computed/Optional) Set the maximum permissible number of channels per connection.
* `consumer_timeout`              - (Computed/Optional) A consumer that has recevied a message and does not acknowledge that message within the timeout in milliseconds
* `vm_memory_high_watermark`      - (Computed/Optional) When the server will enter memory based flow-control as relative to the maximum available memory.
* `queue_index_embed_msgs_below`  - (Computed/Optional) Size in bytes below which to embed messages in the queue index.
* `max_message_size`              - (Computed/Optional) The largest allowed message payload size in bytes.
* `log_exchange_level`            - (Computed/Optional) Log level for the logger used for log integrations and the CloudAMQP Console log view.

  ***Note: Requires a RabbitMQ restart to be applied.***

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Argument threshold values

| Argument                     | Type   | Default   | Min   | Max       | Note                                                              |
|------------------------------|--------|-----------|-------|-----------|-------------------------------------------------------------------|
| heartbeat                    | int    | 120       | 0     | -         |                                                                   |
| connection_max               | int    | -1        | 1     | -         | -1 in the provider corresponds to INFINITY in the RabbitMQ config |
| channel_max                  | int    | 128       | 0     | -         | 0 means "no limit"                                                |
| consumer_timeout             | int    | 900000    | 10000 | 25000000  | Timeout in milliseconds                                           |
| vm_memory_high_watermark     | float  | 0.81      | 0.4   | 0.9       |                                                                   |
| queue_index_embed_msgs_below | int    | 4096      | 1     | 10485760  |                                                                   |
| max_message_size             | int    | 134217728 | 1     | 536870912 | Size in bytes                                                     |
| log_exchange_level           | string | error     | -     | -         | debug, info, warning, error, critical                             |

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_rabbit_configuration` must be imported using CloudAMQP instance identifier. The RabbitMQ config can then be updated with new configuration.

`terraform import cloudamqp_rabbit_configuration.rabbit_config <instance_id>`

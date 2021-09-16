---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_webhook"
description: |-
  Enable or disable webhook for a vhost and queue
---

# cloudamqp_webhook

This resource allows you to enable or disable webhooks for a specific vhost and queue.

Only available for dedicated subscription plans.

## Example Usage

```hcl
resource "cloudamqp_webhook" "webhook_queue" {
  instance_id = cloudamqp_instance.instance.id
  vhost = "myvhost"
  queue = "webhook-queue"
  webhook_uri = "https://example.com/webhook?key=secret"
  retry_interval = 5
  concurrency = 5
}
```

## Argument Reference

The following arguments are supported:

* `instance_id`     - (Required) The CloudAMQP instance ID.
* `vhost`           - (Required) The vhost the queue resides in.
* `queue`           - (Required) A (durable) queue on your RabbitMQ instance.
* `webhook_uri`     - (Required) A POST request will be made for each message in the queue to this endpoint.
* `retry_interval`  - (Required) How often we retry if your endpoint fails (in seconds).
* `concurrency`     - (Required) Max simultaneous requests to the endpoint.

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

Note: If the endpoint returns a HTTP status code in the 200 range the message will be acknowledged and removed from the queue, otherwise retried.

Note 2: Each argument reference is also set to `ForceNew`, since there is no support for updating the resource. Gives that each change in any of the argument, will force the provider to destroy and re-create the resource.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_webhook` can be imported using the resource identifier together with CloudAMQP instance identifier. The identifiers are CSV separated, see example below.

`terraform import cloudamqp_webhook.webhook_queue <id>,<instance_id>`

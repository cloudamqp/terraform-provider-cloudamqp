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

<details>
 <summary>
    <b>
      <i>Enable webhook before v1.30.0</i>
    </b>
  </summary>

Doesn't support updating the resource, makes all argument using `ForceNew` behaviour.
Any changes to an argument will destroy and re-create the resource. Also no longer support
for using `retry_interval` argument in the backend, even though it's set to be required.

```hcl
resource "cloudamqp_webhook" "webhook_queue" {
  instance_id = cloudamqp_instance.instance.id
  vhost = cloudamqp_instance.instance.vhost
  queue = "webhook-queue"
  webhook_uri = "https://example.com/webhook?key=secret"
  retry_interval = 5
  concurrency = 5
}
```

</details>

<details>
 <summary>
    <b>
      <i>Enable webhook after [v1.30.0](https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.30.0)</i>
    </b>
  </summary>

Support to updating the resource which makes the argument no longer require `ForceNew` behaviour.
The argument `retry_interval` have also been removed.

```hcl
resource "cloudamqp_webhook" "webhook_queue" {
  instance_id = cloudamqp_instance.instance.id
  vhost = cloudamqp_instance.instance.vhost
  queue = "webhook-queue"
  webhook_uri = "https://example.com/webhook?key=secret"
  concurrency = 5
}
```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id`     - (Required) The CloudAMQP instance ID.
* `vhost`           - (Required) The vhost the queue resides in.
* `queue`           - (Required) A (durable) queue on your RabbitMQ instance.
* `webhook_uri`     - (Required) A POST request will be made for each message in the queue to this endpoint.
* `concurrency`     - (Required) Max simultaneous requests to the endpoint.

___

* `retry_interval`  - ~~(Required)~~ No longer supported in the backend and removed as an argument from v1.30.0. Still required to be set before v1.30.0.

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_webhook` can be imported using the resource identifier together with CloudAMQP instance identifier. The identifiers are CSV separated, see example below.

`terraform import cloudamqp_webhook.webhook_queue <id>,<instance_id>`

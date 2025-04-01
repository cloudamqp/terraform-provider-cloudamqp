---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_webhook"
description: |-
  Add, update or remove a webhook for a vhost and queue
---

# cloudamqp_webhook

This resource allows you to add, update or remove a swebhook for a specific vhost and queue.

Only available for dedicated subscription plans.

## Example Usage

<details>
 <summary>
    <b>
      <i>Enable webhook from </i>
      <a href="https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.30.0">v1.30.0</a>
    </b>
  </summary>

Support to updating the resource which makes the argument no longer require `ForceNew` behaviour.
The argument `retry_interval` have also been removed.

```hcl
resource "cloudamqp_webhook" "webhook_queue" {
  instance_id = cloudamqp_instance.instance.id
  vhost       = cloudamqp_instance.instance.vhost
  queue       = "webhook-queue"
  webhook_uri = "https://example.com/webhook?key=secret"
  concurrency = 5
}
```

</details>

<details>
 <summary>
    <b>
      <i>Enable webhook before v1.30.0</i>
    </b>
  </summary>

For more information see below [versions](#versions) section.

```hcl
resource "cloudamqp_webhook" "webhook_queue" {
  instance_id     = cloudamqp_instance.instance.id
  vhost           = cloudamqp_instance.instance.vhost
  queue           = "webhook-queue"
  webhook_uri     = "https://example.com/webhook?key=secret"
  retry_interval  = 5
  concurrency     = 5
}
```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance ID.
* `vhost`       - (Required) The vhost the queue resides in.
* `queue`       - (Required) A (durable) queue on your RabbitMQ instance.
* `webhook_uri` - (Required) A POST request will be made for each message in the queue to this
                  endpoint.
* `concurrency` - (Required) Max simultaneous requests to the endpoint.

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_webhook` can be imported using the resource identifier together with CloudAMQP instance
identifier (CSV separated). To retrieve the resource identifier, use [CloudAMQP API list webhooks].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_webhook.webhook_queue
  id = format("<id>,%s", cloudamqp_instance.instance.id)
}
```

Or use Terraform CLI:

`terraform import cloudamqp_webhook.webhook_queue <id>,<instance_id>`

## Versions

Information for older versions

<details>
  <summary>
    <i>Before v1.30.0</i>
  </summary>

  Versions before v1.30.0 doesn't support updating the resource, therefore all arguments using the
  `ForceNew` behaviour. Any changes to an argument will destroy and re-create the resource. The
  argument `retry_interval` is set to required, even if it's no longer supported in the backend.

  <b>Example Usage</b>
  
  ```hcl
    resource "cloudamqp_webhook" "webhook_queue" {
    instance_id     = cloudamqp_instance.instance.id
    vhost           = cloudamqp_instance.instance.vhost
    queue           = "webhook-queue"
    webhook_uri     = "https://example.com/webhook?key=secret"
    retry_interval  = 5
    concurrency     = 5
  }
  ```

  **Argument Reference**

  The following arguments are supported:

  > * `instance_id`     - (Required/ForceNew) The CloudAMQP instance ID.
  > * `vhost`           - (Required/ForceNew) The vhost the queue resides in.
  > * `queue`           - (Required/ForceNew) A (durable) queue on your RabbitMQ instance.
  > * `webhook_uri`     - (Required/ForceNew) A POST request will be made for each message in the
                          queue to this endpoint.
  > * `retry_interval`  - (Required/ForceNew) How often we retry if your endpoint fails (in seconds).
  > * `concurrency`     - (Required/ForceNew) Max simultaneous requests to the endpoint.

</details>

[CloudAMQP API list webhooks]: https://docs.cloudamqp.com/cloudamqp_api.html#list-webhooks

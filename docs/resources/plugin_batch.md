---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_plugin_batch"
description: |-
  Enable and disable multiple RabbitMQ plugins in a single batch operation.
---

<!-- markdownlint-disable MD033 -->

# cloudamqp_plugin_batch

This resource allows you to enable or disable multiple RabbitMQ plugins in a single batch
operation. Compared to [cloudamqp_plugin], this resource manages a set of plugins together and
minimises the number of API calls by sending only what has changed on update.

Only available for dedicated subscription plans running ***RabbitMQ***.

## Example Usage

```hcl
resource "cloudamqp_instance" "instance" {
  name   = "terraform-cloudamqp-instance"
  plan   = "bunny-1"
  region = "amazon-web-services::us-east-1"
  tags   = ["terraform"]
}

resource "cloudamqp_plugin_batch" "plugins" {
  instance_id = cloudamqp_instance.instance.id

  plugins = {
    rabbitmq_stomp           = true
    rabbitmq_top             = true
    rabbitmq_web_mqtt        = true
    rabbitmq_random_exchange = false
  }
}
```

<details>
  <summary>
    <b>
      <i>Faster instance destroy when running `terraform destroy`</i>
    </b>
  </summary>

Set `enable_faster_instance_destroy` to ***true*** in the provider configuration to skip disabling
plugins when destroying the instance.

```hcl
provider "cloudamqp" {
  apikey                         = var.cloudamqp_customer_api_key
  enable_faster_instance_destroy = true
}

resource "cloudamqp_instance" "instance" {
  name   = "terraform-cloudamqp-instance"
  plan   = "bunny-1"
  region = "amazon-web-services::us-east-1"
  tags   = ["terraform"]
}

resource "cloudamqp_plugin_batch" "plugins" {
  instance_id = cloudamqp_instance.instance.id

  plugins = {
    rabbitmq_stomp    = true
    rabbitmq_top      = true
    rabbitmq_web_mqtt = true
  }
}
```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance ID.
* `plugins`     - (Required) A map of plugin name to enabled state. Set a plugin to `true` to
                  enable it, or `false` to disable it. Only the plugins listed in this map are
                  managed by this resource; all other plugins on the instance are left untouched.
* `sleep`       - (Optional) Configurable sleep time (seconds) used when polling for job completion.
                  Default set to 10 seconds.
* `timeout`     - (Optional) Configurable timeout time (seconds) for the operation. Default set to
                  1800 seconds.

## Attributes Reference

All attributes reference are computed

* `id` - The identifier for this resource.

## Behaviour

### Create

Sends an enable request containing all plugins in the map with value `true`.

### Update

Computes the minimal diff between the previous and desired state:

* Plugins changed to `true` (or newly added as `true`) are sent in the **enable** list.
* Plugins changed to `false` (or removed from the map, treated as disabled) are sent in the
  **disable** list.

Only the changed plugins are sent to the API, keeping the request as small as possible.

### Delete

Sends a disable request for all plugins in the map with value `true`. Plugins already set to
`false` are ignored.

## Dependency

This resource depends on the CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

Not possible to import this resource.

## Enable faster instance destroy

When running `terraform destroy` this resource will try to disable the managed plugins before
deleting `cloudamqp_instance`. This is not necessary since the servers will be deleted.

Set `enable_faster_instance_destroy` to ***true*** in the provider configuration to skip this.

## Required plugins

The following plugins are always enabled by CloudAMQP and do not need to be managed:

| Name                      | Version |
|---------------------------|---------|
| rabbitmq_management       | all     |
| rabbitmq_management_agent | all     |
| rabbitmq_prometheus       | 3.10.0  |
| rabbitmq_web_dispatch     | all     |

[cloudamqp_plugin]: plugin.md

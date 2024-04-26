---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_plugin"
description: |-
  Enable and disable Rabbit MQ plugin.
---

# cloudamqp_plugin

This resource allows you to enable or disable Rabbit MQ plugins.

Only available for dedicated subscription plans running ***RabbitMQ***.

## Example Usage

```hcl
resource "cloudamqp_plugin" "rabbitmq_top" {
  instance_id = cloudamqp_instance.instance.id
  name        = "rabbitmq_top"
  enabled     = true
}
```

<details>
  <summary>
    <b>
      <i>Enable multiple plugins v1.19.1 and older versions</i>
    </b>
  </summary>

Rabbit MQ can only change one plugin at a time. It will fail if multiple plugins resources are used, unless by creating dependencies with `depend_on` between the resources. Once one plugin has been enabled, the other will continue. See example below.

```hcl
resource "cloudamqp_plugin" "rabbitmq_top" {
  instance_id = cloudamqp_instance.instance.id
  name        = "rabbitmq_top"
  enabled     = true
}

resource "cloudamqp_plugin" "rabbitmq_amqp1_0" {
  instance_id = cloudamqp_instance.instance.id
  name        = "rabbitmq_amqp1_0"
  enabled     = true

  depends_on = [
    cloudamqp_plugin.rabbitmq_top
  ]
}
```

</details>

<details>
  <summary>
    <b>
      <i>Enable multiple plugins from </i>
      <a href="https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.19.2">v1.19.2</a>
    </b>
  </summary>

CloudAMQP Terraform provider [v1.19.2](https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.19.2) support asynchronous request for plugin actions.

```hcl
resource "cloudamqp_plugin" "rabbitmq_top" {
  instance_id = cloudamqp_instance.instance.id
  name        = "rabbitmq_top"
  enabled     = true
}

resource "cloudamqp_plugin" "rabbitmq_amqp1_0" {
  instance_id = cloudamqp_instance.instance.id
  name        = "rabbitmq_amqp1_0"
  enabled     = true
}
```

</details>

<details>
  <summary>
    <b>
      <i>Faster instance destroy when running `terraform destroy` from </i>
      <a href="https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.27.0">v1.27.0</a>
    </b>
  </summary>

CloudAMQP Terraform provider [v1.27.0](https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.27.0) enables faster `cloudamqp_instance` destroy when running `terraform destroy`.

```hcl
# Configure the CloudAMQP Provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
  enable_faster_instance_destroy = true
}

resource "cloudamqp_instance" "instance" {
  name    = "terraform-cloudamqp-instance"
  plan    = "bunny-1"
  region  = "amazon-web-services::us-west-1"
  tags    = ["terraform"]
}

resource "cloudamqp_plugin" "rabbitmq_top" {
  instance_id = cloudamqp_instance.instance.id
  name        = "rabbitmq_top"
  enabled     = true
}

resource "cloudamqp_plugin" "rabbitmq_amqp1_0" {
  instance_id = cloudamqp_instance.instance.id
  name        = "rabbitmq_amqp1_0"
  enabled     = true
}
```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance ID.
* `name`        - (Required) The name of the Rabbit MQ plugin.
* `enabled`     - (Required) Enable or disable the plugins.
* `sleep` - (Optional) Configurable sleep time (seconds) for retries when requesting information
about plugins. Default set to 10 seconds. *Available from v1.29.0*
* `timeout` - (Optional) - Configurable timeout time (seconds) for retries when requesting
information about plugins. Default set to 1800 seconds. *Available from v1.29.0*

## Attributes Reference

All attributes reference are computed

* `id`          - The identifier for this resource.
* `description` - The description of the plugin.
* `version`     - The version of the plugin.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

If multiple plugins should be enable, create dependencies between the plugin resources. See example above.

## Import

`cloudamqp_plugin` can be imported using the name argument of the resource together with CloudAMQP instance identifier. The name and identifier are CSV separated, see example below.

`terraform import cloudamqp_plugin.rabbitmq_management rabbitmq_management,<instance_id>`

## Required plugins

Plugins that is not needed to be managed by the provider since they will always be enabled. We have made `rabbitmq_prometheus` required for all clusters since RabbitMQ version 3.10.0.

| Name                      | Version |
|---------------------------|---------|
| rabbitmq_management       | all     |
| rabbitmq_management_agent | all     |
| rabbitmq_prometheus       | 3.10.0  |

## Enable faster instance destroy

When running `terraform destroy` this resource will try to disable the managed plugin before
deleting `cloudamqp_instance`. This is not necessary since the servers will be deleted.

Set `enable_faster_instance_destroy` to ***true*** in the provider configuration to skip this.

---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_plugin_commiunity"
description: |-
  Install or uninstall community plugin.
---

# cloudamqp_plugin_community

This resource allows you to install or uninstall community plugins. Once installed the plugin will be available in `cloudamqp_plugin`.

Only available for dedicated subscription plans running ***RabbitMQ***.

~> CloudAMQP Terraform provider [v1.11.0](https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.11.0) there is support for multiple retries when requesting information about community plugins. This was introduced to avoid `ReadPluginCommunity error 400: Timeout talking to backend`.

~> CloudAMQP Terraform provider [v1.19.2](https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.19.2) support asynchronous request for plugin/community actions. Solve issues reported when enable multiple plugins.

## Example Usage

```hcl
resource "cloudamqp_plugin_community" "rabbitmq_delayed_message_exchange" {
  instance_id = cloudamqp_instance.instance.id
  name        = "rabbitmq_delayed_message_exchange"
  enabled     = true
}
```

<details>
  <summary>
    <b>
      <i>Faster instance destroy when running `terraform destroy` from v1.27.0
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

resource "cloudamqp_plugin_community" "rabbitmq_delayed_message_exchange" {
  instance_id = cloudamqp_instance.instance.id
  name        = "rabbitmq_delayed_message_exchange"
  enabled     = true
}
```
</details>

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance ID.
* `name`        - (Required) The name of the Rabbit MQ community plugin.
* `enabled`     - (Required) Enable or disable the plugins.

## Attributes Reference

* `id`          - The identifier for this resource.
* `description` - The description of the plugin.
* `require`     - Required version of RabbitMQ.

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_plugin` can be imported using the name argument of the resource together with CloudAMQP instance identifier. The name and identifier are CSV separated, see example below.

`terraform import cloudamqp_plugin.<resource_name> <plugin_name>,<instance_id>`

## Enable faster instance destroy

When running `terraform destroy` this resource will try to uninstall the managed community plugin before deleting `cloudamqp_instance`. This is not necessary since the servers will be deleted.

Set `enable_faster_instance_destroy` to ***true***  in the provider configuration to skip this.

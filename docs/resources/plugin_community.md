---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_plugin_commiunity"
description: |-
  Install or uninstall community plugin.
---

# cloudamqp_plugin_community

This resource allows you to install or uninstall community plugins. Once installed the plugin will
be available in `cloudamqp_plugin`.

Only available for dedicated subscription plans running ***RabbitMQ***.

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
      <i>Faster instance destroy when running `terraform destroy` from [v1.27.0]</i>
    </b>
  </summary>

CloudAMQP Terraform provider [v1.27.0] enables faster `cloudamqp_instance` destroy when running
`terraform destroy`.

```hcl
# Configure the CloudAMQP Provider
provider "cloudamqp" {
  apikey                          = var.cloudamqp_customer_api_key
  enable_faster_instance_destroy  = true
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
* `sleep`       - (Optional) Configurable sleep time (seconds) for retries when requesting
                  information about community plugins. Default set to 10 seconds.

  ***Note:*** Available from [v1.29.0]

* `timeout`     - (Optional) - Configurable timeout time (seconds) for retries when requesting
                  information about community plugins. Default set to 1800 seconds.

  ***Note:*** Available from [v1.29.0]

## Attributes Reference

* `id`          - The identifier for this resource.
* `description` - The description of the plugin.
* `require`     - Required version of RabbitMQ.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_plugin_community` can be imported if it's has already been installed by using the name
argument of the resource together with CloudAMQP instance identifier (CSV separated). To retrieve
list of available community plugins, use [CloudAMQP API list community plugins].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_plugin_community.rabbitmq_delayed_message_exchange
  id = format("rabbitmq_delayed_message_exchange,%s", cloudamqp_instance.instance.id)
}
```

Or use Terraform CLI:

`terraform import cloudamqp_plugin.rabbitmq_delayed_message_exchange <plugin_name>,<instance_id>`

## Enable faster instance destroy

When running `terraform destroy` this resource will try to uninstall the managed community plugin
before deleting `cloudamqp_instance`. This is not necessary since the servers will be deleted.

Set `enable_faster_instance_destroy` to ***true***  in the provider configuration to skip this.

[CloudAMQP API list community plugins]: https://docs.cloudamqp.com/instance-api.html#tag/plugins/get/plugins/community
[v1.27.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.27.0
[v1.29.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.29.0

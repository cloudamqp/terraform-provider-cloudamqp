---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_instance"
description: |-
  Creates and manages a Rabbit MQ instance within CloudAMQP.
---

# cloudamqp_instance

This resource allows you to create and manage a CloudAMQP instance running either [**RabbitMQ**] or
[**LavinMQ**] and can be deployed to multiple cloud platforms provider and regions, see
[instance regions] for more information.

Once the instance is created it will be assigned a unique identifier. All other resources and data
sources created for this instance needs to reference this unique instance identifier.

Pricing is available at [CloudAMQP plans].

## Example Usage

<details>
  <summary>
    <b>
      <i>Shared and dedicated instances running LavinMQ</i>
    </b>
  </summary>

```hcl
# Minimum free lemming instance running LavinMQ
resource "cloudamqp_instance" "lemur_instance" {
  name    = "cloudamqp-free-instance-01"
  plan    = "lemming"
  region  = "amazon-web-services::us-west-1"
  tags    = ["lavinmq"]
}

# Minimum free lemur instance running RabbitMQ
resource "cloudamqp_instance" "lemming_instance" {
  name    = "cloudamqp-free-instance-02"
  plan    = "lemur"
  region  = "amazon-web-services::us-west-1"
  tags    = ["rabbitmq"]
}

# Dedicated penguin instance running LavinMQ
resource "cloudamqp_instance" "penguin_instance" {
  name    = "terraform-cloudamqp-instance-01"
  plan    = "penguin-1"
  region  = "amazon-web-services::us-west-1"
  tags    = ["lavinmq"]
}

# Dedicated bunny instance running RabbitMQ
resource "cloudamqp_instance" "bunny_instance" {
  name    = "terraform-cloudamqp-instance-02"
  plan    = "bunny-1"
  region  = "amazon-web-services::us-west-1"
  tags    = ["rabbitmq"]
}
```

</details>

<details>
  <summary>
    <b>
      <i>Dedicated instance using attribute vpc_subnet to create VPC, before v1.16.0</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_instance" "instance" {
  name        = "terraform-cloudamqp-instance"
  plan        = "penguin-1"
  region      = "amazon-web-services::us-west-1"
  tags        = ["terraform"]
  vpc_subnet  = "10.56.72.0/24"
}
```

</details>

<details>
  <summary>
    <b>
      <i>
        Dedicated instance using attribute vpc_subnet to create VPC and then import managed VPC,
        from v1.16.0 (Managed VPC)
      </i>
    </b>
  </summary>

```hcl
# Dedicated instance that also creates VPC
resource "cloudamqp_instance" "instance_01" {
  name        = "terraform-cloudamqp-instance-01"
  plan        = "penguin-1"
  region      = "amazon-web-services::us-west-1"
  tags        = ["terraform"]
  vpc_subnet  = "10.56.72.0/24"
}
```

Once the instance and the VPC are created, the VPC can be imported as managed VPC and added to the
configuration file. Set attribute `vpc_id` to the managed VPC identifier. To keep the managed VPC
when deleting the instance, set attribute `keep_associated_vpc` to true. For more information see
guide [Managed VPC].

```hcl
# Imported managed VPC
resource "cloudamqp_vpc" "vpc" {
  name    = "<vpc-name>"
  region  = "amazon-web-services::us-east-1"
  subnet  = "10.56.72.0/24"
  tags    = []
}

# Add vpc_id and keep_associated_vpc attributes
resource "cloudamqp_instance" "instance_01" {
  name                = "terraform-cloudamqp-instance-01"
  plan                = "penguin-1"
  region              = "amazon-web-services::us-west-1"
  tags                = ["terraform"]
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}
```

</details>

<details>
  <summary>
    <b>
      <i>Dedicated instances and managed VPC, from v1.16.0 (Managed VPC)</i>
    </b>
  </summary>

```hcl
# Managed VPC
resource "cloudamqp_vpc" "vpc" {
  name    = "<vpc-name>"
  region  = "amazon-web-services::us-east-1"
  subnet  = "10.56.72.0/24"
  tags    = []
}

# First instance added to managed VPC
resource "cloudamqp_instance" "instance_01" {
  name                = "terraform-cloudamqp-instance-01"
  plan                = "penguin-1"
  region              = "amazon-web-services::us-west-1"
  tags                = ["terraform"]
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

# Second instance added to managed VPC
resource "cloudamqp_instance" "instance_02" {
  name                = "terraform-cloudamqp-instance-02"
  plan                = "penguin-1"
  region              = "amazon-web-services::us-west-1"
  tags                = ["terraform"]
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}
```

Set attribute `keep_associated_vpc` to true, will keep managed VPC when deleting the instances.

</details>

## Argument Reference

The following arguments are supported:

* `name`    - (Required) Name of the CloudAMQP instance.
* `plan`    - (Required) The subscription plan. See available [plans].
* `region`  - (Required) The region to host the instance in. See available [regions].

  ***Note:*** Changing region will force the instance to be destroyed and a new created in the new
              region. All data will be lost and a new name assigned.

* `nodes`   - (Optional/Computed) Number of nodes, 1, 3 or 5 depending on plan used. Only needed for
              legacy plans, will otherwise be computed.

  ***Deprecated:*** Legacy subscriptions plan can still change this to scale up or down the instance.
                    New subscriptions plans use the plan to determine number of nodes. In order to
                    change number of nodes the `plan` needs to be updated.

* `tags`        - (Optional) One or more tags for the CloudAMQP instance, makes it possible to
                  categories multiple instances in console view. Default there is no tags assigned.
* `rmq_version` - (Optional/Computed) The Rabbit MQ version. Can be left out, will then be set to
                  default value used by CloudAMQP API.

  ***Note:*** There is not yet any support in the provider to change the RMQ version. Once it's set
              in the initial creation, it will remain.

* `vpc_id`      - (Optional/Computed) The VPC ID. Use this to create your instance in an existing
                  VPC. See available [example].
* `vpc_subnet`  - (Optional/Computed) Creates a dedicated VPC subnet, shouldn't overlap with other
                  VPC subnet, default subnet used 10.56.72.0/24.

  ***Deprecated:*** Will be removed in next major version (v2.0)

  ***Note:*** Extra fee will be charged when using VPC, see [CloudAMQP] for more information.

* `no_default_alarm`    - (Optional/Computed) Set to true to discard creating default alarms when
                          the instance is created. Can be left out, will then use default
                          value = false.
* `keep_associated_vpc` - (Optional) Keep associated VPC when deleting instance. Default set to
                          false.
* `copy_settings`       - (Optional) Copy settings from one CloudAMQP instance to a new. Consists of
                          the block documented below.

___

The `copy_settings` block consists of:

* `subscription_id` - (Required) Instance identifier of the CloudAMQP instance to copy the settings
                      from.
* `settings`        - (Required) Array of one or more settings to be copied. Allowed values:
                      [alarms, config, definitions, firewall, logs, metrics, plugins]

See more below, [copy settings].

## Attributes Reference

All attributes reference are computed

* `id`            - The identifier (instance_id) for this resource, used as a reference by almost
                    all other resource and data sources
* `url`           - The AMQP URL (uses the internal hostname if the instance was created with VPC).
                    Has the format: `amqps://{username}:{password}@{hostname}/{vhost}`
* `apikey`        - API key needed to communicate to CloudAMQP's second API. The second API is used
                    to manage alarms, integration and more, full description [CloudAMQP API].
* `host`          - The external hostname for the CloudAMQP instance.
* `host_internal` - The internal hostname for the CloudAMQP instance.
* `vhost`         - The virtual host used by Rabbit MQ.
* `dedicated`     - Information if the CloudAMQP instance is shared or dedicated.
* `backend`       - Information if the CloudAMQP instance runs either RabbitMQ or LavinMQ.

## Import

`cloudamqp_instance`can be imported using resource identifier. To retrieve the resource identifier,
use [CloudAMQP API list instances]

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_instance.instance
  id = <id>
}
```

Or use Terraform CLI:

`terraform import cloudamqp_instance.instance <id>`

## Upgrade and downgrade

It's possible to upgrade or downgrade your subscription plan, this will either increase or decrease
the underlying resource used for by the CloudAMQP instance. To do this, change the argument `plan`
in the configuration and apply the changes. See available [plans].

<details>
  <summary>
    <b>
      <i>Upgrade the subscription plan</i>
    </b>
  </summary>

```hcl
# Initial CloudAMQP instance configuration
resource "cloudamqp_instance" "instance" {
  name    = "instance"
  plan    = "puffin-1"
  region  = "amazon-web-services::us-west-1"
  tags    = ["terraform"]
}

# Upgraded CloudAMQP instance configuration
resource "cloudamqp_instance" "instance" {
  name    = "instance"
  plan    = "penguin-1"
  region  = "amazon-web-services::us-west-1"
  tags    = ["terraform"]
}
```

</details>

<details>
  <summary>
    <b>
      <i>Downgrade number of nodes from 3 to 1</i>
    </b>
  </summary>

```hcl
# Initial CloudAMQP instance configuration
resource "cloudamqp_instance" "instance" {
  name    = "instance"
  plan    = "penguin-3"
  region  = "amazon-web-services::us-west-1"
  tags    = ["terraform"]
}

# Downgraded CloudAMQP instance configuration
resource "cloudamqp_instance" "instance" {
  name    = "instance"
  plan    = "penguin-1"
  region  = "amazon-web-services::us-west-1"
  tags    = ["terraform"]
}
```

</details>

## Copy settings to a new dedicated instance

With copy settings it's possible to create a new dedicated instance with settings such as alarms,
config, etc. from another dedicated instance. This can be done by adding the `copy_settings` block
to this resource and populate `subscription_id` with a CloudAMQP instance identifier from another
already existing instance.

~> `rmq_version` argument is required when doing this action. Must match the broker version of the
dedicated instance to be copied from.

Then add the settings to be copied over to the new dedicated instance. Settings that can be copied

### Settings supported by LavinMQ

***Allowed values:*** alarms, definitions, firewall, metrics

<details>
  <summary>
    <b>
      <i>Copy settings from a dedicated instance to a new dedicated instance</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_instance" "penguin_instance" {
  name        = "terraform-cloudamqp-instance-01"
  plan        = "penguin-1"
  region      = "amazon-web-services::us-west-1"
  rmq_version = "2.2.0"
  tags        = ["terraform"]

  copy_settings {
    subscription_id = var.instance_id
    settings = ["alarms", "definitions", "firewall", "metrics"]
  }
}
```

</details>

### Settings supported by RabbitMQ

***Allowed values:*** alarms, config, definitions, firewall, logs, metrics, plugins

<details>
  <summary>
    <b>
      <i>Copy settings from a dedicated instance to a new dedicated instance</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_instance" "bunny_instance" {
  name        = "terraform-cloudamqp-instance-02"
  plan        = "bunny-1"
  region      = "amazon-web-services::us-west-1"
  rmq_version = "3.12.2"
  tags        = ["terraform"]

  copy_settings {
    subscription_id = var.instance_id
    settings = ["alarms", "config", "definitions", "firewall", "logs", "metrics", "plugins"]
  }
}
```

</details>

[CloudAMQP]: https://cloudamqp.com
[CloudAMQP API]: https://docs.cloudamqp.com/cloudamqp_api.html
[CloudAMQP API list instances]: https://docs.cloudamqp.com/#list-instances
[CloudAMQP plans]: https://www.cloudamqp.com/plans.html
[copy settings]: #copy-settings-to-a-new-dedicated-instance
[example]: ../guides/info_vpc_existing.md
[regions]: ../guides/info_region.md
[**LavinMQ**]: https://lavinmq.com
[Managed VPC]: ../guides/info_managed_vpc#dedicated-instance-and-vpc_subnet
[plans]: ../guides/info_plan.md
[**RabbitMQ**]: https://www.rabbitmq.com

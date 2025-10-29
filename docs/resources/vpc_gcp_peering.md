---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_vpc_gcp_peering"
description: |-
  Create VPC peering configuration to another VPC network hosted in GCP
---

# cloudamqp_vpc_gcp_peering

This resouce creates a VPC peering configuration for the CloudAMQP instance. The configuration will
connect to another VPC network hosted on Google Cloud Platform (GCP). See the [GCP documentation]
for more information on how to create the VPC peering configuration.

~> **Note:** Creating a VPC peering will automatically add firewall rules for the peered subnet.

<details>
 <summary>
    <i>Default VPC peering firewall rule</i>
  </summary>

For LavinMQ:

```hcl
rules {
  Description = "VPC peer request"
  ip          = "<VPC peered subnet>"
  ports       = [15672, 5552, 5551]
  services    = ["AMQP", "AMQPS", "HTTPS"]
}
```

For RabbitMQ:

```hcl
rules {
  Description = "VPC peer request"
  ip          = "<VPC peered subnet>"
  ports       = [15672]
  services    = ["AMQP", "AMQPS", "HTTPS", "STREAM", "STREAM_SSL"]
}
```

</details>

Pricing is available at [CloudAMQP plans].

Only available for dedicated subscription plans.

## Example Usage

<details>
  <summary>
    <b>
      <i>VPC peering before v1.16.0</i>
    </b>
  </summary>

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# CloudAMQP instance
resource "cloudamqp_instance" "instance" {
  name        = "terraform-vpc-peering"
  plan        = "penguin-1"
  region      = "google-compute-engine::europe-north1"
  tags        = ["terraform"]
  vpc_subnet  = "10.40.72.0/24"
}

# VPC information
data "cloudamqp_vpc_gcp_info" "vpc_info" {
  instance_id = cloudamqp_instance.instance.id
}

# VPC peering configuration
resource "cloudamqp_vpc_gcp_peering" "vpc_peering_request" {
  instance_id       = cloudamqp_instance.instance.id
  peer_network_uri  = "https://www.googleapis.com/compute/v1/projects/PROJECT-NAME/global/networks/VPC-NETWORK-NAME"
}
```

</details>

<details>
  <summary>
    <b>
      <i>VPC peering from [v1.16.0] (Managed VPC)</i>
    </b>
  </summary>

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# Managed VPC resource
resource "cloudamqp_vpc" "vpc" {
  name    = "<VPC name>"
  region  = "google-compute-engine::europe-north1"
  subnet  = "10.56.72.0/24"
  tags    = []
}

# CloudAMQP instance
resource "cloudamqp_instance" "instance" {
  name   = "terraform-vpc-peering"
  plan   = "penguin-1"
  region = "google-compute-engine::europe-north1"
  tags   = ["terraform"]
  vpc_id = cloudamqp_vpc.vpc.id
}

# VPC information
data "cloudamqp_vpc_gcp_info" "vpc_info" {
  vpc_id = cloudamqp_vpc.vpc.info
  # or
  # instance_id = cloudamqp_instance.instance.id
}

# VPC peering configuration
resource "cloudamqp_vpc_gcp_peering" "vpc_peering_request" {
  vpc_id = cloudamqp_vpc.vpc.id
  # or
  # instance_id = cloudamqp_instance.instance.id
  peer_network_uri = "https://www.googleapis.com/compute/v1/projects/PROJECT-NAME/global/networks/VPC-NETWORK-NAME"
}
```

</details>

<details>
  <summary>
    <b>
      <i>VPC peering from [v1.28.0], wait_on_peering_status </i>
    </b>
  </summary>

Default peering request, no need to set `wait_on_peering_status`. It's default set to false and will
not wait on peering status. Create resource will be considered completed, regardless of the status
of the state.

```hcl
resource "cloudamqp_vpc_gcp_peering" "vpc_peering_request" {
  vpc_id            = cloudamqp_vpc.vpc.id
  peer_network_uri  = "https://www.googleapis.com/compute/v1/projects/ROJECT-NAME/global/networks/VPC-NETWORK-NAME"
}
```

Peering request and waiting for peering status of the state to change to ACTIVE before the create
resource is consider complete. This is done once both side have done the peering.

```hcl
resource "cloudamqp_vpc_gcp_peering" "vpc_peering_request" {
  vpc_id                  = cloudamqp_vpc.vpc.id
  wait_on_peering_status  = true
  peer_network_uri        = "https://www.googleapis.com/compute/v1/projects/PROJECT-NAME/global/networks/VPC-NETWORK-NAME"
}
```

</details>

## Argument Reference

* `instance_id`             - (Optional) The CloudAMQP instance identifier.

  ***Deprecated:*** from [v1.16.0], will be removed in next major version (v2.0)

* `vpc_id`                  - (Optional) The managed VPC identifier.

  ***Note:*** Available from [v1.16.0], will be required in next major version (v2.0)

* `peer_network_uri`        - (Required) Network URI of the VPC network to which you will peer with.
                              See examples above for the format.
* `wait_on_peering_status`  - (Optional) Makes the resource wait until the peering is connected.
                              Default set to false.

  ***Note:*** Available from [v1.28.0]

* `sleep`                   - (Optional) Configurable sleep time (seconds) between retries when
                              requesting or reading peering. Default set to 10 seconds.

  ***Note:*** Available from [v1.29.0]

* `timeout`                 - (Optional) Configurable timeout time (seconds) before retries times
                              out. Default set to 1800 seconds.

  ***Note:*** Available from [v1.29.0]

## Attributes Reference

All attributes reference are computed

* `id`                  - The identifier for this resource.
* `state`               - VPC peering state
* `state_details`       - VPC peering state details
* `auto_create_routes`  - VPC peering auto created routes

## Dependency

***From v1.16.0:***
This resource depends on CloudAMQP managed VPC identifier, `cloudamqp_vpc.vpc.id` or instance
identifier, `cloudamqp_instance.instance.id`.

***Before v1.16.0:***
This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

***From v1.32.2:***
`cloudamqp_vpc_gcp_peering` can be imported while using the resource type, with CloudAMQP VPC
identifier or instance identifier together with *peering_network_uri* (CSV seperated).

### Resource type VPC

To use the CloudAMQP managed VPC identifier set the resource type to *vpc*.

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_vpc_gcp_peering.this
  id = "vpc,<vpc_id>,<peering_network_uri>"
}
```

Or use Terraform CLI:

```hcl
terraform import cloudamqp_vpc_gcp_peering.vpc_peering_request vpc,<vpc_id>,<peer_network_uri>
```

### Resource type instance

To use the Cloudamqp instance identifier set the resource type to *instance*

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_vpc_gcp_peering.this
  id = "instance,<instance_id>,<peering_network_uri>"
}
```

Or use Terraform CLI:

```hcl
terraform import cloudamqp_vpc_gcp_peering.vpc_peering_request instance,<instance_id>,<peer_network_uri>
```

***Before v1.32.2:***
Not possible to import this resource.

### Peering network URI

This is required to be able to import the correct peering. Following the same format as the argument
reference.

```hcl
https://www.googleapis.com/compute/v1/projects/PROJECT-NAME/global/networks/VPC-NETWORK-NAME
```

## Create VPC peering with additional firewall rules

To create a VPC peering configuration with additional firewall rules, it's required to chain the
[cloudamqp_security_firewall] resource to avoid parallel conflicting resource calls. This is done by
adding dependency from the firewall resource to the VPC peering resource.

Furthermore, since all firewall rules are overwritten, the otherwise automatically added rules for
the VPC peering also needs to be added.

See example below.

## Example usage with additional firewall rules

<details>
  <summary>
    <b>
      <i>VPC peering before v1.16.0</i>
    </b>
  </summary>

```hcl
# VPC peering configuration
resource "cloudamqp_vpc_gcp_peering" "vpc_peering_request" {
  instance_id       = cloudamqp_instance.instance.id
  peer_network_uri  = var.peer_network_uri
}

# Firewall rules
resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  # Default VPC peering rule
  rules {
    ip          =  var.peer_subnet
    ports       = [15672, 5552, 5551]
    services    = ["AMQP","AMQPS"]
    description = "VPC peering for <NETWORK>"
  }

  rules {
    ip        = "192.168.0.0/24"
    ports     = [4567, 4568]
    services  = ["AMQP","AMQPS", "HTTPS"]
  }

  depends_on = [
    cloudamqp_vpc_gcp_peering.vpc_peering_request
  ]
}
```

</details>

<details>
  <summary>
    <b>
      <i>VPC peering from [v1.16.0] (Managed VPC)</i>
    </b>
  </summary>

```hcl
# VPC peering configuration
resource "cloudamqp_vpc_gcp_peering" "vpc_peering_request" {
  vpc_id = cloudamqp_vpc.vpc.id
  # vpc_id prefered over instance_id
  # instance_id = cloudamqp_instance.instance.id
  peer_network_uri = var.peer_network_uri
}

# Firewall rules
resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  # Default VPC peering rule
  rules {
    ip          =  var.peer_subnet
    ports       = [15672, 5552, 5551]
    services    = ["AMQP","AMQPS"]
    description = "VPC peering for <NETWORK>"
  }

  rules {
    ip          = "0.0.0.0/0"
    ports       = []
    services    = ["HTTPS"]
    description = "MGMT interface"
  }

  depends_on = [
    cloudamqp_vpc_gcp_peering.vpc_peering_request
  ]
}
```

</details>

[CloudAMQP plans]: https://www.cloudamqp.com/plans.html
[cloudamqp_security_firewall]: https://registry.terraform.io/providers/cloudamqp/cloudamqp/latest/docs/resources/security_firewall
[GCP documentation]: https://cloud.google.com/vpc/docs/using-vpc-peering
[v1.16.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.16.0
[v1.28.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.28.0
[v1.29.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.29.0

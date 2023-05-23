---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_vpc_gcp_peering"
description: |-
  Create VPC peering configuration to another VPC network hosted in GCP
---

# cloudamqp_vpc_gcp_peering

This resouce creates a VPC peering configuration for the CloudAMQP instance. The configuration will connect to another VPC network hosted on Google Cloud Platform (GCP). See the [GCP documentation](https://cloud.google.com/vpc/docs/using-vpc-peering) for more information on how to create the VPC peering configuration.

~> **Note:** Creating a VPC peering will automatically add firewall rules for the peered subnet.
<details>
 <summary>
    <i>Default VPC peering firewall rule</i>
  </summary>
```hcl
rules {
  Description = "VPC peer request"
  ip          = "<VPC peered subnet>"
  ports       = []
  services    = ["AMQP", "AMQPS", "HTTPS", "STREAM", "STREAM_SSL", "STOMP", "STOMPS", "MQTT", "MQTTS"]
}
```
</details>

Pricing is available at [cloudamqp.com](https://www.cloudamqp.com/plans.html).

Only available for dedicated subscription plans.

## Example Usage

<details>
  <summary>
    <b>
      <i>VPC peering pre v1.16.0</i>
    </b>
  </summary>

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# CloudAMQP instance
resource "cloudamqp_instance" "instance" {
  name   = "terraform-vpc-peering"
  plan   = "bunny-1"
  region = "google-compute-engine::europe-north1"
  tags   = ["terraform"]
  vpc_subnet = "10.40.72.0/24"
}

# VPC information
data "cloudamqp_vpc_gcp_info" "vpc_info" {
  instance_id = cloudamqp_instance.instance.id
}

# VPC peering configuration
resource "cloudamqp_vpc_gcp_peering" "vpc_peering_request" {
  instance_id = cloudamqp_instance.instance.id
  peer_network_uri = "https://www.googleapis.com/compute/v1/projects/<PROJECT-NAME>/global/networks/<NETWORK-NAME>"
}
```
</details>

<details>
  <summary>
    <b>
      <i>VPC peering post v1.16.0 (Managed VPC)</i>
    </b>
  </summary>

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# Managed VPC resource
resource "cloudamqp_vpc" "vpc" {
  name = "<VPC name>"
  region = "google-compute-engine::europe-north1"
  subnet = "10.56.72.0/24"
  tags = []
}

# CloudAMQP instance
resource "cloudamqp_instance" "instance" {
  name   = "terraform-vpc-peering"
  plan   = "bunny-1"
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
  peer_network_uri = "https://www.googleapis.com/compute/v1/projects/<PROJECT-NAME>/global/networks/<NETWORK-NAME>"
}
```
</details>

## Argument Reference

 *Note: this resource require either `instance_id` or `vpc_id` from v1.16.0*

* `instance_id` - (Optional) The CloudAMQP instance identifier.

 ***Depreacted: Changed from required to optional in v1.16.0, will be removed in next major version (v2.0)***

* `vpc_id` - (Optional) The managed VPC identifier.

 ***Note: Added as optional in version v1.16.0, will be required in next major version (v2.0)***

* `peer_network_uri`- (Required) Network uri of the VPC network to which you will peer with.

## Attributes Reference

All attributes reference are computed

* `id` - The identifier for this resource.
* `state` - VPC peering state
* `state_details` - VPC peering state details
* `auto_create_routes` - VPC peering auto created routes

## Depedency

*Pre v1.16.0*
This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

*Post v1.16.0*
This resource depends on CloudAMQP managed VPC identifier, `cloudamqp_vpc.vpc.id` or instance identifier, `cloudamqp_instance.instance.id`.

## Import

Not possible to import this resource.

## Create VPC Peering with additional firewall rules

To create a VPC peering configuration with additional firewall rules, it's required to chain the [cloudamqp_security_firewall](https://registry.terraform.io/providers/cloudamqp/cloudamqp/latest/docs/resources/security_firewall)
resource to avoid parallel conflicting resource calls. This is done by adding dependency from the firewall resource to the VPC peering resource.

Furthermore, since all firewall rules are overwritten, the otherwise automatically added rules for the VPC peering also needs to be added.

See example below.

## Example Usage with additional firewall rules

<details>
  <summary>
    <b>
      <i>VPC peering pre v1.16.0</i>
    </b>
  </summary>

```hcl
# VPC peering configuration
resource "cloudamqp_vpc_gcp_peering" "vpc_peering_request" {
  instance_id = cloudamqp_instance.instance.id
  peer_network_uri = var.peer_network_uri
}

# Firewall rules
resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  # Default VPC peering rule
  rules {
    ip          =  var.peer_subnet
    ports       = [15672]
    services    = ["AMQP","AMQPS", "STREAM", "STREAM_SSL"]
    description = "VPC peering for <NETWORK>"
  }

  rules {
    ip          = "192.168.0.0/24"
    ports       = [4567, 4568]
    services    = ["AMQP","AMQPS", "HTTPS"]
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
      <i>VPC peering post v1.16.0 (Managed VPC)</i>
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
    ports       = [15672]
    services    = ["AMQP","AMQPS", "STREAM", "STREAM_SSL"]
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

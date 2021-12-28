---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_vpc_gcp_peering"
description: |-
  Create VPC peering configuration to another VPC network hosted in GCP
---

# cloudamqp_vpc_gcp_peering

This resouce creates a VPC peering configuration for the CloudAMQP instance. The configuration will connect to another VPC network hosted on Google Cloud Platform (GCP). See the [GCP documentation](https://cloud.google.com/vpc/docs/using-vpc-peering) for more information on how to create the VPC peering configuration. 

Only available for dedicated subscription plans.

## Example Usage

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
  nodes  = 1
  tags   = ["terraform"]
  rmq_version = "3.8.4"
  vpc_subnet = "10.40.72.0/24"
}

# VPC information
data "cloudamqp_vpc_gcp info" "vpc_info" {
  instance_id = cloudamqp_instance.instance.id
}

# VPC peering configuration
resource "cloudamqp_vpc_gcp_peering" "vpc_peering_request" {
  instance_id = cloudamqp_instance.instance.id
  peer_network_uri = "https://www.googleapis.com/compute/v1/projects/<PROJECT-NAME>/global/networks/<NETWORK-NAME>"
}
```

**Note: Creating a VPC peering configuration will trigger a firewall change to automatically add rules for the peered subnet.**

```
rules {
  ip          = "<PEERED-NETWORK-SUBNET>"
  ports       = [15672]
  services    = ["AMQP","AMQPS", "STREAM", "STREAM_SSL"]
  description = "VPC peering for <NETWORK-NAME>"
  }
```

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance ID.
* `peer_network_uri`- (Required) Network uri of the VPC network to which you will peer with.

## Attributes Reference

All attributes reference are computed

* `id` - The identifier for this resource.
* `state` - VPC peering state
* `state_details` - VPC peering state details
* `auto_create_routes` - VPC peering auto created routes

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Create VPC Peering with additional firewall rules

To create a VPC peering configuration with additional firewall rules, it's required to chain the [cloudamqp_security_firewall](https://registry.terraform.io/providers/cloudamqp/cloudamqp/latest/docs/resources/security_firewall)
resource to avoid parallel conflicting resource calls. This is done by adding dependency from the firewall resource to the VPC peering resource.

Furthermore, since all firewall rules are overwritten, the otherwise automatically added rules for the VPC peering also needs to be added.

See example below.

## Example Usage with additional firewall rules

```hcl
# VPC peering configuration
resource "cloudamqp_vpc_gcp_peering" "vpc_peering_request" {
  instance_id = cloudamqp_instance.instance.id
  peer_network_uri = var.peer_network_uri
}

# Firewall rules
resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

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

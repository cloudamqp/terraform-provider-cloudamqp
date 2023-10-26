terraform {
  required_providers {
    cloudamqp = {
      source = "localhost/cloudamqp/cloudamqp"
      version = "~> 1.0"
    }
  }
}

locals {
  cloudamqp_customer_api_key = "<apikey>"
  gcp_network_uri="https://www.googleapis.com/compute/v1/projects/<PROJECT-ID>/global/networks/<VPC-NAME>"
  gcp_subnet="10.56.73.0/24"
}

provider "cloudamqp" {
  apikey  = local.cloudamqp_customer_api_key
}

resource "cloudamqp_vpc" "vpc" {
  name = "terraform-vpc-peering"
  region = "google-compute-engine::europe-north1"
  subnet = "10.56.72.0/24"
}

resource "cloudamqp_instance" "instance" {
 name   = "terraform-vpc-peering"
 plan   = "bunny-1"
 region = "google-compute-engine::europe-north1"
 vpc_id = cloudamqp_vpc.vpc.id
 keep_associated_vpc = true
}

data "cloudamqp_vpc_gcp_info" "vpc_info" {
  vpc_id = cloudamqp_vpc.vpc.id
}

output "network_uri" {
  value = data.cloudamqp_vpc_gcp_info.vpc_info.network
}

resource "cloudamqp_vpc_gcp_peering" "vpc_peering_request" {
  vpc_id = cloudamqp_vpc.vpc.id
  peer_network_uri = var.gcp_network_uri
  // Can be used to wait until peering is connected
  wait_on_peering_status = false
}

resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  rules {
    ip          = var.gcp_subnet
    ports       = [15672]
    services    = ["AMQPS", "HTTPS", "MQTTS"]
    description = "VPC peering for <NETWORK>"
  }

  depends_on = [
    cloudamqp_vpc_gcp_peering.vpc_peering_request
  ]
}
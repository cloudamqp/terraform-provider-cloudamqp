terraform {
  required_providers {
    cloudamqp = {
      source  = "cloudamqp/cloudamqp"
      version = "~>1.0"
    }
  }
}

provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

resource "cloudamqp_instance" "instance" {
  name   = "aws-instance"
  plan   = "bunny-1"
  region = "amazon-web-services::us-east-1"
  tags   = []
}

resource "cloudamqp_custom_certificate" "cert" {
  instance_id = cloudamqp_instance.instance.id
  ca          = file(var.custom_ca_path)
  cert        = file(var.custom_cert_path)
  private_key = file(var.custom_private_key_path)
  sni_hosts   = "my.custom.domain"
}

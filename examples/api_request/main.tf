terraform {
  required_providers {
    cloudamqp = {
      source = "cloudamqp/cloudamqp"
      version = "~> 1.0"
    }
    # Unofficial Terraform provider to access the RabbitMQ management HTTP API.
    # More information can be found at: 
    # https://registry.terraform.io/providers/cyrilgdn/rabbitmq/latest
    rabbitmq = {
      source = "cyrilgdn/rabbitmq"
      version = "~> 1.0"
    }
  }
}

provider "cloudamqp" {
  apikey  = var.cloudamqp_customer_api_key
}

resource "cloudamqp_instance" "instance" {
 name   = "instance"
 plan   = "bunny-1"
 region = "amazon-web-services::us-east-1"
}

data "cloudamqp_credentials" "local-admin" {
  instance_id = cloudamqp_instance.instance.id
}

provider "rabbitmq" {
  endpoint = format("https://%s", cloudamqp_instance.instance.host)
  username = data.cloudamqp_credentials.local-admin.username
  password = data.cloudamqp_credentials.local-admin.password
}

resource "rabbitmq_user" "test_user" {
  name     = "test_user"
  password = "foobar"
  tags     = ["management", "policymaker"]
}
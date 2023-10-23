terraform {
  required_providers {
    cloudamqp = {
      source = "cloudamqp/cloudamqp"
      version = "~>1.0"
    }
  }
}

provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

resource "cloudamqp_instance" "instance" {
  name 				= "terraform-plugin-test"
  nodes 			= 1
  plan  			= "bunny-1"
  region 			= "amazon-web-services::us-east-1"
}

resource "cloudamqp_plugin" "mqtt_plugin" {
  instance_id = cloudamqp_instance.instance.id
  name = "rabbitmq_web_mqtt"
  enabled = true
}

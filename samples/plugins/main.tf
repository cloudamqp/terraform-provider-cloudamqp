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
  # Use to skip teardown of plugins for faster overall 'terraform destroy'
  # enable_faster_instance_destroy = true
}

resource "cloudamqp_instance" "instance" {
  name 				= "terraform-plugin-test"
  plan  			= "bunny-1"
  region 			= "amazon-web-services::us-east-1"
}

resource "cloudamqp_plugin" "mqtt_plugin" {
  instance_id = cloudamqp_instance.instance.id
  name = "rabbitmq_web_mqtt"
  enabled = true
}

resource "cloudamqp_plugin_community" "delayed_message_exchange" {
  instance_id = cloudamqp_instance.instance.id
  name = "rabbitmq_delayed_message_exchange"
  enabled = true
}
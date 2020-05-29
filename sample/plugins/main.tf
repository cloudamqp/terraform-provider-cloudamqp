provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

resource "cloudamqp_instance" "instance" {
  name 				= "terraform-plugin-test"
  nodes 			= 1
  plan  			= "bunny"
  region 			= "amazon-web-services::us-east-1"
  rmq_version = "3.8.2"
  tags 				= ["terraform"]
  vpc_subnet = "192.168.0.1/24"
}

resource "cloudamqp_plugin" "mqtt_plugin" {
  instance_id = cloudamqp_instance.instance.id
  name = "rabbitmq_web_mqtt"
  enabled = true
}

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

# Import default firewall settings
# "import cloudamqp_security_firewall.firewall <instance_id>"
# resource "cloudamqp_security_firewall" "firewall" {}

# Once imported, populate instance_id and rules attributes
# resource "cloudamqp_security_firewall" "firewall" {
#   instance_id = cloudamqp_instance.instance.id
#   rules {
#     ip = "0.0.0.0/0"
#     ports = []
#     services = ["STOMP", "AMQP", "MQTTS", "STOMPS", "MQTT", "AMQPS"]
#   }
# }

# Or overwrite the firewall settings directly with new settings.
resource "cloudamqp_security_firewall.firewall" {
  instance_id = cloudamqp_insntance.instance.id
  rules {
    ip = "192.168.0.0/0"
    ports = [4567]
    services = ["MQTTS", "STOMPS", "AMQPS"]
  }
}


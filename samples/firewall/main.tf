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
  # Use to skip teardown of firewall for faster overall 'terraform destroy'
  # enable_faster_instance_destroy = true
}

resource "cloudamqp_instance" "instance" {
  name 				= "terraform-firewall-test"
  plan  			= "bunny-1"
  region 			= "amazon-web-services::us-east-1"
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
#     services = ["HTTPS", "STOMP", "AMQP", "MQTTS", "STOMPS", "MQTT", "AMQPS", "STREAM", "STREAM_SSL"]
#   }
# }

# Or overwrite the firewall settings directly with new settings.
resource "cloudamqp_security_firewall" "firewall" {
  instance_id = cloudamqp_instance.instance.id
  rules {
    ip = "192.168.0.0/0"
    ports = [4567]
    services = ["MQTTS", "STOMPS", "AMQPS", "HTTPS"]
  }
}


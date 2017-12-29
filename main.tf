provider "cloudamqp" {}

resource "cloudamqp_instance" "instance1" {
  name   = "terraform-provider-test-instance-1"
  plan   = "rabbit"
  region = "amazon-web-services::us-east-1"
  nodes = 2
  vpc_subnet = "10.201.0.0/24"
  rmq_version = "3.6.12"
}

output "instance_name" {
  value = "${cloudamqp_instance.instance1.name}"
}

output "instance_url" {
  value = "${cloudamqp_instance.instance1.url}"
}

output "instance_apikey" {
  value = "${cloudamqp_instance.instance1.apikey}"
}

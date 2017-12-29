provider "cloudamqp" {}

resource "cloudamqp_instance" "instance1" {
  name   = "terraform-provider-test-instance-1"
  plan   = "lemur"
  region = "amazon-web-services::us-east-1"
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

provider "cloudamqp" {}

resource "cloudamqp_instance" "rmq_bunny" {
  name   = "terraform-provider-test"
  plan   = "bunny"
  region = "amazon-web-services::us-east-1"
  vpc_subnet = "10.201.0.0/24"
}

output "rmq_url" {
  value = "${cloudamqp_instance.rmq_bunny.url}"
}

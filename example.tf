provider "cloudamqp" {}

resource "cloudamqp_instance" "rmq_lemur" {
  name   = "terraform-provider-test-lemur"
  plan   = "lemur"
  region = "amazon-web-services::us-east-1"
}

output "rmq_lemur_url" {
  value = "${cloudamqp_instance.rmq_lemur.url}"
}


#resource "cloudamqp_instance" "rmq_bunny" {
#  name   = "terraform-provider-test-vpc-bunny"
#  plan   = "bunny"
#  region = "amazon-web-services::us-east-1"
#  vpc_subnet = "10.201.0.0/24"
#}
#
#output "rmq_bunny_url" {
#  value = "${cloudamqp_instance.rmq_bunny.url}"
#}

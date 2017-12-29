provider "cloudamqp" {}

resource "cloudamqp_instance" "my_instance" {
  name = "terraform-provider-test-instance-1"
  plan = "lemur"
  region = "amazon-web-services::us-east-1"
}

resource "cloudamqp_instance" "my_instance2" {
  name = "terraform-provider-test-instance-2"
  plan = "lemur"
  region = "amazon-web-services::us-east-1"
}

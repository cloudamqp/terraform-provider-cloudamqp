# Terraform provider for CloudAMQP

Setup your CloudAMQP cluster from Terraform

## Install

```sh
git clone https://github.com/cloudamqp/terraform-provider.git
cd terraform-provider
make depupdate
make init
```

Now the provider is installed in the terraform plugins folder and ready to be used.

## Example

```hcl
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

resource "cloudamqp_notification" "recipient_01" {
  instance_id = "${cloudamqp_instance.rmq_bunny.id}"
  type = "email"
  value = "alarm@example.com"
}

resource "cloudamqp_alarm" "alarm_01" {
  instance_id = "${cloudamqp_instance.rmq_bunny.id}"
  type = "cpu"
  value_threshold = 90
  time_threshold = 600
  notifications = ["${cloudamqp_notification.recipient_01.id}"]
}
```




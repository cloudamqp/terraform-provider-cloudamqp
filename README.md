# Terraform provider for CloudAMQP

Setup your CloudAMQP cluster from Terraform

## Install

```sh
# Go libraries
go get github.com/84codes/go-api/api
go get github.com/hashicorp/terraform/helper/schema
go get github.com/hashicorp/terraform/plugin
go get github.com/hashicorp/terraform/terraform
# Provider
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
```




# Terraform provider for CloudAMQP

Setup your CloudAMQP cluster from Terraform

## Prerequisite

Golang, Dep

### Mac

- brew install golang
- brew install dep

## Install

```sh
cd $GOPATH/src/github.com
mkdir cloudamqp
cd cloudamqp
git clone https://github.com/cloudamqp/terraform-provider.git
cd terraform-provider
make depupdate
make init
```

Now the provider is installed in the terraform plugins folder and ready to be used.

## Example

```hcl
provider "cloudamqp" {}

resource "cloudamqp_instance" "rmq_url"{
  name = "rmq_url"
  plan = "lemur"
  nodes = 1
  region = "amazon-web-services::us-east-1"
  rmq_version = "3.6.16"
  vpc_subnet = "10.201.0.0/24"
}

output "rmq_url" {
  value = "${cloudamqp_instance.rmq_url.url}"
}
```

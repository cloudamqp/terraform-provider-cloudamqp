# Terraform provider for CloudAMQP

Setup your CloudAMQP cluster from Terraform

## Getting Started (As a Terraform User)
### Prerequisites

* Install golang: https://golang.org/dl/

* Install go dep: https://golang.github.io/dep/docs/installation.html

* Install terraform: https://learn.hashicorp.com/terraform/getting-started/install.html

* Create a CloudAMQP account if you haven't already:

    * Go to https://www.cloudamqp.com/

    * Click "Sign Up"

    * Sign in

    * Go to API access (https://customer.cloudamqp.com/apikeys) and create a key.
      (note that this is the API key for one of the two APIs CloudAMQP supports.  
      See https://docs.cloudamqp.com/cloudamqp_api.html.  We will discuss the other
      later.)

### Install CloudAMQP Terraform Provider
```sh
git clone https://github.com/cloudamqp/terraform-provider.git
cd terraform-provider
make depupdate
make init
```

Now the provider is installed in the terraform plugins folder and ready to be used.

### Example Usage: Deploying a First CloudAMQP RMQ server

(See the examples.tf file in the repo.  It has a bunny VPC example and a simple lemur example.)

```sh
cd terraform-provider  #This is the root of the repo where examples.tf lives.
terraform plan
```
When prompted paste in your CloudAMQP API key (created above).

This will give you output on stdout that tells you what would have been created:
* rmq_lemur

Next run--
```sh
terraform apply
```

Again, paste in your API key.  This should create an actual CloudAMQP instance.

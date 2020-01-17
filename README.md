# Terraform provider for CloudAMQP

Setup your CloudAMQP cluster from Terraform

## Getting Started (As a Terraform User)

### Prerequisites

Golang, make, Terraform

## Install

* Install golang: https://golang.org/dl/
  Example with default paths
  * Download latest version and extract to `/usr/local/go`
  * Set environmental variable `export GOROOT=/usr/local/go`
  * Set environmental variable `export GOPATH=$HOME/go`
  * Set environmental variable `export PATH=$GOROOT/bin:$GOPATH:$PATH`
  * Activate module mode `export GO111MODULE=on` (Very importent!)
* Install make
  * `sudo apt install make`
* Install terraform: https://learn.hashicorp.com/terraform/getting-started/install.html
  * Download the latest version and extract to /usr/local/terraform
  * Set environmental variable `export PATH=/usr/local/terraform:$PATH`

* Create a CloudAMQP account if you haven't already:
  * Go to https://www.cloudamqp.com/
  * Click "Sign Up"
  * Sign in
  * Go to API access (https://customer.cloudamqp.com/apikeys) and create a key. (note that this is the API key for one of the two APIs CloudAMQP supports.)

`
The two APIs supported can be found at https://docs.cloudamqp.com (called customer) and https://docs.cloudamqp.com/cloudamqp_api.html (called api). The API key created gain access to the customer API (used to handle the instance). While the second API handles different resources on the instace (such as alarms, notification etc.). The customer API also has a proxy service, which makes it possible for the provider to access the second API through customer API using the same created API key.
`

### Install CloudAMQP Terraform Provider

```sh
go get -d -u -v github.com/cloudamqp/terraform-provider-cloudamqp
cd $GOPATH/src/github.com/cloudamqp/terraform-provider-cloudamqp
go get -u
make install
```

Now the provider is installed in the terraform plugins folder and ready to be used.

To update the dependencies, then run again.
`go get -u`

To clean up the .mod and .sum files from unused dependencies, run.
`go mod tidy`

More detailed documentation of the provider can be found at: https://docs.cloudamqp.com/cloudamqp_terraform.html

### Example Usage: Deploying a First CloudAMQP RMQ server

(See the examples.tf file in the repo.  It has a bunny VPC example and a simple lemur example.)

```sh
cd $GOPATH/src/github.com/cloudamqp/terraform-provider-cloudamqp  #This is the root of the repo where examples.tf lives.
terraform plan
```

When prompted paste in your CloudAMQP API key (created above).

This will give you output on stdout that tells you what would have been created:

* rmq_lemur

Next run

```sh
terraform apply
```

Again, paste in your API key.  This should create an actual CloudAMQP instance.

## Versioning

Enabled versioning to the Makefile, which also automatically adds it to the built provider. New name is therefore terraform-provider-cloudamqp_vx.y.z, where x.y.z is the version.

## Debug log

If more information needed, it's possible to increase Terraform log level. Using *DEBUG* will enable both CloudAMQP and underlying go-api debug logging.

To enable Terraform debug logging.
`export TF_LOG=DEBUG`

## Resources

Resource documentation can be found [here](https://docs.cloudamqp.com/cloudamqp_terraform.html)

## Import

Import existing infrastructure into state and bring the resource under Terraform management. Information about the resource will be added to the terraform.state file. Then add manually the given information to the .tf file. Once this is done, run terraform plan to see that the resource is under Terraform management. From here it's possible to add more resources such as alarm.

### Instance:

Import cloudamqp instance and bring it under Terraform management. First declare an empty instance resource in the .tf file. Followed by running the terraform import command

```sh
resource "cloudamqp_instance"."rmq_url" {}
```

Generic form of terraform import command

```sh
terraform import {resource_type}.{resource_name} {resource_id}
```

Example of terraform import command (with resource_id=80)

```sh
terraform import cloudamqp_instance.rmq_url 80
```

### Resources depending on an instance:

All resources depending on the instance resource also needs the instance id when using terraform import, in order to make correct API calls. Resource id and instance id is seperated with ",".

Resource affected by this is:

* cloudamqp_notification
* cloudmaqp_alarm

First declare two empty notification and alarm resources in the .tf file. Followed by running the terraform import command.

```sh
resource "cloudamqp_notification"."recipient_01" {}
resource "cloudamqp_alarm"."alarm_01" {}
```

Generic form of terraform import command

```sh
terraform import {resource_type}.{resource_name} {resource_id},{instance_id}
```

Example of terraform import command (with instance_id=80)

```sh
terraform import cloudamqp_notification.recipient_01 10,80
terraform import cloudamqp_alarm.alarm_01 65,80
```

## AWS VPC Setup

Support for setting up VPC peering connection between AWS instance and CloudAMQP. Requires that the AWS instance is used as the requester and CloudAMQP used as an accepter. More detailed description can be found here: [setup](https://docs.cloudamqp.com/cloudamqp_terraform.html#aws-vpc-setup)

Together with at full example found under *sample/aws_vpc*.

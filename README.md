# Terraform provider for CloudAMQP

Setup your CloudAMQP cluster from Terraform

## Getting Started (As a Terraform User)

### Prerequisites

Golang, make, Terraform

## Setup prerequisites and CloudAMQP account

* Install golang: https://golang.org/dl/
  Example with default paths
  * Download latest version and extract to `/usr/local/go`
  * Set environmental variable `export GOROOT=/usr/local/go`
  * Set environmental variable `export GOPATH=$HOME/go`
  * Set environmental variable `export PATH=$GOROOT/bin:$GOPATH:$GOPATH/bin:$PATH`
  * Activate module mode `export GO111MODULE=on` (Very important!)
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

## Install the CloudAMQP Terraform Provider

### From Terraform Registry

The CloudAMQP provider is available from the registry at https://registry.terraform.io/providers/cloudamqp/cloudamqp

If you are using Terraform 0.13+:

```yaml
terraform {
  required_providers {
    cloudamqp = {
      source = "cloudamqp/cloudamqp"
    }
  }
}
```

Read more at https://www.terraform.io/docs/language/providers/requirements.html

### From source

Clone repository to `$GOPATH/src/github.com/cloudamqp/terraform-provider-cloudamqp`

Change directory and build the provider from make. This will call `go install` and install the plugin under `$GOPATH/bin`.

```sh
$ cd $GOPATH/src/github.com/cloudamqp/terraform-provider-cloudamqp
$ make build
```

Run `terraform init` from the same folder as the tf.file is located. Terraform should also search in `$GOPATH/bin`. If this not the case, the provider needs to be manually installed by moving it to `$HOME/.terraform.d/plugins`. [Install plugins](https://www.terraform.io/docs/plugins/basics.html#installing-plugins). If `$HOME/.terraform.d/plugins` don't exists, the directory needs to be created.

```sh
$ mkdir -p ~/.terraform.d/plugins
$ cp $GOPATH/bin/terraform-provider-cloudamqp $HOME/.terraform.d/plugins/terraform-provider-cloudamqp
$ cd <path_to_tf_file>
$ terraform init
```

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


## Debug log

If more information needed, it's possible to increase Terraform log level. Using *DEBUG* will enable both CloudAMQP and underlying go-api debug logging.

To enable Terraform debug logging.
`export TF_LOG=DEBUG`

## Resources

Resource documentation can be found [here](https://docs.cloudamqp.com/cloudamqp_terraform.html)

### Instance ###

**IMPORTANT - PLAN CHANGES BETWEEN SHARED AND DEDICATED**
`
It’s possible to change between shared and dedicated plans (or vice versa). This will however force a destruction of the old instance, before creating a new one. All data will be lost and a new hostname will be created with corresponding DNS record.
`

## Import

Import existing infrastructure into state and bring the resource under Terraform management. Information about the resource will be added to the terraform.state file. Then add manually the given information to the .tf file. Once this is done, run terraform plan to see that the resource is under Terraform management. From here it's possible to add more resources such as alarm.

You'll need to determine the `resource_id` and other identifiers for the items you intend to import. You can do this using the CloudAMQP API, which is documented here https://docs.cloudamqp.com/#instances. Use a token found on this page https://customer.cloudamqp.com/apikeys.

### Instance:

Import cloudamqp instance and bring it under Terraform management. First declare an empty instance resource in the .tf file. Followed by running the terraform import command

```sh
resource "cloudamqp_instance" "rmq_url" {}
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
* cloudamqp_alarm

First declare two empty notification and alarm resources in the .tf file. Followed by running the terraform import command.

```sh
resource "cloudamqp_notification"."recipient_01" {}
resource "cloudamqp_alarm"."alarm_01" {}
```

Generic form of terraform import command

```sh
terraform import {resource_type}.{resource_name} {resource_id},{instance_id}
```

You can find `{instance_id}` through [the API](https://docs.cloudamqp.com/#instances)

Example of terraform import command (with instance_id=80)

```sh
terraform import cloudamqp_notification.recipient_01 10,80
terraform import cloudamqp_alarm.alarm_01 65,80
```

## AWS VPC Setup

Support for setting up VPC peering connection between AWS instance and CloudAMQP. Requires that the AWS instance is used as the requester and CloudAMQP used as an accepter. More detailed description can be found here: [setup](https://docs.cloudamqp.com/cloudamqp_terraform.html#aws-vpc-setup)

Together with at full example found under *sample/aws_vpc*.

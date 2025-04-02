---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_account_vpcs"
description: |-
  List all standalone VPCs available for an account.
---

# cloudamqp_account_vpcs

Use this data source to retrieve basic information about all standalone VPCs available for an
account. Uses the included apikey in provider configuration to determine which account to read from.

## Example Usage

Can be used in other resources/data sources when the VPC identifier is unknown, while other
attributes are known. E.g. find correct VPC using the `name` you gave your VPC. Then iterate over
VPCs to find the matching one and extract the VPC identifier.

```hcl
provider "cloudamqp" {
  apikey  = "<apikey>"
}

locals {
  my_vpc_name = "<your VPC name>"
}

data "cloudamqp_account_vpcs" "vpc_list" {}

output "vpc_id" {
  value = [for vpc in data.cloudamqp_account_vpcs.vpc_list.vpcs : vpc if vpc["name"] == local.my_vpc_name][0].id
}
```

## Attributes Reference

All attributes reference are computed

* `id`      - The identifier for this data source. Set to `na` since there is no unique identifier.
* `vpcs`    - An array of VPCs. Each `vpcs` block consists of the fields documented below.

___

The `vpcs` block consist of

* `id`          - The VPC identifier.
* `name`        - The VPC instance name.
* `region`      - The region the VPC is hosted in.
* `subnet`      - The VPC subnet.
* `tags`        - Optional tags set for the VPC.
* `vpc_name`    - VPC name given when hosted at the cloud provider.

## Dependency

This data source depends on apikey set in the provider configuration.

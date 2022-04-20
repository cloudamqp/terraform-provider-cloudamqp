---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_account"
description: |-
  List all instances available for an account.
---

# cloudamqp_account

Use this data source to retrieve basic information about all instances available for an account. Uses the included apikey in provider configuration, to determine which account to read from.

## Example Usage

Can be used in other resources/data sources when instance identifier is unknown, while other attributes are known. E.g. find correct instance from `instance name`. Then iterate over instances to find the matching one and extract the instance identifier.

```hcl
provider "cloudamqp" {
  apikey  = "<apikey>"
}

locals {
  instance_name = "<instance_name>"
}

data "cloudamqp_account" "instance_list" {}

output "instance_id" {
  instance_id = [for instance in data.cloudamqp_account.instance_list.instances : instance if instance["name"] == local.instance_name][0].id
}
```

## Attributes reference

All attributes reference are computed

* `id`          - The identifier for this data source. Set to `na` since there is no unique identifier.
* `instances`   - An array of instances. Each `instances` block consists of the fields documented below.

___

The `instances` block consist of

* `id`      - The instance identifier.
* `name`    - The name of the instance.
* `plan`    - The subscription plan used for the instance.
* `region`  - The region were the instanece is located in.
* `tags`    - Optional tags set for the instance.

## Dependency

This data source depends on apikey set in the provider configuration.

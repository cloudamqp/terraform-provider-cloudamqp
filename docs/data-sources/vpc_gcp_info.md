---
layout: "cloudamqp"
page_title: "CloudAMQP: data source vpc_gcp_info"
description: |-
  Get information about VPC hosted in GCP.
---

# cloudamqp_vpc_gcp_info

Use this data source to retrieve information about VPC for a CloudAMQP instance hosted in GCP.

## Example Usage

```hcl
data "cloudamqp_vpc_gcp_info" "vpc_info" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attributes reference

All attributes reference are computed

* `id`                  - The identifier for this resource.
* `name`                - The name of the VPC.
* `vpc_subnet`          - Dedicated VPC subnet.
* `network`             - VPC network uri.

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

---
layout: "cloudamqp"
page_title: "CloudAMQP: data source vpc_info"
description: |-
  Get information about VPC hosted in AWS.
---

# cloudamqp_vpc_info

Use this data source to retrieve information about VPC for a CloudAMQP instance.

Only available for CloudAMQP instances hosted in AWS.

## Example Usage

```hcl
data "cloudamqp_vpc_info" "vpc_info" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attributes reference

All attributes reference are computed

* `id`                  - The identifier for this resource.
* `name`                - The name of the CloudAMQP instance.
* `vpc_subnet`          - Dedicated VPC subnet.
* `owner_id`            - AWS account identifier.
* `security_group_id`   - AWS security group identifier.

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

---
layout: "cloudamqp"
page_title: "CloudAMQP: data source vpc_info"
description: |-
  Get information about VPC hosted in AWS.
---

# cloudamqp_vpc_info

Use this data source to retrieve information about VPC for a CloudAMQP instance. Depens on the identifier of the corresponding `cloudamqp_instance`resource or data source.

Only available for CloudAMQP instances hosted in AWS.

## Example Usage

```hcl
data "cloudamqp_vpc_info" "vpc_info" {
  instance_id = cloudamqp_instance.instance.id
}
```

## Argument reference

* `instance_id` - (Required) The CloudAMQP instance identifier.

## Attribute reference

* `name`                - (Computed) The name of the CloudAMQP instance.
* `vpc_subnet`          - (Computed) Dedicated VPC subnet.
* `owner_id`            - (Computed) AWS account identifier.
* `security_group_id`   - (Computed) AWS security group identifier.

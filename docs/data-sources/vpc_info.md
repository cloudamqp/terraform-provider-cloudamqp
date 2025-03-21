---
layout: "cloudamqp"
page_title: "CloudAMQP: data source vpc_info"
description: |-
  Get information about VPC hosted in AWS.
---

# cloudamqp_vpc_info

Use this data source to retrieve information about VPC for a CloudAMQP instance.

-> **Note:** Only available for CloudAMQP instances/VPCs hosted in AWS.

## Example Usage

<details>
  <summary>
    <b>
      <i>AWS VPC peering pre v1.16.0</i>
    </b>
  </summary>

```hcl
data "cloudamqp_vpc_info" "vpc_info" {
  instance_id = cloudamqp_instance.instance.id
}
```
</details>

<details>
  <summary>
    <b>
      <i>AWS VPC peering post v1.16.0 (Managed VPC)</i>
    </b>
  </summary>

```hcl
data "cloudamqp_vpc_info" "vpc_info" {
  vpc_id = cloudamqp_vpc.vpc.id
  # vpc_id prefered over instance_id
  # instance_id = cloudamqp_instance.instance.id
}
```
</details>


## Argument reference

 *Note: this resource require either `instance_id` or `vpc_id` from v1.16.0*

* `instance_id` - (Optional) The CloudAMQP instance identifier.

 ***Deprecated: Changed from required to optional in v1.16.0 will be removed in next major version (v2.0)***

* `vpc_id` - (Optional) The managed VPC identifier.

 ***Note: Added as optional in version v1.16.0 and will be required in next major version (v2.0)***

## Attributes reference

All attributes reference are computed

* `id`                  - The identifier for this resource.
* `name`                - The name of the CloudAMQP instance.
* `vpc_subnet`          - Dedicated VPC subnet.
* `owner_id`            - AWS account identifier.
* `security_group_id`   - AWS security group identifier.

## Dependency

*Pre v1.16.0*
This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

*Post v1.16.0*
This resource depends on CloudAMQP managed VPC identifier, `cloudamqp_vpc.vpc.id` or instance identifier, `cloudamqp_instance.instance.id`.

---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_vpc_peering"
description: |-
  Accepting VPC peering request from an AWS accepter.
---

# cloudamqp_vpc_peering

This resouce allows you to accepting VPC peering request from an AWS requester. Only available for CloudAMQP instance hosted in AWS. Depends on `cloudamqp_instance`resource and the instance identifier.

Only available for dedicated subscription plans.

## Example Usage

```hcl
resource "cloudamqp_vpc_peering" "vpc_accept_peering" {
  instance_id = cloudamqp_instance.instance_01.id
  peering_id = aws_vpc_peering_connection.aws_vpc_peering.id
}
```

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance ID.
* `peering_id` - (Required) Peering identifier created by AW peering request.

## Import

`cloudamqp_vpc_peering` can be imported using the CloudAMQP instance identifier.

`terraform import cloudamqp_vpc_peering.aws_vpc_peering <instance_id>`

---
layout: "cloudamqp"
page_title: "CloudAMQP: LavinMQ shared to dedicated upgrade"
subcategory: "info"
description: |-
  Guide on how to upgrade a LavinMQ shared instance to a dedicated plan in-place.
---

# LavinMQ shared to dedicated upgrade

From v1.45.0 it is possible to upgrade a [**LavinMQ**] shared instance to a dedicated plan
in-place. The instance ID is preserved, no resource replacement occurs, and existing definitions
are kept. The upgrade is performed by changing the `plan` argument in your configuration and
running `terraform apply`.

~> This upgrade path is only available for **LavinMQ** instances. RabbitMQ shared instances cannot
be upgraded in-place to dedicated plans and will require a new resource.

## Shared LavinMQ plans

The following LavinMQ shared plans are eligible for in-place upgrade to a dedicated plan:

Plan       | Name
-----------|----------------
`lemming`  | Loyal Lemming
`ermine`   | Elegant Ermine

## Constraints

* The cloud provider must remain the same. For example, moving between AWS regions is supported,
  but moving from AWS to GCP is not.
* The target plan must be a dedicated LavinMQ plan. See available [plans].
* The reverse upgrade path (dedicated to shared) is **not** supported and will force a new
  resource to be created.

## What can change during the upgrade

The following attributes can optionally be changed during the same `terraform apply` as the plan
change:

* `plan` — (Required) Must change to a dedicated LavinMQ plan (e.g. `puffin-1`, `penguin-1`,
  `wolverine-1`).
* `region` — (Optional) Can change to a different region within the same cloud provider.
* `vpc_id` — (Optional) Place the dedicated instance in an existing VPC.
* `vpc_subnet` — (Optional) Create a new VPC with the given subnet.
* `preferred_az` — (Optional) Preferred availability zones for the dedicated nodes.

## Upgrade example

<details>
  <summary>
    <b>
      <i>Upgrade a lemming instance to a dedicated wolverine plan</i>
    </b>
  </summary>

Start with a shared LavinMQ instance:

```hcl
resource "cloudamqp_instance" "instance" {
  name   = "terraform-cloudamqp-instance"
  plan   = "lemming"
  region = "amazon-web-services::eu-north-1"
  tags   = ["terraform"]
}
```

Change `plan` to a dedicated LavinMQ plan. Optionally change `region` within the same cloud
provider:

```hcl
resource "cloudamqp_instance" "instance" {
  name   = "terraform-cloudamqp-instance"
  plan   = "wolverine-1"
  region = "amazon-web-services::eu-west-1"
  tags   = ["terraform"]
}
```

Run `terraform apply`. The instance ID will remain the same and no resource replacement occurs.

</details>

<details>
  <summary>
    <b>
      <i>Upgrade a lemming instance to dedicated and place it in a managed VPC</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_vpc" "vpc" {
  name   = "my-vpc"
  region = "amazon-web-services::eu-west-1"
  subnet = "10.56.72.0/24"
  tags   = []
}

resource "cloudamqp_instance" "instance" {
  name                = "terraform-cloudamqp-instance"
  plan                = "wolverine-1"
  region              = "amazon-web-services::eu-west-1"
  tags                = ["terraform"]
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}
```

</details>

[**LavinMQ**]: https://lavinmq.com
[cloudamqp_instance]: ../resources/instance.md
[plans]: ../guides/info_plan.md

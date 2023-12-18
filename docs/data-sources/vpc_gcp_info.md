---
layout: "cloudamqp"
page_title: "CloudAMQP: data source vpc_gcp_info"
description: |-
  Get information about VPC hosted in GCP.
---

# cloudamqp_vpc_gcp_info

Use this data source to retrieve information about VPC for a CloudAMQP instance hosted in GCP.

## Example Usage

<details>
  <summary>
    <b>
      <i>AWS VPC peering pre v1.16.0</i>
    </b>
  </summary>

```hcl
data "cloudamqp_vpc_gcp_info" "vpc_info" {
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
data "cloudamqp_vpc_gcp_info" "vpc_info" {
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

* `sleep` - (Optional) Configurable sleep time (seconds) between retries when reading peering. Default set to 10 seconds.
* `timeout` - (Optional) - Configurable timeout time (seconds) before retries times out. Default set to 1800 seconds.

## Attributes reference

All attributes reference are computed

* `id`                  - The identifier for this resource.
* `name`                - The name of the VPC.
* `vpc_subnet`          - Dedicated VPC subnet.
* `network`             - VPC network uri.

## Dependency

*Pre v1.16.0*
This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

*Post v1.16.0*
This resource depends on CloudAMQP managed VPC identifier, `cloudamqp_vpc.vpc.id` or instance identifier, `cloudamqp_instance.instance.id`.

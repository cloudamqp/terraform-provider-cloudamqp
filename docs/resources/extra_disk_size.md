---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_extra_disk_size"
description: |-
  Resize the disk with extra storage capacity.
---

# cloudamqp_extra_disk_size

This resource allows you to resize the disk with additional storage capacity.

***Pre v1.25.0***: Only available for Amazon Web Services (AWS) and it done without downtime

***Post v1.25.0***: Now also available for Google-Compute-Engine (GCE) and Azure.

Introducing an additional optional argument called `allow_downtime`, default set to *false*. This will resize the disk without downtime for *AWS* and *GCE*. While Azure only support swapping the disk, and this argument needs to be set to *true*.

Allow downtime also makes it possible to circumvent the time rate limit or shrinking the disk.

| Cloud Platform        | allow_downtime=false | allow_downtime=true           |
|-----------------------|----------------------|-------------------------------|
| amazon-web-services   | Expand current disk  | Try to expand, otherwise swap |
| google-compute-engine | Expand current disk  | Try to expand, otherwise swap |
| azure-arm             | Not supported        | Swap disk to new size         |

~> **WARNING:** Due to restrictions from cloud providers, it's only possible to resize the disk every 8 hours. Unless the `allow_downtime=true` is set. Then the disk will be swapped for a new.

Pricing is available at [cloudamqp.com](https://www.cloudamqp.com/).

Only available for dedicated subscription plans.

## Example Usage

<details>
  <summary>
    <b>
      <i>AWS extra disk size (pre v1.25.0)</i>
    </b>
  </summary>

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# Instance
resource "cloudamqp_instance" "instance" {
  name   = "Instance"
  plan   = "bunny-1"
  region = "amazon-web-services::us-west-2"
  rmq_version = "3.10.1"
}

# Resize disk with 25 extra GB
resource "cloudamqp_extra_disk_size" "resize_disk" {
  instance_id = cloudamqp_instance.instance.id
  extra_disk_size = 25
}

# Refresh nodes info after disk resize
data "cloudamqp_nodes" "nodes" {
  instance_id = cloudamqp_instance.instance.id
  depends_on = [
    cloudamqp_extra_disk_size.resize_disk,
  ]
}
```

</details>

<details>
  <summary>
    <b>
      <i>AWS extra disk size without downtime</i>
    </b>
  </summary>

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# Instance
resource "cloudamqp_instance" "instance" {
  name   = "Instance"
  plan   = "bunny-1"
  region = "amazon-web-services::us-west-2"
  rmq_version = "3.10.1"
}

# Resize disk with 25 extra GB, without downtime
resource "cloudamqp_extra_disk_size" "resize_disk" {
  instance_id = cloudamqp_instance.instance.id
  extra_disk_size = 25
  allow_downtime = false
}

# Refresh nodes info after disk resize
data "cloudamqp_nodes" "nodes" {
  instance_id = cloudamqp_instance.instance.id
  depends_on = [
    cloudamqp_extra_disk_size.resize_disk,
  ]
}
```

</details>

<details>
  <summary>
    <b>
      <i>GCE extra disk size without downtime</i>
    </b>
  </summary>

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# Instance
resource "cloudamqp_instance" "instance" {
  name   = "Instance"
  plan   = "bunny-1"
  region = "google-compute-engine::us-central1"
  rmq_version = "3.10.1"
}

# Resize disk with 25 extra GB, without downtime
resource "cloudamqp_extra_disk_size" "resize_disk" {
  instance_id = cloudamqp_instance.instance.id
  extra_disk_size = 25
  allow_downtime = false
}

# Refresh nodes info after disk resize
data "cloudamqp_nodes" "nodes" {
  instance_id = cloudamqp_instance.instance.id
  depends_on = [
    cloudamqp_extra_disk_size.resize_disk,
  ]
}
```

</details>

<details>
  <summary>
    <b>
      <i>Azure extra disk size with downtime</i>
    </b>
  </summary>

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# Instance
resource "cloudamqp_instance" "instance" {
  name   = "Instance"
  plan   = "bunny-1"
  region = "azure-arm::centralus"
  rmq_version = "3.10.1"
}

# Resize disk with 25 extra GB, with downtime
resource "cloudamqp_extra_disk_size" "resize_disk" {
  instance_id = cloudamqp_instance.instance.id
  extra_disk_size = 25
  allow_downtime = true
}

# Refresh nodes info after disk resize
data "cloudamqp_nodes" "nodes" {
  instance_id = cloudamqp_instance.instance.id
  depends_on = [
    cloudamqp_extra_disk_size.resize_disk,
  ]
}
```

</details>

## Argument Reference

* `instance_id`       - (Required) The CloudAMQP instance ID.
* `extra_disk_size`   - (Required) Extra disk size in GB. Supported values: 25, 50, 100, 250, 500, 1000, 2000
* `allow_downtime`    - (Optional) When resizing the disk, allow downtime to do so. Default set to false. (Only Available from v1.25.0).
* `sleep`             - (Optional) Configurable sleep time in seconds between retries for firewall configuration. Default set to 30 seconds.
* `timeout`           - (Optional) Configurable timeout time in seconds for firewall configuration. Default set to 1800 seconds.

## Attributes reference

All attributes reference are computed

* `id`    - The identifier for this resource.
* `nodes` - An array of node information. Each `nodes` block consists of the fields documented below.

___

The `nodes` block consist of

* `name`                  - Name of the node.
* `disk_size`             - Subscription plan disk size
* `additional_disk_size`  - Additional added disk size

***Note:*** *Total disk size = disk_size + additional_disk_size*

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

Not possible to import this resource.

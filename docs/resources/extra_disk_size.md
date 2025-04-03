---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_extra_disk_size"
description: |-
  Resize the disk with extra storage capacity.
---

# cloudamqp_extra_disk_size

This resource allows you to resize the disk with additional storage capacity.

***Before v1.25.0***: Only available for Amazon Web Services (AWS) without downtime.

***From v1.25.0***: Google Compute Engine (GCE) and Azure available.

Introducing a new optional argument called `allow_downtime`. Leaving it out or set it to false will
proceed to try and resize the disk without downtime, available for *AWS*, *GCE* and *Azure*.

`allow_downtime` also makes it possible to circumvent the time rate limit or shrinking the disk.

| Cloud Platform        | allow_downtime=false | allow_downtime=true           | Possible to resize |
|-----------------------|----------------------|-------------------------------|--------------------|
| amazon-web-services   | Expand current disk* | Try to expand, otherwise swap | Every 6 hour       |
| google-compute-engine | Expand current disk* | Try to expand, otherwise swap | Every 4 hour       |
| azure-arm             | Expand current disk* | Expand current disk           | No time rate limit |

*Preferable method to use.

-> **Note:** Due to restrictions from cloud providers, it's only possible to resize the disk after
the rate time limit. See `Possible to resize` column above for the different cloud platforms.

-> **Note:** Shrinking the disk will always need to swap the old disk to a new one and require
`allow_downtime` set to *true*.

Pricing is available at [CloudAMQP] and only available for dedicated subscription plans.

## Example Usage

<details>
  <summary>
    <b>
      <i>AWS extra disk size (before v1.25.0)</i>
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
  plan   = "penguin-1"
  region = "amazon-web-services::us-west-2"
}

# Resize disk with 25 extra GB
resource "cloudamqp_extra_disk_size" "resize_disk" {
  instance_id     = cloudamqp_instance.instance.id
  extra_disk_size = 25
}

# Optional, refresh nodes info after disk resize by adding dependency
# to cloudamqp_extra_disk_size.resize_disk resource
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
  plan   = "penguin-1"
  region = "amazon-web-services::us-west-2"
}

# Resize disk with 25 extra GB, without downtime
resource "cloudamqp_extra_disk_size" "resize_disk" {
  instance_id     = cloudamqp_instance.instance.id
  extra_disk_size = 25
}

# Optional, refresh nodes info after disk resize by adding dependency
# to cloudamqp_extra_disk_size.resize_disk resource
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
  plan   = "penguin-1"
  region = "google-compute-engine::us-central1"
}

# Resize disk with 25 extra GB, without downtime
resource "cloudamqp_extra_disk_size" "resize_disk" {
  instance_id     = cloudamqp_instance.instance.id
  extra_disk_size = 25
}

# Optional, refresh nodes info after disk resize by adding dependency
# to cloudamqp_extra_disk_size.resize_disk resource
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
      <i>Azure extra disk size without downtime</i>
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
  plan   = "penguin-1"
  region = "azure-arm::centralus"
}

# Resize disk with 25 extra GB, with downtime
resource "cloudamqp_extra_disk_size" "resize_disk" {
  instance_id     = cloudamqp_instance.instance.id
  extra_disk_size = 25
}

# Optional, refresh nodes info after disk resize by adding dependency
# to cloudamqp_extra_disk_size.resize_disk resource
data "cloudamqp_nodes" "nodes" {
  instance_id = cloudamqp_instance.instance.id

  depends_on = [
    cloudamqp_extra_disk_size.resize_disk,
  ]
}
```

</details>

## Argument Reference

Any changes to the arguments will destroy and recreate this resource.

* `instance_id`     - (ForceNew/Required) The CloudAMQP instance ID.
* `extra_disk_size` - (ForceNew/Required) Extra disk size in GB. Supported values: 0, 25, 50, 100,
                        250, 500, 1000, 2000
* `allow_downtime`  - (Optional) When resizing the disk, allow cluster downtime if necessary.
                      Default set to false.
* `sleep`           - (Optional) Configurable sleep time in seconds between retries for resizing the
                      disk. Default set to 30 seconds.
* `timeout`         - (Optional) Configurable timeout time in seconds for resizing the disk. Default
                      set to 1800 seconds.

  ***Note:*** `allow_downtime`, `sleep`, `timeout` only available from [v1.25.0].

## Attributes Reference

All attributes reference are computed

* `id`    - The identifier for this resource.
* `nodes` - An array of node information. Each `nodes` block consists of the fields documented below.

___

The `nodes` block consist of

* `name`                  - Name of the node.
* `disk_size`             - Subscription plan disk size
* `additional_disk_size`  - Additional added disk size

  ***Note:*** Total disk size = disk_size + additional_disk_size

## Dependency

This data source depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

Not possible to import this resource.

[CloudAMQP]: https://www.cloudamqp.com/
[v1.25.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.25.0

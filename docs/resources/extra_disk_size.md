---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_extra_disk_size"
description: |-
  Resize the disk with extra storage capacity.
---

# cloudamqp_extra_disk_size

This resource allows you to expand the disk with additional storage capacity.

Only available for dedicated subscription plans hosted at Amazon Web Services (AWS) at this time.

⚠️  Due to restrictions from cloud providers, it's only possible to resize the disk every 8 hours.

## Example Usage

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# Instance
resource "cloudamqp_instance" "instance" {
  name   = "Instance"
  plan   = "squirrel-1"
  region = "amazon-web-services::us-west-2"
  rmq_version = "3.10.1"
}

# Resize disk with 25 extra GB
resource "cloudamqp_extra_disk_size" "resize_disk" {
  instance_id = cloudamqp_instance.instance.id
  extra_disk_size = 25
  depends_on = [
    cloudamqp_instance.instance,
  ]
}

# Refresh nodes info after disk resize
data "cloudamqp_nodes" "nodes" {
  instance_id = cloudamqp_instance.instance.id
  depends_on = [
    cloudamqp_extra_disk_size.resize_disk,
  ]
}
```

## Argument Reference

* `instance_id`       - (Required/ForceNew) The CloudAMQP instance ID.
* `extra_disk_size`   - (Required/ForceNew) Extra disk size in GB. Supported values: 25, 50, 100, 250, 500, 1000, 2000

## Import

Not possible to import this resource.

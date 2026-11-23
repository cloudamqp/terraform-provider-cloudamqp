---
layout: "cloudamqp"
page_title: "CloudAMQP: instance regions"
subcategory: "info"
description: |-
  Available regions for CloudAMQP instances
---

# Instance region

CloudAMQP supports hosting by multiple cloud platform providers and over multiple regions. Below are
a few examples of supported platforms and regions. For a fully updated list, see [CloudAMQP plans],
scroll to the bottom and expand `List all available regions`.

Platforms and regions with shared servers are also listed, for AWS we
try to have at least one shared server supported for each region.

-> To retrieve an up to date list check out [CloudAMQP API regions]

The format used for instance regions is as follows: `{provider}::{region}`

```hcl
# Example of Amazon Web Services regions
amazon-web-services::us-east-1
amazon-web-services::us-west-1
amazon-web-services::eu-central-1
amazon-web-services::ap-east-1

# Example of Azure-arm regions
azure-arm::australiacentral
azure-arm::southeastasia

# Example of Google Compute Engine regions
google-compute-engine::us-central1
google-compute-engine::us-east1
google-compute-engine::europe-west1
google-compute-engine::asia-east1

# Example of Digital Ocean regions
digital-ocean::nyc3
digital-ocean::sgp1
```

[CloudAMQP plans]: https://www.cloudamqp.com/plans.html
[CloudAMQP API regions]: https://docs.cloudamqp.com/index.html#tag/regions

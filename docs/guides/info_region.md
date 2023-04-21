---
layout: "cloudamqp"
page_title: "CloudAMQP: instance regions"
subcategory: "info"
description: |-
  Available regions for CloudAMQP instances
---

# Instance region

CloudAMQP support hosting by multiple cloud platform providers and over multiple regions. Below a few examples of supported platforms and regions. For fully updated list see [CloudAMQP plans](https://www.cloudamqp.com/plans.html) and scroll to the bottom and extend `List all available regions`. Platforms and regions with shared servers are also listed, for AWS we try to have at least one shared server supported for each region.

*To retrieve an up to date list check out [cloudamqp-docs](https://docs.cloudamqp.com/#regions)*

Format used on instance regions are as follow `{provider}::{region}`

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

---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_integration_metric"
description: |-
  Creates and manages third party metrics integration for a CloudAMQP instance.
---

# cloudamqp_integration_metric

This resource allows you to create and manage, forwarding metrics to third party integrations for a CloudAMQP instance. Once configured, the metrics produced will be forward to corresponding integration. This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

Only available for dedicated subscription plans.

## Example usage

```hcl
resource "cloudamqp_integration_metric" "cloudwatch" {
  instance_id = cloudamqp_instance.instance.id
  name = "cloudwatch"
  access_key_id = var.aws_access_key_id
  secret_access_key = var_aws_secret_acccess_key
  region = "us-east-1"
}
```

## Argument references

The following arguments are supported:

* `name`              - (Required) The name of the third party log integration. See `Integration service reference`
* `region`            - (Optional) Region hosting the integration service.
* `access_key_id`     - (Optional) AWS access key identifier.
* `secret_access_key` - (Optional) AWS secret access key.
* `api_key`           - (Optional) The API key for the integration service.
* `email`             - (Optional) The email address registred for the integration service.
* `licene_key`        - (Optional) The license key registred for the integration service.
* `tags`              - (Optional) Tags. e.g. env=prod, region=europe.
* `queue_whitelist`   - (Optional) Whitelist queuesa using regular expression.
* `vhost_whitelist`   - (Optional) Whitelist vhost using regular expression.

This is the full list of all arguments. But a subset of arguments are used based on which type of integration used. See table below for more information.

## Integration service references

Valid names for third party log integration.

| Name          | Description |
|---------------|---------------------------------------------------------------|
| cloudwatch    | Create an IAM with programmatic access. |
| cloudwatch_v2 | Create an IAM with programmatic access. |
| datadog       | Create a Datadog API key at app.datadoghq.com |
| datadog_v2    | Create a Datadog API key at app.datadoghq.com
| librato       | Create a new API token (with record only permissions) here: https://metrics.librato.com/tokens |
| newrelic      | Deprecated! |
| newrelic_v2   | Find or register an Insert API key for your account: Go to insights.newrelic.com > Manage data > API keys. |

## Integration type reference

Valid arguments for third party log integrations.

Required arguments for all integrations: name
Optional arguments for all integrations: tags, queue_whitelist, vhost_whitelist

| Name | Type | Required arguments |
| ---- | ---- | ---- |
| Cloudwatch             | cloudwatch     | region, access_key_id, secret_access_key |
| Cloudwatch v2          | cloudwatch_v2  | region, access_key_id, secret_access_key |
| Datadog                | datadog        | api_key, region |
| Datadog v2             | datadog_v2     | api_key, region |
| Librato                | librato        | email, api_key |
| New relic (deprecated) | newrelic       | - |
| New relic v2           | newrelic_v2    | license_key |

## Import
`cloudamqp_integration_metric`can be imported using the name argument of the resource together with CloudAMQP instance identifier. The name and identifier are CSV separated, see example below.

`terraform import cloudamqp_integration_metric.cloudwatch <name>,<instance_id>`

---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_integration_metric"
description: |-
  Creates and manages third party metrics integration for a CloudAMQP instance.
---

# cloudamqp_integration_metric

This resource allows you to create and manage, forwarding metrics to third party integrations for a CloudAMQP instance. Once configured, the metrics produced will be forward to corresponding integration.

Only available for dedicated subscription plans.

## Example usage

<details>
  <summary>
    <b>
      <i>Cloudwatch v1 and v2 metric integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_metric" "cloudwatch" {
  instance_id = cloudamqp_instance.instance.id
  name = "cloudwatch"
  access_key_id = var.aws_access_key_id
  secret_access_key = var_aws_secret_acccess_key
  region = var.aws_region
}

resource "cloudamqp_integration_metric" "cloudwatch_v2" {
  instance_id = cloudamqp_instance.instance.id
  name = "cloudwatch_v2"
  access_key_id = var.aws_access_key_id
  secret_access_key = var_aws_secret_acccess_key
  region = var.aws_region
}
```
</details>

<details>
  <summary>
    <b>
      <i>Datadog v1 and v2 metric integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_metric" "datadog" {
  instance_id = cloudamqp_instance.instance.id
  name = "datadog"
  api_key = var.datadog_api_key
  region = var.datadog_region
}

resource "cloudamqp_integration_metric" "datadog_v2" {
  instance_id = cloudamqp_instance.instance.id
  name = "datadog_v2"
  api_key = var.datadog_api_key
  region = var.datadog_region
}
```
</details>

<details>
  <summary>
    <b>
      <i>Librato metric integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_metric" "librato" {
  instance_id = cloudamqp_instance.instance.id
  name = "librato"
  email = var.librato_email
  api_key = var.librato_api_key
}
```
</details>

<details>
  <summary>
    <b>
      <i>New relic v2 metric integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_metric" "newrelic" {
  instance_id = cloudamqp_instance.instance.id
  name = "newrelic_v2"
  api_key = var.newrelic_api_key
  region = var.newrelic_region
}
```
</details>

<details>
  <summary>
    <b>
      <i>Stackdriver metric integration (v1.20.2 or earlier versions)</i>
    </b>
  </summary>

Use variable file populated with project_id, private_key and client_email

```hcl
resource "cloudamqp_integration_metric" "stackdriver" {
  instance_id = cloudamqp_instance.instance.id
  name = "stackdriver"
  project_id = var.stackdriver_project_id
  private_key = var.stackdriver_private_key
  client_email = var.stackriver_email
}
```

or by using google_service_account_key resource from Google provider

```hcl
resource "google_service_account" "service_account" {
  account_id = "<account_id>"
  description = "<description>"
  display_name = "<display_name>"
}

resource "google_service_account_key" "service_account_key" {
  service_account_id = google_service_account.service_account.name
}

resource "cloudamqp_integration_metric" "stackdriver" {
  instance_id = cloudamqp_instance.instance.id
  name = "stackdriver"
  project_id = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).project_id
  private_key = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).private_key
  client_email = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).client_email
}
```
</details>

<details>
  <summary>
    <b>
      <i>Stackdriver metric integration (v1.30.0 or newer versions)</i>
    </b>
  </summary>

Use credentials argument and let the provider do the Base64decode and internally populate, *project_id, client_name, private_key*

```hcl
resource "google_service_account" "service_account" {
  account_id = "<account_id>"
  description = "<description>"
  display_name = "<display_name>"
}

resource "google_service_account_key" "service_account_key" {
  service_account_id = google_service_account.service_account.name
}

resource "cloudamqp_integration_metric" "stackdriver" {
  instance_id = cloudamqp_instance.instance.id
  name = "stackdriver"
  credentials = google_service_account_key.service_account_key.private_key
}
```

or use the same as earlier version and decode the google service account key

```hcl
resource "google_service_account" "service_account" {
  account_id = "<account_id>"
  description = "<description>"
  display_name = "<display_name>"
}

resource "google_service_account_key" "service_account_key" {
  service_account_id = google_service_account.service_account.name
}

resource "cloudamqp_integration_metric" "stackdriver" {
  instance_id = cloudamqp_instance.instance.id
  name = "stackdriver"
  project_id = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).project_id
  private_key = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).private_key
  client_email = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).client_email
}
```
</details>

## Argument references

The following arguments are supported:

* `name`              - (Required) The name of the third party log integration. See `Integration service reference`
* `region`            - (Optional) Region hosting the integration service.
* `access_key_id`     - (Optional) AWS access key identifier.
* `secret_access_key` - (Optional) AWS secret access key.
* `api_key`           - (Optional) The API key for the integration service.
* `email`             - (Optional) The email address registred for the integration service.
* `credentials`       - (Optional) Google Service Account private key credentials.
* `project_id`        - (Optional/Computed) The project identifier.
* `private_key`       - (Optional/Computed) The private access key.
* `client_email`      - (Optional/Computed) The client email registered for the integration service.
* `tags`              - (Optional) Tags. e.g. env=prod, region=europe.
* `queue_allowlist`   - (Optional) Allowlist queues using regular expression. Leave empty to include all queues.
* `vhost_allowlist`   - (Optional) Allowlist vhost using regular expression. Leave empty to include all vhosts.
* `queue_whitelist`   - **Deprecated** Use queue_allowlist instead
* `vhost_whitelist`   - **Deprecated** Use vhost_allowlist instead

This is the full list of all arguments. Only a subset of arguments are used based on which type of integration used. See [Integration type reference](#integration-type-reference) below for more information.

## Integration service references

Valid names for third party log integration.

| Name          | Description |
|---------------|---------------------------------------------------------------|
| cloudwatch    | Create an IAM with programmatic access. |
| cloudwatch_v2 | Create an IAM with programmatic access. |
| datadog       | Create a Datadog API key at app.datadoghq.com |
| datadog_v2    | Create a Datadog API key at app.datadoghq.com |
| librato       | Create a new API token (with record only permissions) here: https://metrics.librato.com/tokens |
| newrelic      | Deprecated! |
| newrelic_v2   | Find or register an Insert API key for your account: Go to insights.newrelic.com > Manage data > API keys. |
| stackdriver   | Create a service account and add 'monitor metrics writer' role from your Google Cloud Account |

## Integration type reference

Valid arguments for third party log integrations.

Required arguments for all integrations: *name*</br>
Optional arguments for all integrations: *tags*, *queue_allowlist*, *vhost_allowlist*

| Name | Type | Required arguments |
| ---- | ---- | ---- |
| Cloudwatch             | cloudwatch     | region, access_key_id, secret_access_key |
| Cloudwatch v2          | cloudwatch_v2  | region, access_key_id, secret_access_key |
| Datadog                | datadog        | api_key, region |
| Datadog v2             | datadog_v2     | api_key, region |
| Librato                | librato        | email, api_key |
| New relic (deprecated) | newrelic       | - |
| New relic v2           | newrelic_v2    | api_key, region |
| Stackdriver            | stackdriver    | credentials |

***Note:*** Stackdriver (v1.20.2 or earlier versions) required arguments  : project_id, private_key, client_email

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_integration_metric`can be imported using the resource identifier together with CloudAMQP instance identifier. The name and identifier are CSV separated, see example below.

`terraform import cloudamqp_integration_metric.<resource_name> <resource_id>,<instance_id>`

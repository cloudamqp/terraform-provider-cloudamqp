---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_integration_log"
description: |-
  Creates and manages third party log integration for a CloudAMQP instance.
---

# cloudamqp_integration_log

This resource allows you to create and manage third party log integrations for a CloudAMQP instance.
Once configured, the logs produced will be forward to corresponding integration.

Only available for dedicated subscription plans.

## Example Usage

<details>
  <summary>
    <b>
      <i>Azure monitor log integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log" "azure_monitor" {
  instance_id = cloudamqp_instance.instance.id
  name        = "azure_monitor"
  tenant_id = var.azm_tentant_id
  application_id = var.azm_application_id
  application_secret = var.azm_application_secret
  dce_uri = var.azm_dce_uri
  table = var.azm_table
  dcr_id = var.azm_dcr_id
}
```

</details>

<details>
  <summary>
    <b>
      <i>Cloudwatch log integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log" "cloudwatch" {
  instance_id = cloudamqp_instance.instance.id
  name = "cloudwatchlog"
  access_key_id = var.aws_access_key_id
  secret_access_key = var.aws_secret_access_key
  region = var.aws_region
}
```

</details>

<details>
  <summary>
    <b>
      <i>Coralogix log integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log" "coralogix" {
  instance_id = cloudamqp_instance.instance.id
  name        = "coralogix"
  private_key = var.coralogix_send_data_key
  endpoint    = var.coralogix_endpoint
  application = var.coralogix_application
  subsystem   = cloudamqp_instance.instance.host
}
```

</details>

<details>
  <summary>
    <b>
      <i>Datadog log integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log" "datadog" {
  instance_id = cloudamqp_instance.instance.id
  name = "datadog"
  region = var.datadog_region
  api_key = var.datadog_api_key
  tags = var.datadog_tags
}
```

</details>

<details>
  <summary>
    <b>
      <i>Logentries log integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log" "logentries" {
  instance_id = cloudamqp_instance.instance.id
  name = "logentries"
  token = var.logentries_token
}
```

</details>

<details>
  <summary>
    <b>
      <i>Loggly log integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log" "loggly" {
  instance_id = cloudamqp_instance.instance.id
  name = "loggly"
  token = var.loggly_token
}
```
</details>

<details>
  <summary>
    <b>
      <i>Papertrail log integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log" "papertrail" {
  instance_id = cloudamqp_instance.instance.id
  name = "papertrail"
  url = var.papertrail_url
}
```

</details>

<details>
  <summary>
    <b>
      <i>Scalyr log integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log" "scalyr" {
  instance_id = cloudamqp_instance.instance.id
  name = "scalyr"
  token = var.scalyr_token
  host = var.scalyr_host
}
```

<details>
  <summary>
    <b>
      <i>Splunk log integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log" "splunk" {
  instance_id = cloudamqp_instance.instance.id
  name = "splunk"
  token = var.splunk_token
  host_port = var.splunk_host_port
  source_type = "generic_single_line"
}
```

</details>

</details>

<details>
  <summary>
    <b>
      <i>Stackdriver log integration (v1.20.2 or older versions)</i>
    </b>
  </summary>

Use variable file populated with project_id, private_key and client_email

```hcl
resource "cloudamqp_integration_log" "stackdriver" {
  instance_id = cloudamqp_instance.instance.id
  name = "stackdriver"
  project_id = var.stackdriver_project_id
  private_key = var.stackdriver_private_key
  client_email = var.stackdriver_client_email
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

resource "cloudamqp_integration_log" "stackdriver" {
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
      <i>Stackdriver log integration (v1.21.0 or newer versions)</i>
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

resource "cloudamqp_integration_log" "stackdriver" {
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

resource "cloudamqp_integration_log" "stackdriver" {
  instance_id = cloudamqp_instance.instance.id
  name = "stackdriver"
  project_id = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).project_id
  private_key = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).private_key
  client_email = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).client_email
}
```

</details>

## Argument Reference

The following arguments are supported:

* `name`              - (Required) The name of the third party log integration. See
[Integration type reference](#integration-type-reference)

* `url`               - (Optional) Endpoint to log integration.
* `host_port`         - (Optional) Destination to send the logs.
* `token`             - (Optional/Sensitive) Token used for authentication.
* `region`            - (Optional) Region hosting the integration service.
* `access_key_id`     - (Optional/Sensitive) AWS access key identifier.
* `secret_access_key` - (Optional/Sensitive) AWS secret access key.
* `api_key`           - (Optional/Sensitive) The API key.
* `tags`              - (Optional) Tag the integration, e.g. env=prod,region=europe.
* `credentials`       - (Optional/Sensitive) Google Service Account private key credentials.
* `project_id`        - (Optional/Computed) The project identifier.
* `private_key`       - (Optional/Computed/Sensitive) The private access key.
* `client_email`      - (Optional/Computed) The client email registered for the integration service.
* `host`              - (Optional) The host for Scalyr integration. (app.scalyr.com, app.eu.scalyr.com)
* `sourcetype`        - (Optional) Assign source type to the data exported, eg. generic_single_line. (Splunk)
* `endpoint`          - (Optional) The syslog destination to send the logs to for Coralogix.
* `application`       - (Optional) The application name for Coralogix.
* `subsystem`         - (Optional) The subsystem name for Coralogix.
* `tenant_id`         - (Optional) The tenant identifier for Azure monitor.
* `application_id`    - (Optional) The application identifier for Azure monitor.
* `application_secret` - (Optional/Sensitive) The application secret for Azure monitor.
* `dce_uri`           - (Optional) The data collection endpoint for Azure monitor.
* `table`             - (Optional) The table name for Azure monitor.
* `dcr_id`            - (Optional) ID of data collection rule that your DCE is linked to for Azure Monitor.

This is the full list of all arguments. Only a subset of arguments are used based on which type of integration used. See [Integration Type reference](#integration-type-reference) table below for more information.

## Integration Type reference

Valid arguments for third party log integrations. See more information at [docs.cloudamqp.com](https://docs.cloudamqp.com/cloudamqp_api.html#add-log-integration)

Required arguments for all integrations: name

| Integration | name | Required arguments |
| ---- | ---- | ---- |
| Azure monitor | azure_monitor | tenant_id, application_id, application_secret, dce_uri, table, dcr_id |
| CloudWatch | cloudwatchlog | access_key_id, secret_access_key, region |
| Coralogix | coralogix | private_key, endpoint, application, subsystem |
| Data Dog | datadog | region, api_keys, tags |
| Log Entries | logentries | token |
| Loggly | loggly | token |
| Papertrail | papertrail | url |
| Scalyr | scalyr | token, host |
| Splunk | splunk | token, host_port, sourcetype |
| Stackdriver | stackdriver | credentials |

***Note:*** Stackdriver (v1.20.2 or earlier versions) required arguments  : project_id, private_key, client_email

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_integration_log`can be imported using the resource identifier together with CloudAMQP instance identifier. The name and identifier are CSV separated, see example below.

`terraform import cloudamqp_integration_log.<resource_name> <id>,<instance_id>`

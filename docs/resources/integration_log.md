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
  instance_id         = cloudamqp_instance.instance.id
  name                = "azure_monitor"
  tenant_id           = var.azm_tentant_id
  application_id      = var.azm_application_id
  application_secret  = var.azm_application_secret
  dce_uri             = var.azm_dce_uri
  table               = var.azm_table
  dcr_id              = var.azm_dcr_id
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
  instance_id       = cloudamqp_instance.instance.id
  name              = "cloudwatchlog"
  access_key_id     = var.aws_access_key_id
  secret_access_key = var.aws_secret_access_key
  region            = var.aws_region
}
```

</details>

<details>
  <summary>
    <b>
      <i>Cloudwatch log integration with retention and tags (from [v1.38.0])</i>
    </b>
  </summary>

Use retention and/or tags on the integration to make changes to `CloudAMQP` Log Group.

```hcl
resource "cloudamqp_integration_log" "cloudwatch" {
  instance_id       = cloudamqp_instance.instance.id
  name              = "cloudwatchlog"
  access_key_id     = var.aws_access_key_id
  secret_access_key = var.aws_secret_access_key
  region            = var.aws_region
  retention         = 14
  tags              = "Project=A,Environment=Development"
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
  name        = "datadog"
  region      = var.datadog_region
  api_key     = var.datadog_api_key
  tags        = "env=prod,region=us1,version=v1.0"
}
```

</details>

<details>
  <summary>
    <b>
      <i>Log entries log integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log" "logentries" {
  instance_id = cloudamqp_instance.instance.id
  name        = "logentries"
  token       = var.logentries_token
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
  name        = "loggly"
  token       = var.loggly_token
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
  name        = "papertrail"
  url         = var.papertrail_url
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
  name        = "scalyr"
  token       = var.scalyr_token
  host        = var.scalyr_host
}
```

</details>

<details>
  <summary>
    <b>
      <i>Splunk log integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log" "splunk" {
  instance_id = cloudamqp_instance.instance.id
  name        = "splunk"
  token       = var.splunk_token
  host_port   = var.splunk_host_port
  source_type = "generic_single_line"
}
```

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
  instance_id    = cloudamqp_instance.instance.id
  name          = "stackdriver"
  project_id    = var.stackdriver_project_id
  private_key   = var.stackdriver_private_key
  client_email  = var.stackdriver_client_email
}
```

or by using google_service_account_key resource from Google provider

```hcl
resource "google_service_account" "service_account" {
  account_id    = "<account_id>"
  description   = "<description>"
  display_name  = "<display_name>"
}

resource "google_service_account_key" "service_account_key" {
  service_account_id = google_service_account.service_account.name
}

resource "cloudamqp_integration_log" "stackdriver" {
  instance_id   = cloudamqp_instance.instance.id
  name          = "stackdriver"
  project_id    = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).project_id
  private_key   = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).private_key
  client_email  = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).client_email
}
```

</details>

<details>
  <summary>
    <b>
      <i>Stackdriver log integration (v1.21.0 or newer versions)</i>
    </b>
  </summary>

Use credentials argument and let the provider do the Base64decode and internally populate,
*project_id, client_name, private_key*

```hcl
resource "google_service_account" "service_account" {
  account_id    = "<account_id>"
  description   = "<description>"
  display_name  = "<display_name>"
}

resource "google_service_account_key" "service_account_key" {
  service_account_id = google_service_account.service_account.name
}

resource "cloudamqp_integration_log" "stackdriver" {
  instance_id = cloudamqp_instance.instance.id
  name        = "stackdriver"
  credentials = google_service_account_key.service_account_key.private_key
}
```

or use the same as earlier version and decode the google service account key

```hcl
resource "google_service_account" "service_account" {
  account_id    = "<account_id>"
  description   = "<description>"
  display_name  = "<display_name>"
}

resource "google_service_account_key" "service_account_key" {
  service_account_id = google_service_account.service_account.name
}

resource "cloudamqp_integration_log" "stackdriver" {
  instance_id   = cloudamqp_instance.instance.id
  name          = "stackdriver"
  project_id    = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).project_id
  private_key   = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).private_key
  client_email  = jsondecode(base64decode(google_service_account_key.service_account_key.private_key)).client_email
}
```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id` -  (Required) Instance identifier for the CloudAMQP instance.

Valid arguments for each third party log integrations below. Corresponding API backend documentation can be
found here [CloudAMQP API add integration].

<details>
  <summary>
    <b>Azure monitoring</b>
  </summary>

The following arguments used by Azure monitoring.

* `name`               - (Required) The name of the third party log integration (`azure_monitoring`).
* `application_id`     - (Required) The application identifier.
* `application_secret` - (Required/Sensitive) The application secret.
* `dcr_id`             - (Required) ID of data collection rule that your DCE is linked to.
* `dce_uri`            - (Required) The data collection endpoint.
* `tenant_id`          - (Required) The tenant identifier.
* `table`              - (Required) The table name.

Use Azure portal to configure external access for Azure Monitor. [Tutorial to find/create all arguments]

</details>

<details>
  <summary>
    <b>Cloudwatch</b>
  </summary>

The following arguments used by CloudWatch.

* `name`              - (Required) The name of the third party log integration (`cloudwatchlog`).
* `access_key_id`     - (Required/Sensitive) AWS access key identifier.
* `secret_access_key` - (Required/Sensitive) AWS secret access key.
* `region`            - (Required) AWS region hosting the integration service.

Optional arguments introduced in version [v1.38.0].

* `retention` - (Optional) Number of days to retain log events in `CloudAMQP` log group.

  ***Note:*** Possible values are: 0 (never expire) or between 1-3653, read more about valid values in
  the [Cloudwatch Log retention].

* `tags` - (Optional) Enter tags to `CloudAMQP` log group like this: `Project=A,Environment=Development`.

  ***Note:*** Tags are only added, unwanted tags needs to be removed manually in the AWS console.
  Read more about tags format in the [Cloudwatch Log tags]

#### IAM permissions

Create an IAM user with programmatic access and the following permissions: `CreateLogGroup`, `CreateLogStream`, `DescribeLogGroups`, `DescribeLogStreams` and `PutLogEvents`.

Optional arguments requires IAM permission: `PutRetentionPolicy`, `DeleteRetentionPolicy` and `TagResource`.

</details>

<details>
  <summary>
    <b>Coralogix</b>
  </summary>

The following arguments used by Coralogix.

* `name`        - (Required) The name of the third party log integration (`coralogix`).
* `application` - (Required) The application name for Coralogix.
* `endpoint`    - (Required) The syslog destination to send the logs to for Coralogix.
  - `syslog.eu1.coralogix.com:6514` (Europe - Ireland)
  - `syslog.eu2.coralogix.com:6514` (Europe - Stockholm)
  - `syslog.ap1.coralogix.com:6514` (Asia Pacific - Mumbai)
  - `syslog.ap2.coralogix.com:6514` (Asia Pacific - Singapore)
  - `syslog.ap3.coralogix.com:6514` (Asia Pacific - Sydney)
  - `syslog.us1.coralogix.com:6514` (US - Ohio)
  - `syslog.us2.coralogix.com:6514` (US - Oregon)
* `private_key` - (Required/Sensitive) The private access key.
* `subsystem`   - (Required) The subsystem name for Coralogix.

Create a 'Send-Your-Data' private API key, [Coralogix documentation]

~> ***Important:*** As of December 12, 2025, Coralogix has deprecated legacy endpoints. If you're using an old endpoint (e.g., `syslog.coralogix.com`, `syslog.coralogix.us`, `syslog.coralogix.in`, `syslog.cx498.coralogix.com`, or `syslog.coralogixsg.com`), you must migrate to the appropriate regional endpoint. See the [Coralogix endpoint deprecation notice] for the complete migration mapping.

Existing integrations created before this change will show a configuration drift in `terraform plan`. Update your configuration to use the new regional endpoint corresponding to your Coralogix region.

</details>

<details>
  <summary>
    <b>Datadog</b>
  </summary>

The following arguments used by Data dog.

* `name`    - (Required) The name of the third party log integration (`datadog`).
* `api_key` - (Required/Sensitive) The API key.

  ***Note:*** Create a Datadog API key at, [app.datadoghq.com]

* `region`  - (Required) Region hosting the integration service. Valid regions, `us1`, `us3`, `us5`
              and `eu`.

Optional arguments:

* `tags` - (Optional) Tags. e.g. `env=prod,region=europe`.

  ***Note:*** If tags are used with Datadog. The value part (prod, europe, ...) must start with a
              letter, read more about tags format in the [Datadog documentation].

</details>

<details>
  <summary>
    <b>Log Entries</b>
  </summary>

The following arguments used by Log entries.

* `name`  - (Required) The name of the third party log integration (`logentries`).
* `token` - (Required/Sensitive) Token used for authentication.

Create a Logentries token at [logentries add-log]

</details>

<details>
  <summary>
    <b>Loggly</b>
  </summary>

The following arguments used by Loggly.

* `name`  - (Required) The name of the third party log integration (`loggly`).
* `token` - (Required/Sensitive) Token used for authentication.

Create a Loggly token at `https://{your-company}.loggly.com/tokens`

</details>

<details>
  <summary>
    <b>Papertrail</b>
  </summary>

The following arguments used by Papertrail.

* `name` - (Required) The name of the third party log integration (`papertrail`).
* `url`  - (Required) Endpoint to log integration.

Create a Papertrail endpoint at [papertrail setup]

</details>

<details>
  <summary>
    <b>Scalyr</b>
  </summary>

The following arguments used by Scalyr.

* `name`  - (Required) The name of the third party log integration (`scalyr`).
* `token` - (Required/Sensitive) Token used for authentication.
* `host`  - (Required) The host for Scalyr integration. Valid hosts `app.scalyr.com` and
            `app.eu.scalyr.com`

Create a Log write token at [Scalyr keys]

</details>

<details>
  <summary>
    <b>Splunk</b>
  </summary>

The following arguments used by Splunk.

* `name`       - (Required) The name of the third party log integration (`splunk`).
* `token`      - (Required/Sensitive) Token used for authentication.
* `host_port`  - (Required) Destination to send the logs.
* `sourcetype` - (Required) Assign source type to the data exported, eg. generic_single_line.

Create a HTTP Event Collector token at `https://<your-splunk>.cloud.splunk.com/en-US/manager/search/http-eventcollector`

</details>

<details>
  <summary>
    <b>Stackdriver</b>
  </summary>

The following arguments used by Stackdriver.

* `name`        - (Required) The name of the third party log integration (`stackdriver`).
* `credentials` - (Required/Sensitive) Google Service Account private key credentials.

  ***Note:*** The service Account needs to have `log writer` role added.

Optional arguments for older provider versions.

* `project_id`   - (Optional/Computed) The project identifier.
* `private_key`  - (Optional/Computed/Sensitive) The private access key.
* `client_email` - (Optional/Computed) The client email registered for the integration service.

</details>

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_integration_log`can be imported using the resource identifier together with CloudAMQP
instance identifier. The identifiers are CSV separated, see example below. To retrieve the resource,
use [CloudAMQP API list integration].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_integration_log.this
  id = format("<id>,%s", cloudamqp_instance.instance.id)
}
```

`terraform import cloudamqp_integration_log.this <id>,<instance_id>`

[v1.38.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.38.0
[CloudAMQP API add integration]: https://docs.cloudamqp.com/instance-api.html#tag/integrations/post/integrations/logs/{system}
[Tutorial to find/create all arguments]: https://learn.microsoft.com/en-us/azure/azure-monitor/logs/tutorial-logs-ingestion-portal
[Cloudwatch Log retention]: https://docs.aws.amazon.com/AmazonCloudWatchLogs/latest/APIReference/API_PutRetentionPolicy.html#API_PutRetentionPolicy_RequestSyntax
[Cloudwatch Log tags]: https://docs.aws.amazon.com/AmazonCloudWatchLogs/latest/APIReference/API_TagLogGroup.html#API_TagLogGroup_RequestSyntax
[Coralogix documentation]: https://coralogix.com/docs/send-your-data-api-key/
[app.datadoghq.com]: https://app.datadoghq.com/
[Datadog documentation]: https://docs.datadoghq.com/getting_started/tagging/#define-tags
[logentries add-log]: https://logentries.com/app#/add-log/manual
[CloudAMQP API list integration]: https://docs.cloudamqp.com/instance-api.html#tag/integrations/get/integrations/logs

---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_integration_metric_prometheus"
description: |-
  Creates and manages third party Prometheus metrics integrations for a CloudAMQP instance.
---

<!-- markdownlint-disable MD024 -->
<!-- markdownlint-disable MD033 -->

# cloudamqp_integration_metric_prometheus

This resource allows you to create and manage Prometheus-compatible metric integrations for CloudAMQP instances. Currently supported integrations include New Relic v3, Datadog v3, Azure Monitor, Splunk v2, Dynatrace, CloudWatch v3, Stackdriver v2, Grafana Cloud (Mimir), and Prometheus Remote Write.

## Example Usage

### New Relic v3

```hcl
resource "cloudamqp_integration_metric_prometheus" "newrelic_v3" {
  instance_id = cloudamqp_instance.instance.id

  newrelic_v3 {
    api_key = var.newrelic_api_key
    region  = "us"
    tags    = "key=value,key2=value2"
  }
}
```

### Datadog v3

```hcl
resource "cloudamqp_integration_metric_prometheus" "datadog_v3" {
  instance_id = cloudamqp_instance.instance.id

  datadog_v3 {
    api_key                            = var.datadog_api_key
    region                             = "us1"
    tags                               = "key=value,key2=value2"
    rabbitmq_dashboard_metrics_format  = true
  }
}
```

### Azure Monitor

```hcl
resource "cloudamqp_integration_metric_prometheus" "azure_monitor" {
  instance_id = cloudamqp_instance.instance.id

  azure_monitor {
    connection_string = var.azure_monitor_connection_string
  }
}
```

### Splunk v2

```hcl
resource "cloudamqp_integration_metric_prometheus" "splunk_v2" {
  instance_id = cloudamqp_instance.instance.id

  splunk_v2 {
    token    = var.splunk_token
    endpoint = var.splunk_endpoint
    tags     = "key=value,key2=value2"
  }
}
```

### Dynatrace

```hcl
resource "cloudamqp_integration_metric_prometheus" "dynatrace" {
  instance_id = cloudamqp_instance.instance.id

  dynatrace {
    environment_id = var.dynatrace_environment_id
    access_token   = var.dynatrace_access_token
    tags           = "key=value,key2=value2"
  }
}
```

### CloudWatch v3

```hcl
resource "cloudamqp_integration_metric_prometheus" "cloudwatch_v3" {
  instance_id = cloudamqp_instance.instance.id

  cloudwatch_v3 {
    iam_role        = var.cloudwatch_iam_role
    iam_external_id = var.cloudwatch_iam_external_id
    region          = var.cloudwatch_region
    tags            = "key=value,key2=value2"
  }
}
```

### Stackdriver v2

```hcl
resource "cloudamqp_integration_metric_prometheus" "stackdriver_v2" {
  instance_id = cloudamqp_instance.instance.id

  stackdriver_v2 {
    credentials_file = var.google_service_account_key
    tags             = "key=value,key2=value2"
  }
}
```

**Note:** The `credentials_file` should contain a Base64-encoded Google service account key JSON file. You can create a service account in Google Cloud Console with the "Monitoring Metric Writer" role and download the key file. Then encode it with:

```bash
base64 -i /path/to/service-account-key.json
```

### Grafana Cloud (Mimir)

```hcl
resource "cloudamqp_integration_metric_prometheus" "grafana" {
  instance_id = cloudamqp_instance.instance.id

  grafana {
    endpoint    = var.grafana_endpoint
    instance_id = var.grafana_instance_id
    api_token   = var.grafana_api_token
    tags        = "key=value,key2=value2"
  }
}
```

**Note:** Find all three values in the Grafana Cloud Portal, under **Send Metrics** on the **Prometheus** tile of your stack. The `instance_id` is the numeric Prometheus instance identifier, which is not the same as the Loki instance identifier used by the Grafana Cloud log integration.

### Prometheus Remote Write

```hcl
resource "cloudamqp_integration_metric_prometheus" "prometheus_remote_write" {
  instance_id = cloudamqp_instance.instance.id

  prometheus_remote_write {
    endpoint  = var.remote_write_endpoint
    auth_type = "basic_auth"
    username  = var.remote_write_username
    password  = var.remote_write_password
    headers   = "X-Scope-OrgID: my-tenant"
    tags      = "key=value,key2=value2"
  }
}
```

**Note:** Headers are sent with every request whichever `auth_type` you pick, so a tenant header can accompany basic auth. Multi-tenant Mimir and Cortex read `X-Scope-OrgID`, Thanos reads `THANOS-TENANT`. Set `auth_type` to `headers` when the credentials themselves live in a header, for example `Authorization: Bearer your-token`.

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Instance identifier for the CloudAMQP instance.
* `metrics_filter` - (Optional) List of metrics to include in the integration. If not specified, default metrics are included.
  For more information about metrics filtering, see the [metrics filtering documentation](https://www.cloudamqp.com/docs/monitoring_metrics_splunk_v2.html#metrics-filtering).

Exactly one of the following integration blocks must be specified:

### newrelic_v3

The following arguments are supported:

* `api_key` - (Required) New Relic API key for authentication.
* `region` - (Required) New Relic region code. Valid values: `eu`, `us`.
* `tags` - (Optional) Additional tags to attach to metrics. Format: `key=value,key2=value2`.

### datadog_v3

The following arguments are supported:

* `api_key` - (Required) Datadog API key for authentication.
* `region` - (Required) Datadog region code. Valid values: `us1`, `us3`, `us5`, `eu1`, `ap2`.
* `tags` - (Optional) Additional tags to attach to metrics. Format: `key=value,key2=value2`.
* `rabbitmq_dashboard_metrics_format` - (Optional) Enable metric name transformation to match Datadog's RabbitMQ dashboard format. Default: `false`. **Note:** This option is only available for RabbitMQ clusters, not LavinMQ clusters.

### azure_monitor

The following arguments are supported:

* `connection_string` - (Required) Azure Application Insights Connection String for authentication.

### splunk_v2

The following arguments are supported:

* `token` - (Required) Splunk HEC (HTTP Event Collector) token for authentication.
* `endpoint` - (Required) Splunk HEC endpoint URL. Example: `https://your-instance-id.splunkcloud.com:8088/services/collector`.
* `tags` - (Optional) Additional tags to attach to metrics. Format: `key=value,key2=value2`.

### dynatrace

The following arguments are supported:

* `environment_id` - (Required) Dynatrace environment ID.
* `access_token` - (Required) Dynatrace access token with 'Ingest metrics' permission.
* `tags` - (Optional) Additional tags to attach to metrics. Format: `key=value,key2=value2`.

### cloudwatch_v3

The following arguments are supported:

* `iam_role` - (Required) AWS IAM role ARN with PutMetricData permission for CloudWatch integration.
* `iam_external_id` - (Required) AWS IAM external ID for role assumption.
* `region` - (Required) AWS region for CloudWatch metrics.
* `tags` - (Optional) Additional tags to attach to metrics. Format: `key=value,key2=value2`.

### stackdriver_v2

The following arguments are supported:

* `credentials_file` - (Required) Base64-encoded Google service account key JSON file with 'Monitoring Metric Writer' permission.
* `tags` - (Optional) Additional tags to attach to metrics. Format: `key=value,key2=value2`.

The following computed attributes are available:

* `project_id` - Google Cloud project ID (extracted from credentials file).
* `client_email` - Google service account client email (extracted from credentials file).
* `private_key` - Google service account private key (extracted from credentials file).
* `private_key_id` - Google service account private key ID (extracted from credentials file).

### grafana

The following arguments are supported:

* `endpoint` - (Required) Grafana Cloud Prometheus remote write endpoint. Example: `https://prometheus-prod-01-eu-west-0.grafana.net/api/prom/push`.
* `instance_id` - (Required) Grafana Cloud numeric Prometheus instance identifier, used as the basic auth username.
* `api_token` - (Required) Grafana Cloud API token with the `metrics:write` scope, or a service account token with the `MetricsPublisher` role.
* `tags` - (Optional) Additional tags to attach to metrics. Format: `key=value,key2=value2`.

### prometheus_remote_write

The following arguments are supported:

* `endpoint` - (Required) Remote write endpoint including the path, over HTTPS. Example: `https://mimir.example.com/api/v1/push`.
* `auth_type` - (Optional) Authentication for the endpoint. Valid values: `none`, `basic_auth`, `headers`. Default: `none`.
* `username` - (Optional) Username, used when `auth_type` is `basic_auth`.
* `password` - (Optional) Password or token, used when `auth_type` is `basic_auth`.
* `headers` - (Optional) Headers sent with every request, one `key: value` pair per line. Required when `auth_type` is `headers`.
* `tags` - (Optional) Additional tags to attach to metrics. Format: `key=value,key2=value2`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The integration identifier.

For `stackdriver_v2` integrations, the following computed attributes are also available:

* `stackdriver_v2.0.project_id` - Google Cloud project ID extracted from the credentials file.
* `stackdriver_v2.0.client_email` - Google service account client email extracted from the credentials file.
* `stackdriver_v2.0.private_key` - Google service account private key extracted from the credentials file.
* `stackdriver_v2.0.private_key_id` - Google service account private key ID extracted from the credentials file.

## Import

CloudAMQP Prometheus metric integrations can be imported using the integration identifier together with the instance identifier. The import identifier should be in the format `{integration_id},{instance_id}`.

From Terraform v1.5.0, the `import` block can be used to import this resource:

### New Relic v3

```hcl
import {
  to = cloudamqp_integration_metric_prometheus.newrelic_v3
  id = format("<integration_id>,%s", cloudamqp_instance.instance.id)
}
```

### Datadog v3

```hcl
import {
  to = cloudamqp_integration_metric_prometheus.datadog_v3
  id = format("<integration_id>,%s", cloudamqp_instance.instance.id)
}
```

### Azure Monitor

```hcl
import {
  to = cloudamqp_integration_metric_prometheus.azure_monitor
  id = format("<integration_id>,%s", cloudamqp_instance.instance.id)
}
```

### Splunk v2

```hcl
import {
  to = cloudamqp_integration_metric_prometheus.splunk_v2
  id = format("<integration_id>,%s", cloudamqp_instance.instance.id)
}
```

### Dynatrace

```hcl
import {
  to = cloudamqp_integration_metric_prometheus.dynatrace
  id = format("<integration_id>,%s", cloudamqp_instance.instance.id)
}
```

### CloudWatch v3

```hcl
import {
  to = cloudamqp_integration_metric_prometheus.cloudwatch_v3
  id = format("<integration_id>,%s", cloudamqp_instance.instance.id)
}
```

### Stackdriver v2

```hcl
import {
  to = cloudamqp_integration_metric_prometheus.stackdriver_v2
  id = format("<integration_id>,%s", cloudamqp_instance.instance.id)
}
```

### Grafana Cloud (Mimir)

```hcl
import {
  to = cloudamqp_integration_metric_prometheus.grafana
  id = format("<integration_id>,%s", cloudamqp_instance.instance.id)
}
```

### Prometheus Remote Write

```hcl
import {
  to = cloudamqp_integration_metric_prometheus.prometheus_remote_write
  id = format("<integration_id>,%s", cloudamqp_instance.instance.id)
}
```

Or use Terraform CLI:

```sh
terraform import cloudamqp_integration_metric_prometheus.newrelic_v3 <integration_id>,<instance_id>
terraform import cloudamqp_integration_metric_prometheus.datadog_v3 <integration_id>,<instance_id>
terraform import cloudamqp_integration_metric_prometheus.azure_monitor <integration_id>,<instance_id>
terraform import cloudamqp_integration_metric_prometheus.splunk_v2 <integration_id>,<instance_id>
terraform import cloudamqp_integration_metric_prometheus.dynatrace <integration_id>,<instance_id>
terraform import cloudamqp_integration_metric_prometheus.cloudwatch_v3 <integration_id>,<instance_id>
terraform import cloudamqp_integration_metric_prometheus.stackdriver_v2 <integration_id>,<instance_id>
terraform import cloudamqp_integration_metric_prometheus.grafana <integration_id>,<instance_id>
terraform import cloudamqp_integration_metric_prometheus.prometheus_remote_write <integration_id>,<instance_id>
```

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

# cloudamqp_integration_metric_prometheus

This resource allows you to create and manage Prometheus-compatible metric integrations for CloudAMQP instances. Currently supported integrations include New Relic v3, Datadog v3, Azure Monitor, and Splunk v2.

## Example Usage

### New Relic v3

```hcl
resource "cloudamqp_integration_metric_prometheus" "newrelic_v3" {
  instance_id = cloudamqp_instance.instance.id

  newrelic_v3 {
    api_key = var.newrelic_api_key
    tags    = "key=value,key2=value2"
  }
}
```

### Datadog v3

```hcl
resource "cloudamqp_integration_metric_prometheus" "datadog_v3" {
  instance_id = cloudamqp_instance.instance.id

  datadog_v3 {
    api_key = var.datadog_api_key
    region  = "us1"
    tags    = "key=value,key2=value2"
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

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Instance identifier for the CloudAMQP instance.

Exactly one of the following integration blocks must be specified:

### newrelic_v3

The following arguments are supported:

* `api_key` - (Required) New Relic API key for authentication.
* `tags` - (Optional) Additional tags to attach to metrics. Format: `key=value,key2=value2`.

### datadog_v3

The following arguments are supported:

* `api_key` - (Required) Datadog API key for authentication.
* `region` - (Required) Datadog region code. Valid values: `us1`, `us3`, `us5`, `eu1`.
* `tags` - (Optional) Additional tags to attach to metrics. Format: `key=value,key2=value2`.

### azure_monitor

The following arguments are supported:

* `connection_string` - (Required) Azure Application Insights Connection String for authentication.

### splunk_v2

The following arguments are supported:

* `token` - (Required) Splunk HEC (HTTP Event Collector) token for authentication.
* `endpoint` - (Required) Splunk HEC endpoint URL. Example: `https://your-instance-id.splunkcloud.com:8088/services/collector`.
* `tags` - (Optional) Additional tags to attach to metrics. Format: `key=value,key2=value2`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The integration identifier.

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

Or use Terraform CLI:

```
$ terraform import cloudamqp_integration_metric_prometheus.newrelic_v3 <integration_id>,<instance_id>
$ terraform import cloudamqp_integration_metric_prometheus.datadog_v3 <integration_id>,<instance_id>
$ terraform import cloudamqp_integration_metric_prometheus.azure_monitor <integration_id>,<instance_id>
$ terraform import cloudamqp_integration_metric_prometheus.splunk_v2 <integration_id>,<instance_id>
```

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

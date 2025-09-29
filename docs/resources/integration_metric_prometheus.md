# cloudamqp_integration_metric_prometheus

This resource allows you to create and manage Prometheus-compatible metric integrations for CloudAMQP instances. Currently supported integrations include New Relic v3 and Datadog v3.

## Example Usage

### New Relic v3 Integration

```hcl
resource "cloudamqp_integration_metric_prometheus" "newrelic" {
  instance_id = cloudamqp_instance.instance.id

  newrelic_v3 {
    api_key = var.newrelic_api_key
    tags    = "env=prod,region=us-east-1"
  }
}
```

### Datadog v3 Integration

```hcl
resource "cloudamqp_integration_metric_prometheus" "datadog" {
  instance_id = cloudamqp_instance.instance.id

  datadog_v3 {
    api_key = var.datadog_api_key
    region  = "us1"
    tags    = "env=prod,region=us-east-1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance identifier.
* `newrelic_v3` - (Optional) Configuration block for New Relic v3 integration. Cannot be used with `datadog_v3`.
* `datadog_v3` - (Optional) Configuration block for Datadog v3 integration. Cannot be used with `newrelic_v3`.

### newrelic_v3 Block

The following arguments are supported:

* `api_key` - (Required) New Relic API key for authentication.
* `tags` - (Optional) Additional tags to attach to metrics. Format: `key=value,key2=value2`.

### datadog_v3 Block

The following arguments are supported:

* `api_key` - (Required) Datadog API key for authentication.
* `region` - (Optional) Datadog region code. Defaults to `us1`. Valid values: `us1`, `us3`, `us5`, `eu1`.
* `tags` - (Optional) Additional tags to attach to metrics. Format: `key=value,key2=value2`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The integration identifier.

## Import

CloudAMQP Prometheus metric integrations can be imported using the integration identifier together with the instance identifier. The import identifier should be in the format `{integration_id},{instance_id}`.

```
$ terraform import cloudamqp_integration_metric_prometheus.datadog 12345,67890
```

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

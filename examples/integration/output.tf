// Instance output
output "instance_id" {
  value = cloudamqp_instance.instance.id
}

output "instance_name" {
  value = cloudamqp_instance.instance.name
}

output "instance_plan" {
  value = cloudamqp_instance.instance.plan
}

output "instance_region" {
  value = cloudamqp_instance.instance.region
}

// Integration

output "cloudwatch_log_id" {
  value = cloudamqp_integration_log.cloudwatch.id
}

output "cloudwatch_metric_id" {
  value = cloudamqp_integration_metric.cloudwatch.id
}

// Prometheus integrations
output "newrelic_v3_integration_id" {
  value = cloudamqp_integration_metric_prometheus.newrelic_v3.id
}

output "datadog_v3_integration_id" {
  value = cloudamqp_integration_metric_prometheus.datadog_v3.id
}

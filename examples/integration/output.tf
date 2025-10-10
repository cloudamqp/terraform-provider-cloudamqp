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

output "cloudwatch_log_id" {
  value = cloudamqp_integration_log.cloudwatchlog.id
}

output "logentries_log_id" {
  value = cloudamqp_integration_log.logentries.id
}

output "loggly_log_id" {
  value = cloudamqp_integration_log.loggly.id
}

output "papertrail_log_id" {
  value = cloudamqp_integration_log.papertrail.id
}

output "splunk_log_id" {
  value = cloudamqp_integration_log.splunk.id
}

output "cloudwatch_metric_id" {
  value = cloudamqp_integration_metric.cloudwatch.id
}

output "newrelic_v3_integration_id" {
  value = cloudamqp_integration_metric_prometheus.newrelic_v3.id
}

output "datadog_v3_integration_id" {
  value = cloudamqp_integration_metric_prometheus.datadog_v3.id
}

output "splunk_v2_integration_id" {
  value = cloudamqp_integration_metric_prometheus.splunk_v2.id
}

output "dynatrace_integration_id" {
  value = cloudamqp_integration_metric_prometheus.dynatrace.id
}

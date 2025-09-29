terraform {
  required_providers {
    cloudamqp = {
      source  = "cloudamqp/cloudamqp"
      version = "~>1.0"
    }
  }
}

provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

resource "cloudamqp_instance" "instance" {
  name   = "terraform-integration-test"
  plan   = "bunny-1"
  region = "amazon-web-services::us-east-1"
}

// LOG INTEGRATION
resource "cloudamqp_integration_log" "cloudwatchlog" {
  instance_id       = cloudamqp_instance.instance.id
  name              = "cloudwatchlog"
  access_key_id     = var.aws_access_key
  secret_access_key = var.aws_secret_key
  region            = var.aws_region
}

resource "cloudamqp_integration_log" "logentries" {
  instance_id = cloudamqp_instance.instance.id
  name        = "logentries"
  token       = var.logentries_token
}

resource "cloudamqp_integration_log" "loggly" {
  instance_id = cloudamqp_instance.instance.id
  name        = "loggly"
  token       = var.loggly_token
}

resource "cloudamqp_integration_log" "papertrail" {
  instance_id = cloudamqp_instance.instance.id
  name        = "papertrail"
  url         = var.papertrail_url
}

resource "cloudamqp_integration_log" "splunk" {
  instance_id = cloudamqp_instance.instance.id
  name        = "splunk"
  token       = var.splunk_token
  host_port   = var.splunk_host_port
}

// METRIC INTEGRATION
resource "cloudamqp_integration_metric" "cloudwatch" {
  instance_id       = cloudamqp_instance.instance.id
  name              = "cloudwatch"
  access_key_id     = var.aws_access_key
  secret_access_key = var.aws_secret_key
  region            = var.aws_region
}

resource "cloudamqp_integration_metric" "datadog" {
  instance_id = cloudamqp_instance.instance.id
  name        = "datadog_v2"
  region      = var.datadog_region
  api_key     = var.datadog_apikey
}

resource "cloudamqp_integration_metric" "librato" {
  instance_id = cloudamqp_instance.instance.id
  name        = "librato"
  email       = "integration@example.com"
  api_key     = var.librato_apikey
}

resource "cloudamqp_integration_metric" "newrelic_v2" {
  instance_id = cloudamqp_instance.instance.id
  name        = "newrelic_v2"
  region      = var.newrelic_region
  api_key     = var.newrelic_apikey
}

resource "cloudamqp_integration_metric_prometheus" "newrelic_v3" {
  instance_id = cloudamqp_instance.instance.id
  name        = "newrelic_v3"
  api_key     = var.newrelic_apikey
}

resource "cloudamqp_integration_metric_prometheus" "datadog_v3" {
  instance_id = cloudamqp_instance.instance.id
  name        = "datadog_v3"
  region      = var.datadog_region
  api_key     = var.datadog_apikey
}

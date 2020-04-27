provider "cloudamqp" {
  apikey = "<cloudamqp_apikey>"
}

resource "cloudamqp_instance" "instance" {
  name 				= "terraform-integration-test"
  nodes 			= 1
  plan  			= "bunny"
  region 			= "amazon-web-services::us-east-1"
  rmq_version = "3.8.2"
  tags 				= ["terraform"]
  vpc_subnet = "192.168.0.1/24"
}

// LOG INTEGRATION
resource "cloudamqp_integration_log" "cloudwatchlog" {
  instance_id = cloudamqp_instance.instance.id
  name = "cloudwatchlog"
  access_key_id = "<aws_access_key_id>"
  secret_access_key = "<aws_secret_access_key>"
  region = "us-east-1"
}

resource "cloudamqp_integration_log" "logentries" {
  instance_id = cloudamqp_instance.instance.id
  name = "logentries"
  token = "<token>"
}

resource "cloudamqp_integration_log" "loggly" {
  instance_id = cloudamqp_instance.instance.id
  name = "loggly"
  token = "<token>"
}

resource "cloudamqp_integration_log" "papertrail" {
  instance_id = cloudamqp_instance.instance.id
  name = "papertrail"
  url = "<url>"
}

resource "cloudamqp_integration_log" "splunk" {
  instance_id = cloudamqp_instance.instance.id
  name = "splunk"
  token = "<token>"
  host_port = "<host_port>"
}

// METRIC INTEGRATION
resource "cloudamqp_integration_metric" "cloudwatch" {
  instance_id = cloudamqp_instance.instance.id
  name = "cloudwatch"
  access_key_id = "<aws_access_key_id>"
  secret_access_key = "<aws_secret_access_key>"
  region = "us-east-1"
}

resource "cloudamqp_integration_metric" "datadog" {
  instance_id = cloudamqp_instance.instance.id
  name = "datadog_v2"
  region = "us"
  api_key = "<api_key>"
}

resource "cloudamqp_integration_metric" "librato" {
  instance_id = cloudamqp_instance.instance.id
  name = "librato"
  email = "integration@example.com"
  api_key = "<api_key>"
}

resource "cloudamqp_integration_metric" "newrelic_v2" {
  instance_id = cloudamqp_instance.instance.id
  name = "newrelic_v2"
  region = "us"
  api_key = "<api_key>"
}

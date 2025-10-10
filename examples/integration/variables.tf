// CloudAMQP
variable "cloudamqp_customer_api_key" {
  type = string
}

// AWS
variable "aws_access_key" {
  type = string
}
variable "aws_secret_key" {
  type = string
}
variable "aws_region" {
  type = string
}

// Logentries
variable "logentries_token" {
  type = string
}

// Loggly
variable "loggly_token" {
  type = string
}

// Papertrail
variable "papertrail_url" {
  type = string
}

// Splunk
variable "splunk_host_port" {
  type = string
}

variable "splunk_endpoint" {
  type = string
}

variable "splunk_token" {
  type = string
}

// Datadog
variable "datadog_apikey" {
  type = string
}

variable "datadog_region" {
  type = string

  validation {
    condition     = var.datadog_region == "us1" || var.datadog_region == "us3" || var.datadog_region == "us5" || var.datadog_region == "eu"
    error_message = "Available regions are, us1, us3, us5 and eu"
  }
}

// Librato
variable "librato_apikey" {
  type = string
}

// New Relic
variable "newrelic_apikey" {
  type = string
}

variable "newrelic_region" {
  type = string

  validation {
    condition     = var.newrelic_region == "us" || var.newrelic_region == "eu"
    error_message = "Available regions are, us and eu"
  }
}

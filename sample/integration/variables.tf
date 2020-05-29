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

// Datadog
variable "datadog_apikey" {
  type = string
}

variable "datadog_region" {
  type = string

  validation {
    condition     = var.newrelic_region == "us" || var.newrelic_region == "eu"
    error_message = "Available regions are, us and eu"
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

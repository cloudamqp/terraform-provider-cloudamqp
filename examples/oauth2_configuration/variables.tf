variable "cloudamqp_customer_api_key" {
  description = "CloudAMQP customer API key"
  type        = string
  sensitive   = true
}

variable "instance_id" {
  description = "CloudAMQP instance ID"
  type        = number
}

variable "cloudamqp_baseurl" {
  description = "CloudAMQP base URL"
  type        = string
}

variable "issuer_url" {
  description = "Issuer URL"
  type        = string
}

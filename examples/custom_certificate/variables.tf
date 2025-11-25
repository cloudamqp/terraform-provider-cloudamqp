variable "cloudamqp_customer_api_key" {
  description = "CloudAMQP customer API key"
  type        = string
  sensitive   = true
}

variable "custom_ca_path" {
  description = "Path to custom PEM-encoded Certificate Authority (CA)"
  type        = string
}

variable "custom_cert_path" {
  description = "Path to custom PEM-encoded server certificate"
  type        = string
}

variable "custom_private_key_path" {
  description = "Path to custom PEM-encoded private key"
  type        = string
}
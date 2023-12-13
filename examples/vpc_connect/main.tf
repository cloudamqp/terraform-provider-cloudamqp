terraform {
  required_providers {
    cloudamqp = {
      source = "cloudamqp/cloudamqp"
      version = "~> 1.0"
    }
  }
}

provider "cloudamqp" {
  apikey  = var.cloudamqp_customer_api_key
}

// === AWS ===
resource "cloudamqp_vpc" "aws-vpc" {
  name    = "aws-instance"
  subnet  = "10.56.72.0/24"
  region  = "amazon-web-services::us-east-1"
  tags    = ["aws"]
}

resource "cloudamqp_instance" "aws-instance" {
  name    = "aws-instance"
  plan    = "bunny-1"
  region  = "amazon-web-services::us-east-1"
  tags    = ["aws"]
  vpc_id  = cloudamqp_vpc.aws-vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_vpc_connect" "aws-vpc-connect" {
  instance_id = cloudamqp_instance.aws-instance.id
  region  = cloudamqp_instance.aws-instance.region
  allowed_principals = [
    "arn:aws:iam::<AWS-ACCOUNT-ID>:user/<username>"
  ]
}

# // === AZURE ===
resource "cloudamqp_vpc" "azure-vpc" {
  name    = "Azure-instance"
  subnet  = "10.56.72.0/24"
  region  = "azure-arm::eastus"
  tags    = ["azure"]
}

resource "cloudamqp_instance" "azure-instance" {
  name    = "Azure-instance"
  plan    = "bunny-1"
  region  = "azure-arm::eastus"
  tags    = ["azure"]
  vpc_id  = cloudamqp_vpc.azure-vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_vpc_connect" "azure-vpc-connect" {
  instance_id = cloudamqp_instance.azure-instance.id
  region  = cloudamqp_instance.azure-instance.region
  approved_subscriptions = [
    "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  ]
}

# // === GOOGLE ===
resource "cloudamqp_vpc" "gcp-vpc" {
  name    = "gcp-instance"
  subnet  = "10.56.72.0/24"
  region  = "google-compute-engine::us-east1"
  tags    = ["gcp"]
}

resource "cloudamqp_instance" "gcp-instance" {
  name    = "gcp-instance"
  plan    = "bunny-1"
  region  = "google-compute-engine::us-east1"
  tags    = ["gcp"]
  vpc_id  = cloudamqp_vpc.gcp-vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_vpc_connect" "gcp-vpc-connect" {
  instance_id = cloudamqp_instance.gcp-instance.id
  region  = cloudamqp_instance.gcp-instance.region
  allowed_projects = [
    "<project-id>"
  ]
}
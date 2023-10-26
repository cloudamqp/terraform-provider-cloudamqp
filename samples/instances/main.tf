terraform {
  required_providers {
    cloudamqp = {
      source = "cloudamqp/cloudamqp"
      version = "~>1.0"
    }
  }
}

provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

// === Basic cloudamqp instance ===
resource "cloudamqp_instance" "instance_01" {
  name 				= "terraform-test-01"
  plan  			= "bunny-1"
  region 			= "amazon-web-services::us-east-1"
}

// === Basic cloudamqp instance with specific RabbitMQ version ===
resource "cloudamqp_instance" "instance_02" {
  name 				= "terraform-test-02"
  plan  			= "bunny-1"
  region 			= "amazon-web-services::us-east-1"
  rmq_version = "3.12.6"
}

// === Standalone VPC and with multiple cloudamqp instances ===
resource "cloudamqp_vpc" "vpc" {
  name        = "terraform-vpc"
  region 			= "amazon-web-services::us-east-1"
  subnet      = ["10.56.72.0/24"]
  tags        = ["aws"]
}

resource "cloudamqp_instance" "instance_03" {
  name 				        = "terraform-test-03"
  plan  			        = "bunny-1"
  region 			        = "amazon-web-services::us-east-1"
  tags                = ["aws"]
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_instance" "instance_04" {
  name 				        = "terraform-test-04"
  plan  			        = "bunny-1"
  region 			        = "amazon-web-services::us-east-1"
  tags                = ["aws"]
  vpc_id              = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}
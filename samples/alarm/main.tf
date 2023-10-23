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

// Instance
resource "cloudamqp_instance" "instance" {
  name   = "terraform-provider-test"
  plan   = "bunny-1"
  region = "amazon-web-services::us-east-1"
  no_default_alarms = true
}

// Notification and recipient
// Each instance will get one default recipient.
// The default recipient can either be imported as a resource or loaded as a data source.
data "cloudamqp_notification" "default_recipient" {
  instance_id = cloudamqp_instance.instance.id
  name = "Default"
}

resource "cloudamqp_notification" "alarm_recipient" {
  instance_id = cloudamqp_instance.instance.id
  type        = "email"
  value       = "alarm@example.com"
}

resource "cloudamqp_notification" "notification_recipient" {
  instance_id = cloudamqp_instance.instance.id
  type        = "email"
  value       = "notification@example.com"
}

// Default alarms
// Each instance will get a set of default alarms (cpu, memory and disk) upon creation.
// The default alarms can either be imported as resources or loaded as a data source.
// To disable creating default alarms when creating a new instance. Set `no_default_alarms`
// attribute to true in the instance resource.
data "cloudamqp_alarm" "default_cpu" {
  instance_id = cloudamqp_instance.instance.id
  type 				= "cpu"
}

data "cloudamqp_alarm" "default_memory" {
  instance_id = cloudamqp_instance.instance.id
  type 				= "memory"
}

data "cloudamqp_alarm" "default_disk" {
  instance_id = cloudamqp_instance.instance.id
  type 				= "disk"
}

// New alarms
resource "cloudamqp_alarm" "cpu_alarm" {
  instance_id     = cloudamqp_instance.instance.id
  type            = "cpu"
  value_threshold = 90
  time_threshold  = 600
  recipients      = [
    cloudamqp_notification.alarm_recipient.id,
    cloudamqp_notification.notification_recipient.id
  ]
}

resource "cloudamqp_alarm" "alarm_02" {
  instance_id     = cloudamqp_instance.instance.id
  type            = "memory"
  value_threshold = 90
  time_threshold  = 600
  recipients      = [
    cloudamqp_notification.alarm_recipient.id,
    cloudamqp_notification.notification_recipient.id
  ]
}

resource "cloudamqp_alarm" "disk_alarm" {
  instance_id     = cloudamqp_instance.instance.id
  type            = "disk"
  value_threshold = 10
  time_threshold  = 600
  recipients      = [cloudamqp_notification.alarm_recipient.id]
}

resource "cloudamqp_alarm" "queue_alarm" {
  instance_id     = cloudamqp_instance.instance.id
  type            = "queue"
  value_threshold = 120
  time_threshold  = 120
  enabled         = true
  queue_regex     = ".*"
  vhost_regex     = ".*"
  message_type    = "total"
  recipients      = [cloudamqp_notification.alarm_recipient.id]
}

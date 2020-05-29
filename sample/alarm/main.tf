provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

// Instance
resource "cloudamqp_instance" "rmq_bunny" {
  name   = "terraform-provider-test"
  plan   = "bunny"
  region = "amazon-web-services::us-east-1"
}

// Notification and recipient
// Each instance will get one default recipient,
// can either be imported as a resource or loaded as a data source.
data "cloudamqp_notificaion" "default_recipient" {
  instance_id = cloudamqp_instance.rmq_bunny.id
  name = "Default"
}

resource "cloudamqp_notification" "recipient_01" {
  instance_id = cloudamqp_instance.rmq_bunny.id
  type        = "email"
  value       = "alarm@example.com"
}

resource "cloudamqp_notification" "recipient_02" {
  instance_id = cloudamqp_instance.rmq_bunny.id
  type        = "email"
  value       = "notification@example.com"
}

// Alarm
// Each instance will get a set of alarms upon creation.
// Can either be imported as resources or loaded as a data source.
data "cloudamqp_alarm" "default_cpu" {
  instance_id = cloudamqp_instance.instance.id
  type 				= "cpu"
}

data "cloudamqp_alarm" "default_memory" {
  instance_id = cloudamqp_instance.instance.id
  type 				= "memory"
}

data "cloudamqp_alarm" "default_memory" {
  instance_id = cloudamqp_instance.instance.id
  type 				= "disk"
}

// New alarms
resource "cloudamqp_alarm" "alarm_01" {
  instance_id     = cloudamqp_instance.rmq_bunny.id
  type            = "cpu"
  value_threshold = 90
  time_threshold  = 600
  recipients      = [cloudamqp_notification.recipient_01.id, cloudamqp_notification.recipient_02]
}

resource "cloudamqp_alarm" "alarm_02" {
  instance_id     = cloudamqp_instance.rmq_bunny.id
  type            = "memory"
  value_threshold = 90
  time_threshold  = 600
  recipients      = [cloudamqp_notification.recipient_01.id, cloudamqp_notification.recipient_02]
}

resource "cloudamqp_alarm" "alarm_03" {
  instance_id     = cloudamqp_instance.rmq_bunny.id
  type            = "disk"
  value_threshold = 10
  time_threshold  = 600
  recipients      = [cloudamqp_notification.recipient_01.id]
}

resource "cloudamqp_alarm" "alarm_04" {
  instance_id     = cloudamqp_instance.rmq_bunny.id
  type            = "queue"
  value_threshold = 120
  time_threshold  = 120
  enabled         = true
  queue_regex     = ".*"
  vhost_regex     = ".*"
  message_type    = "total"
  recipients      = [cloudamqp_notification.recipient_01.id]
}

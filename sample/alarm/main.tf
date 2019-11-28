provider "cloudamqp" {
  apikey = "<apikey>"
}

// Instance
resource "cloudamqp_instance" "rmq_bunny" {
  name   = "terraform-provider-test"
  plan   = "bunny"
  region = "amazon-web-services::us-east-1"
}

// Notification and recipient
resource "cloudamqp_notification" "recipient_01" {
  instance_id = cloudamqp_instance.rmq_bunny.id
  type = "email"
  value = "alarm@example.com"
}

resource "cloudamqp_notification" "recipient_02" {
  instance_id = cloudamqp_instance.rmq_bunny.id
  type = "email"
  value = "notification@example.com"
}

// Alarm
resource "cloudamqp_alarm" "alarm_01" {
  instance_id = cloudamqp_instance.rmq_bunny.id
  type = "cpu"
  value_threshold = 90
  time_threshold = 600
  notifications = [cloudamqp_notification.recipient_01.id, cloudamqp_notification.recipient_02]
}

resource "cloudamqp_alarm" "alarm_02" {
  instance_id = cloudamqp_instance.rmq_bunny.id
  type = "memory"
  value_threshold = 90
  time_threshold = 600
  notifications = [cloudamqp_notification.recipient_01.id, cloudamqp_notification.recipient_02]
}

resource "cloudamqp_alarm" "alarm_03" {
  instance_id = cloudamqp_instance.rmq_bunny.id
  type = "disk"
  value_threshold = 10
  time_threshold = 600
  notifications = [cloudamqp_notification.recipient_01.id]
}

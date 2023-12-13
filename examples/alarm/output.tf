// Instance output
output "instance_id" {
  value = "${cloudamqp_instance.rmq_bunny.id}"
}

output "instance_name" {
  value = "${cloudamqp_instance.rmq_bunny.name}"
}

output "instance_plan" {
  value = "${cloudamqp_instance.rmq_bunny.plan}"
}

output "instance_region" {
  value = "${cloudamqp_instance.rmq_bunny.region}"
}

// Recipient/notification output
output "notification_type" {
  value = "${cloudamqp_notification.recipient_01.type}"
}

output "notification_value" {
  value = "${cloudamqp_notification.recipient_01.value}"
}

// Alarm output
output "alarm_id" {
  value = "${cloudamqp_alarm.alarm_01_cpu.id}"
}

output "alarm_type" {
  value = "${cloudamqp_alarm.alarm_01_cpu.type}"
}

output "alarm_value_threshold" {
  value = "${cloudamqp_alarm.alarm_01_cpu.value_threshold}"
}

output "alarm_time_threshold" {
  value = "${cloudamqp_alarm.alarm_01_cpu.time_threshold}"
}

output "alarm_notifications" {
  value = "${cloudamqp_alarm.alarm_01_cpu.notifications}"
}

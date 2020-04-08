// Instance output
output "instance_id" {
  value = cloudamqp_instance.instance.id
}

output "instance_name" {
  value = cloudamqp_instance.instance.name
}

output "instance_plan" {
  value = cloudamqp_instance.instance.plan
}

output "instance_region" {
  value = cloudamqp_instance.instance.region
}

// Plugin
output "plugin_name" {
  value = cloudamqp_plugin.mqtt_plugin.name
}

output "plugin_enabled" {
  value = cloudamqp_plugin.mqtt_plugin.enabled
}

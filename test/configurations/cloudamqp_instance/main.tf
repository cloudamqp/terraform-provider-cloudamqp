resource "cloudamqp_instance" "instance" {
	name = "{{.InstanceName}}"
	plan = "{{.InstancePlan}}"
	region = "amazon-web-services::us-east-1"
	tags = ["terraform"]  
}
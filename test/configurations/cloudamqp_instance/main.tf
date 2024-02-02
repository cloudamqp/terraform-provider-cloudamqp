resource "cloudamqp_instance" "instance" {
	name = "{{or .InstanceName `TestAccInstance`}}"
	plan = "{{or .InstancePlan `bunny-1`}}"
	region = "{{or .InstanceRegion `amazon-web-services::us-east-1`}}"
	tags = ["{{or .InstanceTags `terraform`}}"]
}

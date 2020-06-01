---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_integration_log"
description: |-
  Creates and manages third party log integration for a CloudAMQP instance.
---

# cloudamqp_integration_log

This resource allows you to create and manage third party log integrations for a CloudAMQP instance. Once configured, the logs produced will be forward to corresponding integration. This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

Only available for dedicated subscription plans.

## Example Usage

```hcl
resource "cloudamqp_integration_log" "cloudwatch" {
  instance_id = cloudamqp_instance.instance.id
  access_key_id = var.aws_access_key_id
  region = "eu-north-1"
  secret_access_key = var.aws_secret_access_key
  name = "cloudwatchlog"
}
```

## Argument Reference

The following arguments are supported:

* `name`              - (Required) The name of the third party log integration. See
* `url`               - (Optional) Endpoint to log integration.
* `host_port`         - (Optional) Destination to send the logs.
* `token`             - (Optional/Sensitive) Token used for authentication.
* `region`            - (Optional) Region hosting the integration service.
* `access_key_id`     - (Optional/Sensitive) AWS access key identifier.
* `secret_access_key` - (Optional/Sensitive) AWS secret access key.

This is the full list of all arguments. But a subset of arguments are used based on which type of integration used. See table below for more information.

## Argument Reference (cloudwatchlog)

Cloudwatch argument reference and example. Create an IAM user with programmatic access and the following permissions:

* CreateLogGroup
* CreateLogStream
* DescribeLogGroups
* DescribeLogStreams
* PutLogEvents

## Integration service reference

Valid names for third party log integration.

| Name       | Description |
|------------|---------------------------------------------------------------|
| cloudwatchlog | Create a IAM with programmatic access. |
| logentries | Create a Logentries token at https://logentries.com/app#/add-log/manual  |
| loggly     | Create a Loggly token at https://{your-company}.loggly.com/tokens |
| papertrail | Create a Papertrail endpoint https://papertrailapp.com/systems/setup |
| splunk     | Create a HTTP Event Collector token at https://.cloud.splunk.com/en-US/manager/search/http-eventcollector |

## Integration Type reference

Valid arguments for third party log integrations.

Required arguments for all integrations: name

| Name | Type | Required arguments |
| ---- | ---- | ---- |
| CloudWatch | cloudwatchlog | access_key_id, secret_access_key, region |
| Log Entries | logentries | token |
| Loggly | loggly | token |
| Papertrail | papertrail | url |
| Splunk | splunk | token, host_port

## Import
`cloudamqp_integration_log`can be imported using the name argument of the resource together with CloudAMQP instance identifier. The name and identifier are CSV separated, see example below.

`terraform import cloudamqp_integration_log.cloudwatchlog <name>,<instance_id>`

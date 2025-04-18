---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_integration_aws_eventbridge"
description: |-
  Creates and manages an AWS EventBridge for a CloudAMQP instance.
---

# cloudamqp_integration_aws_eventbridge

This resource allows you to create and manage, an [AWS EventBridge] for a CloudAMQP instance. Once
created, continue to map the EventBridge in the [AWS Eventbridge console].

~>  Our consumer needs to have exclusive usage to the configured queue and the maximum body size
allowed on msgs by AWS is 256kb. The message body has to be valid JSON for AWS Eventbridge to accept
it. If messages are too large or are not valid JSON, they will be rejected (tip: setup a dead-letter
queue to catch them).

Not possible to update this resource. Any changes made to the argument will destroy and recreate the
resource. Hence why all arguments use ForceNew.

Only available for dedicated subscription plans.

## Example Usage

```hcl
resource "cloudamqp_instance" "instance" {
  name        = "Test instance"
  plan        = "penguin-1"
  region      = "amazon-web-services::us-west-1"
  rmq_version = "3.11.5"
  tags        = ["aws"]
}

resource "cloudamqp_integration_aws_eventbridge" "this" {
  instance_id     = cloudamqp_instance.instance.id
  vhost           = cloudamqp_instance.instance.vhost
  queue           = "<QUEUE-NAME>"
  aws_account_id  = "<AWS-ACCOUNT-ID>"
  aws_region      = "us-west-1"
  with_headers    = true
}
```

## Argument References

The following arguments are supported:

* `aws_account_id`  - (ForceNew/Required) The 12 digit AWS Account ID where you want the events to
                      be sent to.
* `aws_region`      - (ForceNew/Required) The AWS region where you the events to be sent to.
                      (e.g. us-west-1, us-west-2, ..., etc.)
* `vhost`           - (ForceNew/Required) The VHost the queue resides in.
* `queue`           - (ForceNew/Required) A (durable) queue on your RabbitMQ instance.
* `with_headers`    - (ForceNew/Required) Include message headers in the event data.
                      `({ "headers": { }, "body": { "your": "message" } })`

## Attributes Reference

All attributes reference are computed

* `id`      - The identifier for this resource.
* `status`  - Always set to null, unless there is an error starting the EventBridge.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_integration_aws_eventbridge` can be imported using the resource identifier together with
CloudAMQP instance identifier (CSV separated). To retrieve the resource identifier, use
[CloudAMQP API list eventbridges].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_integration_aws_eventbridge.this
  id = format("<id>,%s", cloudamqp_instance.instance.id)
}
```

Or with Terraform CLI:

`terraform import cloudamqp_integration_aws_eventbridge.this <id>,<instance_id>`

[AWS EventBridge]: https://aws.amazon.com/eventbridge
[AWS Eventbridge console]: https://console.aws.amazon.com/events/home
[CloudAMQP API list eventbridges]: https://docs.cloudamqp.com/cloudamqp_api.html#list-eventbridges

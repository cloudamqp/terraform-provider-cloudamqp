---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_notification"
description: |-
  Creates and manages recipients to receive alarm notifications.
---

# cloudamqp_notification

This resource allows you to create and manage recipients to receive alarm notifications. There will
always be a default recipient created upon instance creation. This recipient will use team email and
receive notifications from default alarms.

Available for all subscription plans.

## Example Usage

<details>
  <summary>
    <b>Email recipient</b>
  </summary>

```hcl
resource "cloudamqp_notification" "email_recipient" {
  instance_id = cloudamqp_instance.instance.id
  type        = "email"
  value       = "alarm@example.com"
  name        = "alarm"
}
```

</details>

<details>
  <summary>
    <b>OpsGenie recipient with optional responders</b>
  </summary>

```hcl
resource "cloudamqp_notification" "opsgenie_recipient" {
  instance_id = cloudamqp_instance.instance.id
  type        = "opsgenie" # or "opsgenie-eu"
  value       = "<api-key>"
  name        = "OpsGenie"
  responders {
    type = "team"
    id   = "<team-uuid>"
  }
  responders {
    type      = "user"
    username  = "<username>"
  }
}
```

</details>

<details>
  <summary>
    <b>Pagerduty recipient with optional dedup key</b>
  </summary>

```hcl
resource "cloudamqp_notification" "pagerduty_recipient" {
  instance_id = cloudamqp_instance.instance.id
  type        = "pagerduty"
  value       = "<integration-key>"
  name        = "PagerDuty"
  options     = {
    "dedupkey" = "DEDUPKEY"
  }
}
```

</details>

<details>
  <summary>
    <b>Signl4 recipient</b>
  </summary>

```hcl
resource "cloudamqp_notification" "signl4_recipient" {
  instance_id = cloudamqp_instance.instance.id
  type        = "signl4"
  value       = "<team-secret>"
  name        = "Signl4"
}
```

</details>

<details>
  <summary>
    <b>Teams recipient</b>
  </summary>

```hcl
resource "cloudamqp_notification" "teams_recipient" {
  instance_id = cloudamqp_instance.instance.id
  type        = "teams"
  value       = "<teams-webhook-url>"
  name        = "Teams"
}
```

</details>

<details>
  <summary>
    <b>Victorops recipient with optional routing key (rk)</b>
  </summary>

```hcl
resource "cloudamqp_notification" "victorops_recipient" {
  instance_id = cloudamqp_instance.instance.id
  type        = "victorops"
  value       = "<integration-key>"
  name        = "Victorops"
  options     = {
    "rk" = "ROUTINGKEY"
  }
}
```

</details>

<details>
  <summary>
    <b>Slack recipient</b>
  </summary>

```hcl
resource "cloudamqp_notification" "slack_recipient" {
  instance_id = cloudamqp_instance.instance.id
  type        = "slack"
  value       = "<slack-webhook-url>"
  name        = "Slack webhook recipient"
}
```

</details>

<details>
  <summary>
    <b>Webhook recipient</b>
  </summary>

```hcl
resource "cloudamqp_notification" "webhook_recipient" {
  instance_id = cloudamqp_instance.instance.id
  type        = "webhook"
  value       = "<webhook-url>"
  name        = "Webhook"
}
```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance ID.
* `type`        - (Required) Type of the notification. See valid options below.
* `value`       - (Required) Integration/API key or endpoint to send the notification.
* `name`        - (Optional) Display name of the recipient.
* `options`     - (Optional) Options argument (e.g. `rk` used for VictorOps routing key).
* `responders`  - (Optional) An array of reponders (only for OpsGenie). Each `responders` block
                  consists of the field documented below.

___

The options parameter:

* rk        - (Optional) Routing key to route alarm notification (can be used with Victorops).
* dedupkey  - (Optional) If multiple alarms are triggered using a recipient with this key, only the
              the first alarm will trigger a notification (can be used with PagerDuty). Leave blank
              to use the generated dedup key.

___

The `responders` block consists of:

* `type`      - (Required) Type of responder. [`team`, `user`, `escalation`, `schedule`]
* `id`        - (Optional) Identifier in UUID format
* `name`      - (Optional) Name of the responder
* `username`  - (Optional) Username of the responder

Responders of type `team`, `escalation` and `schedule` can use either id or name.
While `user` can use either id or username.

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Notification type reference

Valid options for notification type.

* email
* opsgenie
* opsgenie-eu
* pagerduty
* signl4
* slack
* teams
* victorops
* webhook

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_notification` can be imported using the resource identifier together with CloudAMQP
instance identifier (CSV separated). To retrieve the resource identifier, use
[CloudAMQP API list recipients].

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_notification.recipient
  id = format("<id>,%s", cloudamqp_instance.instance.id)
}
```

Or use Terraform CLI:

`terraform import cloudamqp_notification.recipient <id>,<instance_id>`

[CloudAMQP API list recipients]: https://docs.cloudamqp.com/cloudamqp_api.html#list-recipients

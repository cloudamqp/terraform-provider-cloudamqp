---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_team_member"
description: |-
  Manages a team member invite/membership.
---

# cloudamqp_team_member

This resource allows you to manage a team member as a single resource across both states it can
be in: a pending invite before it's accepted, and an accepted member afterwards. `status` reports
which one currently applies. The resource transparently keeps working as the invite transitions from one to the other.

`role` defaults to `"member"` if not set. `tags` are only meaningful for the `"member"`/`"monitor"`
roles. A tag name that doesn't exist yet is created automatically, the same way tagging an instance
or VPC works.

For a SAML-configured team, the IdP may overwrite `role`/`tags` again on the member's next login
regardless of what Terraform configures. If that happens, the next plan will show the drift like
any other out-of-band change.

## Example Usage

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# Invite a new team member
resource "cloudamqp_team_member" "member" {
  email = "user@example.com"
  role  = "member"
  tags  = ["production"]
}
```

## Argument Reference

* `email` - (Required) Email address to invite. Changing this forces a new resource.
* `role`  - (Optional/Computed) Role assigned to the team member. Defaults to `"member"` if omitted.
  One of: `admin`, `billing manager`, `member`, `monitor`, `devops`.
* `tags`  - (Optional/Computed) Tag names granting access to tagged instances. A tag that doesn't
  exist yet is created automatically. Only meaningful for the `"member"`/`"monitor"` roles.

## Attributes Reference

All attributes reference are computed

* `id`              - The member's normalized (lowercased) email, used as this resource's
  identifier.
* `status`          - `"pending"` while the invite hasn't been accepted yet, `"member"`
  afterwards.
* `invite_id`       - The pending invite's id. Null once `status` is `"member"`.
* `user_id`         - The accepted member's user id. Null while `status` is `"pending"`.

## Import

`cloudamqp_team_member` can be imported using the member's email as identifier.

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_team_member.member
  id = "user@example.com"
}
```

Or with Terraform CLI:

`terraform import cloudamqp_team_member.member user@example.com`

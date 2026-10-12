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
roles. Either create new ones with [`cloudamqp_team_tag`], or reference existing tags (created outside
Terraform, or in another Terraform configuration) via the [`cloudamqp_team_tags`] data source.

For SAML-configured teams, `role`/`tags` are not managed by Terraform, the IdP controls them on
login. Configured values are still sent once on invite, but any plan/apply against an existing
SAML-configured member ignores configuration changes to `role`/`tags` with a warning, since the
IdP may overwrite them again on the member's next login regardless of what Terraform does. Leave
`role`/`tags` unset in config for a SAML-managed member so they can track whatever the IdP assigns
without Terraform fighting it with a perpetual diff.

## Example Usage

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

resource "cloudamqp_team_tag" "production" {
  name = "production"
}

# Invite a new team member
resource "cloudamqp_team_member" "member" {
  email = "user@example.com"
  role  = "member"
  tags  = [cloudamqp_team_tag.production.name]
}
```

## Argument Reference

* `email` - (Required) Email address to invite. Changing this forces a new resource.
* `role`  - (Optional/Computed) Role assigned to the team member. Defaults to `"member"` if omitted.
  One of: `admin`, `billing manager`, `member`, `monitor`, `devops`.
* `tags`  - (Optional/Computed) Tags granting access to tagged instances. Only meaningful for the
  `"member"`/`"monitor"` roles.

## Attributes Reference

All attributes reference are computed

* `id`              - The member's normalized (lowercased) email, used as this resource's
  identifier.
* `status`          - `"pending"` while the invite hasn't been accepted yet, `"member"`
  afterwards.
* `invite_id`       - The pending invite's id. Null once `status` is `"member"`.
* `user_id`         - The accepted member's user id. Null while `status` is `"pending"`.
* `saml_configured` - True when the team has SAML configured, in which case `role`/`tags` are not
  managed by Terraform.

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

[`cloudamqp_team_tag`]: ./team_tag.md
[`cloudamqp_team_tags`]: ../data-sources/team_tags.md

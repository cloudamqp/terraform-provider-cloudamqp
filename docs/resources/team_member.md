---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_team_member"
description: |-
  Manages a team member invite/membership.
---

# cloudamqp_team_member

This resource manages a team member through both of its states: a pending invite, and an
accepted member once the invite is accepted. `status` shows which state applies.

## Example Usage

```hcl
resource "cloudamqp_team_member" "this" {
  email = "user@example.com"
  role  = "member"
  tags  = ["production"]
}
```

## Argument Reference

* `email` - (Required) Email address to invite. Changing this forces a new resource.
* `role`  - (Required) Role assigned to the team member. One of: `admin`, `billing manager`,
  `member`, `monitor`, `devops`. Changing this forces a new resource while the invite is pending.
* `tags`  - (Optional) Set of tag names granting access to tagged instances. The order
  doesn't matter. A tag that doesn't exist yet is created automatically. Required (at least one)
  for the `member`/`monitor` roles, not allowed for other roles. Changing this forces a new
  resource while the invite is pending.

## Attributes Reference

All attributes are computed.

* `id`              - The member's normalized (lowercased) email, used as this resource's
  identifier.
* `status`          - `pending` while the invite hasn't been accepted yet, `member`
  afterwards.
* `invite_id`       - The pending invite's id. Null once `status` is `member`.
* `user_id`         - The accepted member's user id. Null while `status` is `pending`.

## Import

`cloudamqp_team_member` can be imported using the member's email as identifier.

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_team_member.this
  id = "user@example.com"
}
```

Or with Terraform CLI:

`terraform import cloudamqp_team_member.this user@example.com`

## Notes

### Pending invites

A pending invite can't be updated. If you change `role` or `tags` while `status` is `pending`,
the invite is revoked and a new one is sent. The plan shows this as a replacement. Note that
`terraform apply -auto-approve` skips the confirmation. After the invite is accepted, the same
changes are updated in place.

If the invite is accepted between plan and apply, the apply stops and nothing is changed. Run
plan again to update the member instead.

### Tag letter case

Tag access is case-sensitive. If the team already has `"Prod"` and you configure `"prod"`, the
member gets the existing `"Prod"` tag. Terraform keeps your spelling in state and shows a
warning on apply. For the same reason, tags that only differ in case, like `["Prod", "prod"]`,
are rejected.

### SAML teams

The IdP may overwrite `role` and `tags` on the member's next login. The next plan then shows
this as drift. To let the IdP own these values, ignore changes to them:

```hcl
  lifecycle {
    ignore_changes = [role, tags]
  }
```

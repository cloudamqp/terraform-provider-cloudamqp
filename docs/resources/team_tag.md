---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_team_tag"
description: |-
  Manages a team tag.
---

# cloudamqp_team_tag

This resource allows you to manage team tags, used to control which team members can access
which instances. Create a tag here, then assign it to a member's role with
[`cloudamqp_team_member`] and to an instance's `tags` to grant that member access to just that
instance.

Not supported for SAML-configured teams.

## Example Usage

```hcl
# Configure CloudAMQP provider
provider "cloudamqp" {
  apikey = var.cloudamqp_customer_api_key
}

# Team tag resource
resource "cloudamqp_team_tag" "production" {
  name = "production"
}
```

## Argument Reference

* `name` - (Required) Name of the tag. Can't contain a comma. Changing this forces a new resource,
  since tags can't be renamed (delete and recreate instead).

## Attributes Reference

All attributes reference are computed

* `id`     - The tag's name, used as this resource's identifier.
* `tag_id` - The tag's underlying numeric id, used internally to delete it.

## Import

`cloudamqp_team_tag` can be imported using the tag name as identifier.

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_team_tag.production
  id = "production"
}
```

Or with Terraform CLI:

`terraform import cloudamqp_team_tag.production production`

[`cloudamqp_team_member`]: ./team_member.md

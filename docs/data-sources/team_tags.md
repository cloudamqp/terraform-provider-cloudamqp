---
layout: "cloudamqp"
page_title: "CloudAMQP: data source cloudamqp_team_tags"
description: |-
  List all tags for the team.
---

# cloudamqp_team_tags

Use this data source to list all tags for the team. Uses the included apikey in provider
configuration to determine which team to read from.

Useful for referencing existing tags (e.g. created outside Terraform, or managed in another
Terraform configuration) without managing them via [`cloudamqp_team_tag`].

## Example Usage

```hcl
provider "cloudamqp" {
  apikey = "<apikey>"
}

data "cloudamqp_team_tags" "tags" {}

output "tag_names" {
  value = [for tag in data.cloudamqp_team_tags.tags.tags : tag["name"]]
}
```

## Attributes Reference

All attributes reference are computed

* `id`   - The identifier for this data source. Fixed to `team_tags`, since there is no unique
  identifier.
* `tags` - An array of tags. Each `tags` block consists of the fields documented below.

___

The `tags` block consist of

* `name` - The tag name.

## Dependency

This data source depends on apikey set in the provider configuration.

[`cloudamqp_team_tag`]: ../resources/team_tag.md

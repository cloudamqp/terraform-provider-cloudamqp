---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_custom_domain"
description: |-
  Configure and manage your custom domain
---

# cloudamqp_custom_domain

This resource allows you to configure and manage your custom domain for the CloudAMQP instance.

Adding a custom domain to your instance will generate a TLS certificate from [Let's Encrypt], for the given hostname, and install it on all servers in your cluster. The certificate will be automatically renewed going forward.

⚠️ Please note that when creating, changing or deleting the custom domain, the listeners on your servers will be restarted in order to apply the changes. This will close your current connections.

Your custom domain name needs to point to your CloudAMQP hostname, preferably using a CNAME DNS record.

Only available for dedicated subscription plans.

## Example Usage

```hcl
resource "cloudamqp_custom_domain" "settings" {
  instance_id = cloudamqp_instance.instance.id
  hostname = "myname.mydomain"
}
```

## Argument Reference

Top level argument reference

* `instance_id` - (Required) The CloudAMQP instance ID.
* `hostname`    - (Required) Your custom domain name.

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_custom_domain` can be imported using CloudAMQP instance identifier.

`terraform import cloudamqp_custom_domain.settings <instance_id>`

[Let's Encrypt]: https://letsencrypt.org/

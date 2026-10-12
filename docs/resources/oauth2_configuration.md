---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_oauth2_configuration"
description: |-
  Configure OAuth2 authentication for RabbitMQ or LavinMQ
---

<!-- markdownlint-disable MD033 -->

# cloudamqp_oauth2_configuration

This resource allows you to configure OAuth2 authentication for your RabbitMQ or LavinMQ instance.

Only available for dedicated subscription plans running ***RabbitMQ*** or ***LavinMQ***. LavinMQ
requires version 2.7.0 or newer.

The resource looks up which broker the instance runs when creating or updating the configuration.
Some attributes are only supported by one of the brokers, see [Broker support](#broker-support).
Setting an attribute the broker doesn't support fails the apply with an error.

## Example Usage

<details>
  <summary>
    <b>
      <i>Basic OAuth2 configuration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_oauth2_configuration" "oauth2_config" {
  instance_id       = cloudamqp_instance.instance.id
  resource_server_id = "test-resource-server"
  issuer            = "https://test-issuer.example.com"
  verify_aud        = true
  oauth_client_id   = "test-client-id"
  oauth_scopes      = ["read", "write"]
}
```

</details>

<details>
  <summary>
    <b>
      <i>OAuth2 configuration with all optional fields</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_oauth2_configuration" "oauth2_config" {
  instance_id                = cloudamqp_instance.instance.id
  resource_server_id         = "test-resource-server"
  issuer                     = "https://test-issuer.example.com"
  preferred_username_claims  = ["preferred_username", "username"]
  additional_scopes_key      = ["admin"]
  scope_prefix               = "cloudamqp"
  scope_aliases = {
    read  = "read:all"
    write = "write:all"
  }
  verify_aud         = true
  oauth_client_id    = "test-client-id"
  oauth_scopes       = ["read", "write", "admin"]
  audience           = "https://test-audience.example.com"
  disable_basic_auth = true
}
```

</details>

<details>
  <summary>
    <b>
      <i>Minimal OAuth2 configuration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_oauth2_configuration" "oauth2_config" {
  instance_id        = cloudamqp_instance.instance.id
  resource_server_id = "test-resource-server"
  issuer             = "https://test-issuer.example.com"
}
```

</details>

<details>
  <summary>
    <b>
      <i>LavinMQ OAuth2 configuration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_oauth2_configuration" "oauth2_config" {
  instance_id               = cloudamqp_instance.instance.id
  issuer                    = "https://lavinmq-demo-keycloak.84.codes/realms/lavinmq-demo"
  resource_server_id        = "wamschat-prod"
  verify_aud                = true
  preferred_username_claims = ["preferred_username"]
}
```

</details>

<details>
  <summary>
    <b>
      <i>LavinMQ OAuth2 configuration with management UI SSO login</i>
    </b>
  </summary>

Requires LavinMQ 2.9.0 or newer, and 2.9.1 or newer for `mgmt_scopes`.

```hcl
resource "cloudamqp_oauth2_configuration" "oauth2_config" {
  instance_id           = cloudamqp_instance.instance.id
  issuer                = "https://auth.example.com/realms/lavinmq"
  resource_server_id    = "lavinmq"
  additional_scopes_key = ["roles"]
  oauth_client_id       = "lavinmq-mgmt"
  mgmt_scopes           = "openid profile email"
  jwks_cache_ttl        = 3600
}
```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id`                - (Required) The CloudAMQP instance ID.
* `resource_server_id`         - (Required) Resource server identifier used to identify the resource
                                 server in OAuth2 tokens.
* `issuer`                     - (Required) The issuer URL of the OAuth2 provider. This is typically
                                 the base URL of your OAuth2 provider (e.g., Auth0, Keycloak, etc.).
* `preferred_username_claims`  - (Optional) List of JWT claims to use as the preferred username.
                                 The first claim found in the token will be used as the username.
* `additional_scopes_key`      - (Optional) List of additional JWT claim keys to extract OAuth2
                                 scopes from.
* `scope_prefix`               - (Optional) Prefix to add to scopes. This is useful when scopes in
                                 the JWT token need to be prefixed for RabbitMQ permissions.
* `scope_aliases`              - (Optional) Map of scope aliases to translate scope names. This allows
                                 mapping OAuth2 scopes to RabbitMQ permission tags.
* `verify_aud`                 - (Optional/Computed) Whether to verify the audience claim in the JWT
                                 token. When not set, RabbitMQ is configured with false and LavinMQ
                                 uses the API default, true.
* `oauth_client_id`            - (Optional) OAuth2 client ID used for token validation. For LavinMQ,
                                 this is the client ID for management UI SSO login, which requires
                                 LavinMQ 2.9.0 or newer.
* `oauth_scopes`               - (Optional) List of OAuth2 scopes to request. These scopes will be
                                 used when obtaining access tokens.
* `audience`                   - (Optional) The audience to be passed along to the Oauth2 provider when
                                 logging in to the management interface. Must be configured for Auth0,
                                 cannot be configured for Entra ID v2.
* `disable_basic_auth`         - (Optional/Computed) Disable static username/password management interface access.
                                 RabbitMQ only.
* `jwks_cache_ttl`             - (Optional) How long JWKS signing keys are cached before being re-fetched
                                 from the identity provider, in seconds. Between 0 and 86400. LavinMQ only.
* `mgmt_scopes`                - (Optional) Space-separated OAuth2 scopes requested during management UI
                                 SSO login, used together with `oauth_client_id`. Requires LavinMQ 2.9.1 or newer.
                                 LavinMQ only.
* `sleep`                      - (Optional) Configurable sleep time in seconds between retries for
                                 OAuth2 configuration. Default set to 60 seconds.
* `timeout`                    - (Optional) Configurable timeout time in seconds for OAuth2
                                 configuration. Default set to 3600 seconds.

## Broker support

| Attribute                   | RabbitMQ | LavinMQ                           |
|-----------------------------|----------|-----------------------------------|
| `resource_server_id`        | Yes      | Yes                               |
| `issuer`                    | Yes      | Yes                               |
| `preferred_username_claims` | Yes      | Yes                               |
| `additional_scopes_key`     | Yes      | Yes                               |
| `scope_prefix`              | Yes      | Yes                               |
| `verify_aud`                | Yes      | Yes                               |
| `audience`                  | Yes      | Yes                               |
| `oauth_client_id`           | Yes      | Yes, LavinMQ 2.9.0 or newer       |
| `scope_aliases`             | Yes      | No                                |
| `oauth_scopes`              | Yes      | No                                |
| `disable_basic_auth`        | Yes      | No                                |
| `jwks_cache_ttl`            | No       | Yes                               |
| `mgmt_scopes`               | No       | Yes, LavinMQ 2.9.1 or newer       |

The CloudAMQP API requires LavinMQ 2.7.0 or newer for OAuth2 configuration.

## Attributes Reference

All attributes reference are computed

* `id` - The identifier for this resource.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_oauth2_configuration` can be imported using the CloudAMQP instance identifier.

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_oauth2_configuration.oauth2_config
  id = cloudamqp_instance.instance.id
}
```

Or use Terraform CLI:

`terraform import cloudamqp_oauth2_configuration.oauth2_config <instance_id>`

## Notes

* Changes to `instance_id` will force recreation of the resource.
* OAuth2 configuration changes are applied asynchronously and may take some time to complete. The
  resource will poll for job completion using the configured `sleep` and `timeout` values.
* Only one OAuth2 configuration can exist per instance. Creating a new configuration will replace
  any existing configuration.
* After a configuration has been applied, a restart of RabbitMQ is required for the changes to take effect.
* For LavinMQ with `oauth_client_id` set, the management UI base URL used for the SSO login callback
  defaults to the instance hostname. Overriding it (`mgmt_base_url` in the API) is not supported by
  this resource.

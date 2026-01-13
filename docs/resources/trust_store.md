---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_trust_store"
description: |-
  Configure trust store for RabbitMQ
---

# cloudamqp_trust_store

This resource allows you to configure a trust store for your RabbitMQ broker. The trust store
enables RabbitMQ to fetch and use CA certificates from an external source for validating client
certificates, or upload multiple leaf certificates as an allow list.

The `http.cacert` and `file.certificates` fields use **WriteOnly**, meaning no information is
present in plan phase, logs or stored in the state for security purposes. To update these fields,
increment either the `version` or update the `key_id` attribute.

-> **Note:** Updates to write-only fields (`http.cacert` or `file.certificates`) are only applied
when `version` is incremented or `key_id` is changed. This design allows you to manage certificate
rotation explicitly.

Only available for dedicated subscription plans running ***RabbitMQ***.

## Example Usage

<details>
  <summary>
    <b>
      <i>Trust store configuration with HTTP provider</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_trust_store" "trust_store" {
  instance_id = cloudamqp_instance.instance.id
  http {
    url = "https://example.com/trust-store-certs"
  }
  refresh_interval = 30
}
```

</details>

<details>
  <summary>
    <b>
      <i>Trust store with HTTP provider and CA certificate</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_trust_store" "trust_store" {
  instance_id = cloudamqp_instance.instance.id

  http {
    url    = "https://example.com/trust-store-certs"
    cacert = file("${path.module}/certs/ca.pem")
  }

  refresh_interval = 30
  version          = 1
}
```

</details>

<details>
  <summary>
    <b>
      <i>Trust store with file provider</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_trust_store" "trust_store" {
  instance_id = cloudamqp_instance.instance.id

  file {
    certificates = [
      file("${path.module}/certs/client1.pem"),
      file("${path.module}/certs/client2.pem")
    ]
  }

  refresh_interval = 30
  version          = 1
}
```

</details>

<details>
  <summary>
    <b>
      <i>Certificate rotation with version management</i>
    </b>
  </summary>

Example of incrementing version to trigger update of write-only certificate fields.

```hcl
locals {
  cert_version = 2  # Increment this to update certificates
}

resource "cloudamqp_trust_store" "trust_store" {
  instance_id = cloudamqp_instance.instance.id

  http {
    url    = "https://example.com/trust-store-certs"
    cacert = file("${path.module}/certs/ca-${local.cert_version}.pem")
  }

  refresh_interval = 30
  version          = local.cert_version
}
```

</details>

<details>
  <summary>
    <b>
      <i>Certificate rotation with key identifier</i>
    </b>
  </summary>

Example of using key_id to trigger update of write-only certificate fields. Useful when
integrating with external key management systems like Azure Key Vault.

```hcl
locals {
  cert_key_id = "a918beb8-fee4-4de1-b0d5-873e2cb0eba2"
}

resource "cloudamqp_trust_store" "trust_store" {
  instance_id = cloudamqp_instance.instance.id

  file {
    certificates = [
      file("${path.module}/certs/client-${local.cert_key_id}.pem")
    ]
  }

  refresh_interval = 30
  key_id           = local.cert_key_id
}
```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance identifier.
* `http` - (Optional*) HTTP trust store configuration block. See [HTTP Block](#http-block) below.
* `file` - (Optional*) File trust store configuration block. See [File Block](#file-block) below.
* `refresh_interval` - (Optional/Computed) Interval in seconds for RabbitMQ to refresh the trust
                       store certificates (default: 30).
* `version` - (Optional/Computed) An integer to trigger updates of write-only certificate fields.
              Increment this value to apply changes to ***http.cacert*** or ***file.certificates*** (default: 1).
* `key_id` - (Optional/Computed) A string identifier to trigger updates of write-only certificate fields.
              Change this value to apply changes to ***http.cacert*** or ***file.certificates*** (default: "").
* `sleep`   - (Optional/Computed) Configurable sleep time in seconds between retries for trust store
              operations (default: 10).
* `timeout` - (Optional/Computed) Configurable timeout time in seconds for trust store operations
              (default: 1800).

***Note:*** Either `http` or `file` configuration block must be specified, but not both.

### HTTP Block

The `http` block supports:

* `url`    - (Required) URL to fetch trust store certificates from. RabbitMQ will periodically fetch
             CA certificates from this URL according to the `refresh_interval`.
* `cacert` - (Optional/WriteOnly) PEM-encoded CA certificates used to verify the HTTPS connection to
             the trust store URL. Updates require incrementing `version` or changing `key_id`.

### File Block

The `file` block supports:

* `certificates` - (Required/WriteOnly) List of PEM-encoded x.509 formatted leaf certificates
                   (1-100 certificates). Updates require incrementing `version` or changing `key_id`.

## Attributes Reference

All attributes reference are computed

* `id` - The identifier for this resource (same as `instance_id`).

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_trust_store` can be imported using the CloudAMQP instance identifier.

-> **Note:** Import will read the current trust store configuration but cannot retrieve write-only
fields (`http.cacert` or `file.certificates`). You'll need to set these in your configuration.

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_trust_store.trust_store
  id = cloudamqp_instance.instance.id
}
```

Or use Terraform CLI:

`terraform import cloudamqp_trust_store.trust_store <instance_id>`

## Notes

* Changes to `instance_id` will force recreation of the resource.
* Trust store configuration changes are applied asynchronously and may take some time to complete.
  The resource will poll for job completion using the configured `sleep` and `timeout` values.
* Only one trust store configuration can exist per instance. Creating a new configuration will
  replace any existing configuration.
* RabbitMQ will periodically fetch certificates from the configured URL according to the
  `refresh_interval` setting.
* The trust store is useful for dynamic certificate management where CA certificates may be
  rotated or updated externally.
* Either use `http` or `file` configuration block.
* The `http.cacert` field is write-only. To update the CA certificate, increment the `version` or
  change `key_id` attributes. This triggers the provider to re-apply the certificate.
* The `file.certificates` field is write-only. To update the allow list with certificates, increment
  the `version` or change `key_id` attributes. This triggers the provider to re-apply the certificates.

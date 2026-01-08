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
      <i>Trust store configuration with HTTP provider and CA certificate</i>
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
      <i>Trust store configuration with file system provider</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_trust_store" "trust-store-config" {
  instance_id = cloudamqp_instance.instance.id
  file {
    certificates = [
      file("${path.module}/certs/client.pem"),
      file("${path.module}/certs/client2.pem")
    ]
  }
  version = 1

  sleep   = 10
  timeout = 100
}
```

</details>

<details>
  <summary>
    <b>
      <i>Trust store with custom sleep and timeout</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_trust_store" "trust_store" {
  instance_id = cloudamqp_instance.instance.id
  http {
    url = "https://example.com/trust-store-certs"
  }
  refresh_interval = 60
  sleep            = 30
  timeout          = 3600
}
```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The CloudAMQP instance ID.
* `http` - (Optional*) HTTP trust store configuration block. See [HTTP Block](#http-block) below.
* `file` - (Optional*) File system  trust store configuration block. See [File Block](#file-block) below.
* `refresh_interval` - (Optional/Computed) Interval in seconds to refresh the trust store certificates.
                       Defaults to 30 seconds.
* `version` - (Optional/Computed) Version of write-only certificates. Increment this value to force
              an update of write-only fields like `http.cacert` or `file.certificates`. Defaults to 1.
* `sleep`   - (Optional/Computed) Configurable sleep time in seconds between retries for
              trust store operations. Defaults to 10 seconds.
* `timeout` - (Optional/Computed) Configurable timeout time in seconds for trust store operations.
              Defaults to 1800 seconds (30 minutes).

***Note:*** * Either `http` or `file` configuration block needs to be used.

### HTTP Block

The `http` block supports:

* `url`    - (Required) URL to fetch trust store certificates from. RabbitMQ will periodically
             fetch CA certificates from this URL.
* `cacert` - (Optional) PEM encoded CA certificates used to verify the HTTPS connection to the
             trust store URL. This is a write-only field, updates are only applied when `version`
             is incremented.

### File Block

The `file` block supports:

* `certificates` - (Required) Unencrypted PEM-encoded and x.509 formatted leaf certificates. This is
                    a write-only field, updates are only applied when `version` is incremented.

## Attributes Reference

All attributes reference are computed

* `id` - The identifier for this resource (same as `instance_id`).

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_trust_store` can be imported using the CloudAMQP instance identifier.

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
* The `http.cacert` field is write-only. To update the CA certificate, increment the `version` attribute. This triggers the provider to re-apply the certificate.
* The `file.certificates` field is write-only. To update the allow list with certificates, increment the `version` attribute. This triggers the provider to re-apply the certificates.

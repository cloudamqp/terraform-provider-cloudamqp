---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_custom_certificate"
description: |-
  Upload custom certificate to the cluster
---

# cloudamqp_custom_certificate

This resource allows you to upload a custom certificate to all servers in your cluster. Update is
not supported, all changes require replacement. `ca`, `cert` and `private_key` all use **WriteOnly**,
meaning no information is present in plan phase, logs or stored in the state for security purposes.

~> **WARNING:** Please note that when uploading a custom certificate or restoring to default certificate,
all current connections will be closed.

-> **Note:** Destroying this resource will restore the cluster to use the default CloudAMQP certificate.

Only available for dedicated subscription plans running ***RabbitMQ***.

## Example Usage

<details>
  <summary>
    <b>
      <i>Upload a custom certificate</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_custom_certificate" "cert" {
  instance_id = cloudamqp_instance.instance.id
  
  # Load certificate files from disk
  ca          = file("${path.module}/certs/ca.pem")
  cert        = file("${path.module}/certs/server.crt")
  private_key = file("${path.module}/certs/server.key")
  
  sni_hosts = "cloudamqp.example.com"
}
```

</details>

<details>
  <summary>
    <b>
      <i>With version management for certificate rotation</i>
    </b>
  </summary>

Example of incrementing certificate version, this will trigger a replacement of the current installed
custom certificate and use a newer version.

```hcl
locals {
  cert_version = 2  # Increment this to force certificate replacement
}

resource "cloudamqp_custom_certificate" "cert" {
  instance_id = cloudamqp_instance.instance.id
  
  ca          = file("${path.module}/certs/ca-${local.cert_version}.pem")
  cert        = file("${path.module}/certs/cert-${local.cert_version}.crt")
  private_key = file("${path.module}/certs/key-${local.cert_version}.key")
  
  sni_hosts = "cloudamqp.example.com"
  version   = local.cert_version
}
```

</details>

<details>
  <summary>
    <b>
      <i>With key identifier management for certificate rotation</i>
    </b>
  </summary>

Example of new key identifier, this will trigger a replacement of the current installed
custom certificate and use a another.

```hcl
locals {
  cert_key_id = "a918beb8-fee4-4de1-b0d5-873e2cb0eba2"
}

resource "cloudamqp_custom_certificate" "cert" {
  instance_id = cloudamqp_instance.instance.id

  ca          = file("${path.module}/certs/ca-${local.cert_key_id}.pem")
  cert        = file("${path.module}/certs/cert-${local.cert_key_id}.crt")
  private_key = file("${path.module}/certs/key-${local.cert_key_id}.key")

  sni_hosts = "cloudamqp.example.com"
  key_id    = local.cert_key_id
}
```

Change key identifier to force replacement. E.g. Azure key value identifier.

```hcl
locals {
  cert_key_id = "53f188e8-a81d-4232-b5f1-7379b0223bb1"
}

resource "cloudamqp_custom_certificate" "cert" {
  instance_id = cloudamqp_instance.instance.id

  ca          = file("${path.module}/certs/ca-${local.cert_key_id}.pem")
  cert        = file("${path.module}/certs/cert-${local.cert_key_id}.crt")
  private_key = file("${path.module}/certs/key-${local.cert_key_id}.key")

  sni_hosts = "cloudamqp.example.com"
  key_id    = local.cert_key_id
}
```

</details>

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required/ForceNew) The CloudAMQP instance identifier.
* `ca` - (Required/WriteOnly) The PEM-encoded Certificate Authority (CA).
* `cert` - (Required/WriteOnly) The PEM-encoded server certificate.
* `private_key` - (Required/WriteOnly) The PEM-encoded private key corresponding to the certificate.
* `sni_hosts` - (Required/ForceNew) A hostname (Server Name Indication) that this certificate applies to.
* `version` - (Optional/Computed/ForceNew) An integer based argument to trigger force new (default: 1).
* `key_id` - (Optional/Computed/ForceNew) A string based argument to trigger force new (default: "").

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

This resource cannot be imported due to the WriteOnly nature of the certificate data (ca, cert, private_key). These sensitive values are never stored in state or returned from the API, making import impossible.

# Terraform Provider for CloudAMQP

<!-- markdownlint-disable MD029 -->
<!-- markdownlint-disable MD033 -->
<!-- markdownlint-disable MD034 -->
<!-- markdownlint-disable MD046 -->

Manage your [CloudAMQP](https://www.cloudamqp.com/) LavinMQ and RabbitMQ instances with Terraform. The provider is
published to the [Terraform Registry](https://registry.terraform.io/providers/cloudamqp/cloudamqp).

---

## Prerequisites

1. **Terraform** >= 0.13 — [Installation guide](https://developer.hashicorp.com/terraform/downloads)
2. **CloudAMQP API key** — Sign up at https://www.cloudamqp.com/, then create a key at https://customer.cloudamqp.com/apikeys

> The API key grants access to the Customer API, which manages instances and proxies calls to the per-instance API
(alarms, notifications, plugins, etc.) using the same key.

---

## Example Usage

```hcl
terraform {
  required_providers {
    cloudamqp = {
      source = "cloudamqp/cloudamqp"
    }
  }
}

provider "cloudamqp" {
  apikey = var.cloudamqp_apikey   # or set CLOUDAMQP_APIKEY env var
}

resource "cloudamqp_instance" "this" {
  name   = "my-rabbitmq"
  plan   = "lemur"
  region = "amazon-web-services::us-east-1"
}
```

> **Note on plan changes between shared and dedicated**
> Switching between a shared plan (e.g. `lemur`) and a dedicated plan (e.g. `bunny`) forces destruction of the existing
instance before a new one is created. All data will be lost and a new hostname with a new DNS record will be assigned.

---

### Debug Logging

Enable detailed debug output (includes CloudAMQP provider and underlying HTTP client logs):

```sh
export TF_LOG_PROVIDER=DEBUG
```

## Resources

| Resource | Description |
| --- | --- |
| `cloudamqp_account_actions` | Perform account-level actions such as rotating credentials |
| `cloudamqp_alarm` | Configure alarms that trigger on metric thresholds |
| `cloudamqp_custom_certificate` | Upload a custom TLS certificate |
| `cloudamqp_custom_domain` | Set a custom hostname for the instance |
| `cloudamqp_extra_disk_size` | Resize the disk of a dedicated instance |
| `cloudamqp_instance` | Create and manage a CloudAMQP instance (RabbitMQ or LavinMQ) |
| `cloudamqp_integration_aws_eventbridge` | Connect the instance to AWS EventBridge |
| `cloudamqp_integration_log` | Forward logs to an external log management service |
| `cloudamqp_integration_metric` | Forward metrics to an external metrics service |
| `cloudamqp_integration_metric_prometheus` | Expose a Prometheus-compatible metrics endpoint |
| `cloudamqp_maintenance_window` | Set the preferred window for automatic maintenance |
| `cloudamqp_node_actions` | Perform node-level actions (e.g. restart) |
| `cloudamqp_notification` | Manage notification endpoints (email, Slack, PagerDuty, etc.) |
| `cloudamqp_oauth2_configuration` | Configure OAuth2 authentication for RabbitMQ |
| `cloudamqp_plugin` | Enable or disable built-in RabbitMQ plugins |
| `cloudamqp_plugin_community` | Enable or disable community plugins |
| `cloudamqp_privatelink_aws` | Configure AWS PrivateLink for the instance |
| `cloudamqp_privatelink_azure` | Configure Azure Private Link for the instance |
| `cloudamqp_rabbitmq_configuration` | Tune RabbitMQ broker settings (heartbeat, max connections, etc.) |
| `cloudamqp_security_firewall` | Manage IP allowlist (firewall) rules for the instance |
| `cloudamqp_trust_store` | Manage the instance trust store (CA certificates) |
| `cloudamqp_upgrade_lavinmq` | Trigger a LavinMQ version upgrade |
| `cloudamqp_upgrade_rabbitmq` | Trigger a RabbitMQ version upgrade |
| `cloudamqp_vpc` | Create and manage a dedicated VPC |
| `cloudamqp_vpc_connect` | Manage VPC connect configuration |
| `cloudamqp_vpc_gcp_peering` | Set up GCP VPC peering with a CloudAMQP instance |
| `cloudamqp_vpc_peering` | Set up AWS VPC peering with a CloudAMQP instance |
| `cloudamqp_webhook` | Configure webhooks for instance events |

---

## Data Sources

| Data Source | Description |
| --- | --- |
| `cloudamqp_account` | Retrieve account information |
| `cloudamqp_account_vpcs` | List VPCs associated with the account |
| `cloudamqp_alarm` | Retrieve a single alarm |
| `cloudamqp_alarms` | List all alarms for an instance |
| `cloudamqp_credentials` | Retrieve instance credentials (URL, username, password) |
| `cloudamqp_instance` | Retrieve instance details |
| `cloudamqp_nodes` | List all nodes of an instance |
| `cloudamqp_notification` | Retrieve a single notification endpoint |
| `cloudamqp_notifications` | List all notification endpoints for an instance |
| `cloudamqp_plugins` | List available built-in plugins |
| `cloudamqp_plugins_community` | List available community plugins |
| `cloudamqp_upgradable_versions` | List available upgrade versions for an instance |
| `cloudamqp_vpc_gcp_info` | Retrieve GCP VPC peering information |
| `cloudamqp_vpc_info` | Retrieve VPC information |

---

## Import

Bring existing infrastructure under Terraform management. Find resource IDs via the
[CloudAMQP API](https://docs.cloudamqp.com/index.html#tag/instances) using a key from
https://customer.cloudamqp.com/apikeys.

After importing, run `terraform plan` to confirm the resource is tracked correctly.

### Instance

```hcl
import {
  to = cloudamqp_instance.this
  id = "<instance_id>"
}
```

### Resources that depend on an instance

Resources such as `cloudamqp_alarm` and `cloudamqp_notification` require both the resource ID and
the instance ID, separated by a comma.

```hcl
import {
  to = cloudamqp_notification.recipient
  id = "<resource_id>,<instance_id>"
}

import {
  to = cloudamqp_alarm.alarm
  id = "<resource_id>,<instance_id>"
}
```

> **Terraform CLI (< v1.5.0)**
>
> ```sh
> terraform import cloudamqp_instance.this <instance_id>
> terraform import cloudamqp_notification.recipient <resource_id>,<instance_id>
> terraform import cloudamqp_alarm.alarm <resource_id>,<instance_id>
> ```

---

## Contributor Guide

### Build from Source

**Prerequisites:** [Go](https://go.dev/doc/install) >= 1.24 and `make` must be installed.

1. Clone the repository and build the provider binary:

   ```sh
   git clone https://github.com/cloudamqp/terraform-provider-cloudamqp.git
   cd terraform-provider-cloudamqp
   make clean build
   ```

2. Configure Terraform to use your local binary by adding the following to `~/.terraformrc`:

   ```hcl
   provider_installation {
     dev_overrides {
       "hashicorp/cloudamqp" = "/path/to/terraform-provider-cloudamqp"
     }
     direct {}
   }
   ```

   With `dev_overrides` set, skip `terraform init` and omit the `required_providers` block from your `.tf` files. Run
   `terraform plan` or `terraform apply` directly.

  <details>
    <summary>
      <b>
        <i>Examples of dev_override with HCL configuration with omitted `required_providers`</i>
      </b>
    </summary>

    ```hcl
    provider "cloudamqp" {
      apikey = var.cloudamqp_apikey   # or set CLOUDAMQP_APIKEY env var
    }

    resource "cloudamqp_instance" "this" {
      name   = "my-rabbitmq"
      plan   = "lemur"
      region = "amazon-web-services::us-east-1"
    }
    ```

   </details>

3. After changing the provider code, rebuild and re-run your Terraform command:

### CloudAMQP API

The provider communicates with two CloudAMQP APIs, both using the same* API key:

| API | Base URL | Purpose |
| --- | --- | --- |
| Customer API | https://customer.cloudamqp.com/api | Manage instances (create, delete, plan changes, VPCs) |
| Instance API | https://api.cloudamqp.com | Configure resources on a running instance (alarms, plugins, firewall, integrations, etc.) |

*The Customer API proxies calls to the Instance API, so the provider only needs a single API key.
Full API reference: https://docs.cloudamqp.com

### VCR Testing

The provider uses Terraform Acceptance Tests together with [Go-VCR](https://github.com/dnaeon/go-vcr) to record and
replay HTTP interactions.

**Record a single test** (requires a live CloudAMQP API key in `.env`):

```sh
CLOUDAMQP_RECORD=1 TF_ACC=1 dotenv -f .env go test ./cloudamqp/ -v -run {TestName} -timeout 30m
```

**Replay a single test:**

```sh
TF_ACC=1 go test ./cloudamqp/ -v -run {TestName}
```

**Replay all tests:**

```sh
TF_ACC=1 go test ./cloudamqp/ -v
```

If a test result is cached, pass `-count=1` to force it to re-run. The default timeout is 10 minutes, adjust with
`-timeout`.

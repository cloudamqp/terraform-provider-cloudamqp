---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_integration_log_agent"
description: |-
  Creates and manages agent based log integrations for a CloudAMQP instance.
---

<!-- markdownlint-disable MD033 -->

# cloudamqp_integration_log_agent

~> **Note:** This resource is available from [v1.47.0].

This resource allows you to create and manage agent based log integrations for a CloudAMQP instance.
Once configured, the logs produced will be forwarded to the corresponding integration.

Only available for dedicated subscription plans.

## Example Usage

<details>
  <summary>
    <b>
      <i>CloudWatch log agent integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log_agent" "cloudwatch" {
  instance_id = cloudamqp_instance.instance.id
  cloudwatch {
    iam_role        = var.aws_iam_role
    iam_external_id = var.aws_iam_external_id
    region          = var.aws_region
    log_group       = "CloudAMQP"
    log_stream      = cloudamqp_instance.instance.cluster_name
  }
}
```

* AWS IAM role: `arn:aws:iam::ACCOUNT-ID:role/ROLE-NAME`
* External id: Create your own external identifier that matches the role created. E.g. `cloudamqp-abc123`.

See the [CloudAMQP CloudWatch documentation] for a step-by-step guide on setting up the IAM role and
trust relationship.

</details>

<details>
  <summary>
    <b>
      <i>CloudWatch log agent integration with AWS Terraform provider</i>
    </b>
  </summary>

```hcl
provider "aws" {
  region = var.aws_region
}

resource "cloudamqp_integration_log_agent" "cloudwatch" {
  instance_id = cloudamqp_instance.instance.id
  cloudwatch {
    iam_role        = var.aws_iam_role
    iam_external_id = var.aws_iam_external_id
    region          = var.aws_region
    log_group       = "CloudAMQP"
    log_stream      = cloudamqp_instance.instance.cluster_name
  }
}

resource "aws_cloudwatch_log_group" "this" {
  name              = "CloudAMQP"
  retention_in_days = 30

  tags = {
    Environment = "Production"
  }
}

resource "aws_cloudwatch_log_stream" "this" {
  name           = cloudamqp_instance.instance.cluster_name
  log_group_name = aws_cloudwatch_log_group.this.name
}
```

</details>

<details>
  <summary>
    <b>
      <i>Coralogix log agent integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log_agent" "coralogix" {
  instance_id = cloudamqp_instance.instance.id
  coralogix {
    private_key = var.coralogix_private_key
    region      = "eu2"
    application = "cloudamqp"
    subsystem   = cloudamqp_instance.instance.host
  }
}
```

</details>

<details>
  <summary>
    <b>
      <i>Datadog log agent integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log_agent" "datadog" {
  instance_id = cloudamqp_instance.instance.id
  datadog {
    api_key = var.datadog_api_key
    region  = "us1"
  }
}
```

</details>

<details>
  <summary>
    <b>
      <i>Grafana Cloud log agent integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log_agent" "grafana" {
  instance_id = cloudamqp_instance.instance.id
  grafana {
    endpoint            = var.grafana_endpoint
    grafana_instance_id = var.grafana_instance_id
    api_token           = var.grafana_api_token
  }
}
```

</details>

<details>
  <summary>
    <b>
      <i>Splunk log agent integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log_agent" "splunk" {
  instance_id = cloudamqp_instance.instance.id
  splunk {
    endpoint     = var.splunk_endpoint
    token        = var.splunk_token
    source_type  = "cloudamqp"
  }
}
```

</details>

<details>
  <summary>
    <b>
      <i>Uptrace log agent integration</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_integration_log_agent" "uptrace" {
  instance_id = cloudamqp_instance.instance.id
  uptrace {
    dsn = var.uptrace_dsn
  }
}
```

Find your DSN in the Uptrace project under **Settings → DSN**.
The DSN format is: `https://<token>@otlp.uptrace.dev/<project_id>`

</details>

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) Instance identifier for the CloudAMQP instance.

Exactly one of the following integration blocks must be configured:

<details>
  <summary>
    <b>CloudWatch</b>
  </summary>

The following arguments are used by the `cloudwatch` block.

* `iam_role`        - (Required) AWS IAM role ARN used to assume permissions for the integration.
* `iam_external_id` - (Required) External identifier that matches the trust policy of the IAM role.
* `region`          - (Required) AWS region hosting the CloudWatch log group.
* `log_group`       - (Optional/Computed) The name of the CloudWatch log group. Defaults to `CloudAMQP` if not set.
* `log_stream`      - (Required) The name of the CloudWatch log stream. Recommended to use the cluster name, found in `cloudamqp_instance.instance.cluster_name`.

### IAM permissions

See the [CloudAMQP CloudWatch documentation] for a step-by-step setup guide on configuring the IAM
role and trust relationship, or the [AWS IAM role documentation] for how to create a role with a
cross-account trust policy and an external ID.

</details>

<details>
  <summary>
    <b>Coralogix</b>
  </summary>

The following arguments are used by the `coralogix` block.

* `private_key`         - (Required, Write-only) Coralogix Send-Your-Data API key (starts with `cxtp_...`). Found in Coralogix under **Settings → API Keys**. This value is write-only and will not be stored in state.
* `private_key_version` - (Optional/Computed) Version of the write-only `private_key`. Increment to trigger an update when the key changes (default: `1`).
* `region`              - (Required) Coralogix ingress region. Valid values: `eu1`, `eu2`, `ap1`, `ap2`, `ap3`, `us1`, `us2`, `us3`, `uk1`. See the [Coralogix region documentation] for the region-to-domain mapping.
* `application`         - (Required) Application name used to group logs by environment in Coralogix (e.g. `cloudamqp`).
* `subsystem`           - (Required) Subsystem name used to group logs by service within an application. Recommended to use `cloudamqp_instance.instance.host`.

</details>

<details>
  <summary>
    <b>Datadog</b>
  </summary>

The following arguments are used by the `datadog` block.

* `api_key`         - (Required, Write-only) Datadog API key. Found in Datadog under **Organization Settings → API Keys**. This value is write-only and will not be stored in state.
* `api_key_version` - (Optional/Computed) Version of the write-only `api_key`. Increment to trigger an update when the key changes (default: `1`).
* `region`          - (Required) Datadog ingestion region. Valid values: `us1`, `us3`, `us5`, `eu`, `ap2`.
* `tags`            - (Optional) Comma-separated tags to attach to logs (e.g. `env=prod,region=eu`).

</details>

<details>
  <summary>
    <b>Grafana Cloud</b>
  </summary>

The following arguments are used by the `grafana` block.

* `endpoint`            - (Required) Grafana Cloud OTLP endpoint URL. Format: `https://otlp-gateway-prod-<region>.grafana.net/otlp`. Found in the Grafana Cloud Portal under **Stack → OpenTelemetry → Configure**.
* `grafana_instance_id` - (Required) Grafana Cloud numeric stack instance ID. Found alongside the endpoint in the OpenTelemetry configuration page.
* `api_token`           - (Required, Write-only) Grafana Cloud API token (starts with `glc_eyJ...`). Generate one in the OpenTelemetry configuration page. This value is write-only and will not be stored in state.
* `api_token_version`   - (Optional/Computed) Version of the write-only `api_token`. Increment to trigger an update when the token changes (default: `1`).

To find your credentials:

1. Sign in to the [Grafana Cloud Portal] and open your stack
2. Find the **OpenTelemetry** tile and click **Configure**
3. Copy the **Endpoint URL** and **Instance ID**, and generate an API token

See the [Grafana Cloud OTLP setup guide] for step-by-step instructions.

</details>

<details>
  <summary>
    <b>Splunk</b>
  </summary>

The following arguments are used by the `splunk` block.

* `endpoint`       - (Required) Splunk HTTP Event Collector (HEC) endpoint URL. Format: `https://<instance-id>.splunkcloud.com:443/services/collector`.
* `token`          - (Required, Write-only) Splunk HEC token. This value is write-only and will not be stored in state.
* `token_version`  - (Optional/Computed) Version of the write-only `token`. Increment to trigger an update when the token changes (default: `1`).
* `source_type`    - (Optional) Splunk source type. Leave blank to use the token's default.

To set up:

1. In Splunk, go to **Settings → Data inputs → HTTP Event Collector** and create a new token. Set a **default index** and disable **indexer acknowledgement**.
2. Your HEC URL looks like `https://<your-instance-id>.splunkcloud.com:443/services/collector`.
3. Enter the HEC endpoint, token, and optionally a source type in the resource block.

See the [Splunk HEC documentation] for full setup instructions.

</details>

<details>
  <summary>
    <b>Uptrace</b>
  </summary>

The following arguments are used by the `uptrace` block.

* `dsn`         - (Required, Write-only) Uptrace DSN (Data Source Name) URL. Find this in your Uptrace project under **Settings → DSN** (see [Uptrace DSN documentation]). Format: `https://<token>@otlp.uptrace.dev/<project_id>`. This value is write-only and will not be stored in state.
* `dsn_version` - (Optional/Computed) Version of the write-only `dsn`. Increment this to trigger an update when the DSN changes (default: `1`).

</details>

## Attributes Reference

All attributes reference are computed

* `id` - The identifier for this resource.

## Dependency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

## Import

`cloudamqp_integration_log_agent` can be imported using the resource identifier together with
CloudAMQP instance identifier. The identifiers are CSV separated, see example below.

From Terraform v1.5.0, the `import` block can be used to import this resource:

```hcl
import {
  to = cloudamqp_integration_log_agent.cloudwatch
  id = format("<id>,%s", cloudamqp_instance.instance.id)
}
```

`terraform import cloudamqp_integration_log_agent.cloudwatch <id>,<instance_id>`

[v1.47.0]: https://github.com/cloudamqp/terraform-provider-cloudamqp/releases/tag/v1.47.0
[CloudAMQP CloudWatch documentation]: https://www.cloudamqp.com/docs/monitoring_logs_cloudwatch_v2.html
[AWS IAM role documentation]: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-user_externalid.html
[Coralogix region documentation]: https://coralogix.com/docs/coralogix-domain/
[Grafana Cloud Portal]: https://grafana.com/auth/sign-in
[Grafana Cloud OTLP setup guide]: https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/#manual-opentelemetry-setup-for-advanced-users
[Splunk HEC documentation]: https://docs.splunk.com/Documentation/Splunk/latest/Data/UsetheHTTPEventCollector
[Uptrace DSN documentation]: https://uptrace.dev/get/dsn.html

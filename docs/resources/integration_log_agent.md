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

</details>

<details>
  <summary>
    <b>Uptrace</b>
  </summary>

The following arguments are used by the `uptrace` block.

* `dsn`         - (Required, Write-only) Uptrace DSN (Data Source Name) URL. Find this in your Uptrace project under **Settings → DSN**. Format: `https://<token>@otlp.uptrace.dev/<project_id>`. This value is write-only and will not be stored in state.
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
[Uptrace DSN documentation]: https://uptrace.dev/get/dsn.html

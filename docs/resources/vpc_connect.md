---
layout: "cloudamqp"
page_title: "CloudAMQP: cloudamqp_vpc_connect"
description: |-
  Enable VPC connect (Privatelink or Private Service Connect) for a CloudAMQP instance hosted in
  AWS, Azure or GCP.
---

# cloudamqp_vpc_connect

This resource is a generic way to handle PrivateLink (AWS and Azure) and Private Service Connect (GCP).
Communication between resources can be done just as they were living inside a VPC. CloudAMQP creates an Endpoint
Service to connect the VPC and creating a new network interface to handle the communicate.

If no existing VPC available when enable VPC connect, a new VPC will be created with subnet `10.52.72.0/24`.

More information can be found at: [CloudAMQP VPC Connect](https://www.cloudamqp.com/docs/cloudamqp-vpc-connect.html)

-> **Note:** Enabling VPC Connect will automatically add a firewall rule.

<details>
 <summary>
    <b>
      <i>Default PrivateLink firewall rule [AWS, Azure]</i>
    </b>
  </summary>

```hcl
rules {
  Description = "PrivateLink setup"
  ip          = "<VPC Subnet>"
  ports       = []
  services    = ["AMQP", "AMQPS", "HTTPS", "STREAM", "STREAM_SSL", "STOMP", "STOMPS", "MQTT", "MQTTS"]
}
```

</details>

<details>
 <summary>
    <b>
      <i>Default Private Service Connect firewall rule [GCP]</i>
    </b>
  </summary>

```hcl
rules {
  Description = "Private Service Connect"
  ip          = "10.0.0.0/24"
  ports       = []
  services    = ["AMQP", "AMQPS", "HTTPS", "STREAM", "STREAM_SSL", "STOMP", "STOMPS", "MQTT", "MQTTS"]
}
```

</details>

Only available for dedicated subscription plans.

## Example Usage

<details>
  <summary>
    <b>
      <i>Enable VPC Connect (PrivateLink) in AWS</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_vpc" "vpc" {
  name = "Standalone VPC"
  region = "amazon-web-services::us-west-1"
  subnet = "10.56.72.0/24"
  tags = []
}

resource "cloudamqp_instance" "instance" {
  name   = "Instance 01"
  plan   = "bunny-1"
  region = "amazon-web-services::us-west-1"
  tags   = []
  vpc_id = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_vpc_connect" "vpc_connect" {
  instance_id = cloudamqp_instance.instance.id
  region = cloudamqp_instance.instance.region
  allowed_principals = [
    "arn:aws:iam::aws-account-id:user/user-name"
  ]
}
```

</details>

<details>
  <summary>
    <b>
      <i>Enable VPC Connect (PrivateLink) in Azure</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_vpc" "vpc" {
  name = "Standalone VPC"
  region = "azure-arm::westus"
  subnet = "10.56.72.0/24"
  tags = []
}

resource "cloudamqp_instance" "instance" {
  name   = "Instance 01"
  plan   = "bunny-1"
  region = "azure-arm::westus"
  tags   = []
  vpc_id = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_vpc_connect" "vpc_connect" {
  instance_id = cloudamqp_instance.instance.id
  region = cloudamqp_instance.instance.region
  approved_subscriptions = [
    "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
  ]
}
```

</details>

<details>
  <summary>
    <b>
      <i>Enable VPC Connect (Private Service Connect) in GCP</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_vpc" "vpc" {
  name = "Standalone VPC"
  region = "google-compute-engine::us-west1"
  subnet = "10.56.72.0/24"
  tags = []
}

resource "cloudamqp_instance" "instance" {
  name   = "Instance 01"
  plan   = "bunny-1"
  region = "google-compute-engine::us-west1"
  tags   = []
  vpc_id = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_vpc_connect" "vpc_connect" {
  instance_id = cloudamqp_instance.instance.id
  region = cloudamqp_instance.instance.region
  allowed_projects = [
    "some-project-123456"
  ]
}
```

</details>

## Argument Reference

* `instance_id` - (Required) The CloudAMQP instance identifier.
* `region` - (Required) The region where the CloudAMQP instance is hosted.
* `allowed_principals` - (Optional) List of allowed prinicpals used by AWS, see below table.
* `approved_subscriptions` - (Optional) List of approved subscriptions used by Azure, see below table.
* `allowed_projects` - (Optional) List of allowed projects used by GCP, see below table.
* `sleep` - (Optional) Configurable sleep time (seconds) when enable Private Service Connect.
  Default set to 10 seconds.
* `timeout` - (Optional) Configurable timeout time (seconds) when enable Private Service Connect.
  Default set to 1800 seconds.

___

The `allowed_principals`, `approved_subscriptions` or `allowed_projects` data depends on the provider platform:

| Platform | Description         | Format                                                                                                                             |
|----------|---------------------|------------------------------------------------------------------------------------------------------------------------------------|
| AWS      | IAM ARN principals  | arn:aws:iam::aws-account-id:root<br /> arn:aws:iam::aws-account-id:user/user-name<br /> arn:aws:iam::aws-account-id:role/role-name |
| Azure    | Subscription (GUID) | XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX                                                                                               |
| GCP      | Project IDs*        | 6 to 30 lowercase letters, digits, or hyphens                                                                                      |

*https://cloud.google.com/resource-manager/reference/rest/v1/projects

## Attributes Reference

All attributes reference are computed

* `id`  - The identifier for this resource. Will be same as `instance_id`
* `status`- Private Service Connect status [enable, pending, disable]
* `service_name` - Service name (alias for Azure) of the PrivateLink.
* `active_zones` - Covering availability zones used when creating an endpoint from other VPC. (AWS)

## Depedency

This resource depends on CloudAMQP instance identifier, `cloudamqp_instance.instance.id`.

Since `region` also is required, suggest to reuse the argument from CloudAMQP instance,
`cloudamqp_instance.instance.region`.

## Import

`cloudamqp_vpc_connect` can be imported using CloudAMQP internal identifier.

`terraform import cloudamqp_vpc_connect.vpc_connect <id>`

The resource uses the same identifier as the CloudAMQP instance. To retrieve the identifier for an instance, either use [CloudAMQP customer API](https://docs.cloudamqp.com/#list-instances) or use the data source [`cloudamqp_account`](./data-sources/account.md).

## Create VPC Connect with additional firewall rules

To create a PrivateLink/Private Service Connect configuration with additional firewall rules, it's required to chain the [cloudamqp_security_firewall](https://registry.terraform.io/providers/cloudamqp/cloudamqp/latest/docs/resources/security_firewall)
resource to avoid parallel conflicting resource calls. You can do this by making the firewall
resource depend on the VPC Connect resource, `cloudamqp_vpc_connect.vpc_connect`.

Furthermore, since all firewall rules are overwritten, the otherwise automatically added rules for
the VPC Connect also needs to be added.

## Example usage with additional firewall rules

<details>
  <summary>
    <b>
      <i>CloudAMQP instance in an existing VPC with managed firewall rules</i>
    </b>
  </summary>

```hcl
resource "cloudamqp_vpc" "vpc" {
  name = "Standalone VPC"
  region = "amazon-web-services::us-west-1"
  subnet = "10.56.72.0/24"
  tags = []
}

resource "cloudamqp_instance" "instance" {
  name   = "Instance 01"
  plan   = "bunny-1"
  region = "amazon-web-services::us-west-1"
  tags   = []
  vpc_id = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_vpc_connect" "vpc_connect" {
  instance_id = cloudamqp_instance.instance.id
  allowed_principals = [
    "arn:aws:iam::aws-account-id:user/user-name"
  ]
}

resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id

  rules {
    description = "Custom PrivateLink setup"
    ip          = cloudamqp_vpc.vpc.subnet
    ports       = []
    services    = ["AMQP", "AMQPS", "HTTPS", "STREAM", "STREAM_SSL"]
  }

  rules {
    description = "MGMT interface"
    ip = "0.0.0.0/0"
    ports = []
    services = ["HTTPS"]
  }

  depends_on = [
    cloudamqp_vpc_connect.vpc_connect
   ]
}
```

</details>

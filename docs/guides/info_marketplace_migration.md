---
layout: "cloudamqp"
page_title: "Marketplace migration"
subcategory: "info"
description: |-
  Update resource and dependencies identifiers.
---

# Marketplace migration

If you migrate your CloudAMQP billing from a cloud marketplace to directly with CloudAMQP (or vice
versa), the resource identifiers will change. Below follows an example on how to fetch the new
identifiers and update them.

## Configuration file

Basic example of a standalone managed VPC, CloudAMQP instance and firewall.

```hcl
resource "cloudamqp_vpc" "vpc" {
  name = "instance"
  subnet = "10.56.72.0/24"
  region = "amazon-web-services::us-east-1"
  tags = ["aws"]
}

resource "cloudamqp_instance" "instance" {
  name = "instance"
  plan = "bunny-1"
  region = "amazon-web-services::us-east-1"
  tags = ["aws"]
  vpc_id = cloudamqp_vpc.vpc.id
  keep_associated_vpc = true
}

resource "cloudamqp_security_firewall" "firewall_settings" {
  instance_id = cloudamqp_instance.instance.id
  
  rules {
    ip = "10.56.72.0/24"
    ports = []
    services = ["AMQP","AMQPS"]
  }

  rules {
    ip = "0.0.0.0/0"
    description = "Mgmt interface"
    ports = []
    services = ["HTTPS"]
  }
}
```

## Fetch new identifier

For the basic example above we need to fetch two identifiers to be used when manually update the
state file. The identifiers are both the *VPC ID* and *Instance ID*.

This can be done by using the [CloudAMQP API list instances] and listing all the CloudAMQP instances
in the new team. From the response, the new CloudAMQP instance have the identifier `209` and the
VPC `208`.

```json
[
  {
    "id": 209,
    "name": "instance",
    "plan": "bunny-1",
    "region": "amazon-web-services::us-east-1",
    "tags": [
      "aws"
  ],
    "providerid": "***",
    "vpc_id": 208
  }
]
```

## State file

Update the already populated `id`, `instance_id` and `vpc_id` for each resource with the new
identifier.

### VPC resource

The VPC resource identifier needs to be update with, `id: "208"`.

```json
{
  "mode": "managed",
  "type": "cloudamqp_vpc",
  "name": "vpc",
  "provider": "provider[\"localhost/cloudamqp/cloudamqp\"]",
  "instances": [
    {
      "schema_version": 0,
      "attributes": {
        "id": "208",
        "name": "instance",
        "region": "amazon-web-services::us-east-1",
        "subnet": "10.56.72.0/24",
        "tags": [
          "aws"
        ],
        "vpc_name": "vpc-mfbztwps"
      }
    }
  ]
}
```

### CloudAMQP instance resource

Two identifier needs to be updated,  first `id: "209"` and then `vpc_id: 208`.

```json
{
  "mode": "managed",
  "type": "cloudamqp_instance",
  "name": "instance",
  "provider": "provider[\"localhost/cloudamqp/cloudamqp\"]",
  "instances": [
    {
      "schema_version": 0,
      "attributes": {
        "apikey": "***",
        "backend": "rabbitmq",
        "copy_settings": [],
        "dedicated": true,
        "host": "***.rmq6.dev.cloudamqp.com",
        "host_internal": "***.in.rmq6.dev.cloudamqp.com",
        "id": "209",
        "keep_associated_vpc": true,
        "name": "instance",
        "no_default_alarms": null,
        "nodes": 1,
        "plan": "bunny-1",
        "ready": true,
        "region": "amazon-web-services::us-east-1",
        "rmq_version": "4.0.5",
        "tags": [
          "aws"
        ],
        "url": "amqp://***@***.in.rmq6.dev.cloudamqp.com/***",
        "vhost": "***",
        "vpc_id": 208,
        "vpc_subnet": null
      }
    }
  ]
}
```

### Firewall resource

The firewall resource shares the same `id` as `instance_id` update these to `id: "209"` and
`instance_id: 209`.

```json
{
  "mode": "managed",
  "type": "cloudamqp_security_firewall",
  "name": "firewall_settings",
  "provider": "provider[\"localhost/cloudamqp/cloudamqp\"]",
  "instances": [
    {
      "schema_version": 0,
      "attributes": {
      "id": "209",
      "instance_id": 209,
      "rules": [
        {
          "description": "",
          "ip": "10.56.72.0/24",
          "ports": [],
          "services": [
            "AMQP",
            "AMQPS"
          ]
        },
        {
          "description": "Mgmt interface",
          "ip": "0.0.0.0/0",
          "ports": [],
          "services": [
            "HTTPS"
          ]
        }
      ],
        "sleep": 30,
        "timeout": 1800
      },
      "sensitive_attributes": [],
      "private": "bnVsbA==",
      "dependencies": [
        "cloudamqp_instance.instance",
        "cloudamqp_vpc.vpc"
      ]
    }
  ]
}
```

### More resources

Depending on the resources to update even more identifiers needs to be fetched from the API.

### Apply new indentifiers

Run `terraform apply` to use the new identifiers.

[CloudAMQP API list instances]: https://docs.cloudamqp.com/index.html#tag/instances/get/instances

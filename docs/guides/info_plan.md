---
layout: "cloudamqp"
page_title: "CloudAMQP: subscription plans"
subcategory: "info"
description: |-
  Available subscription plans for CloudAMQP.
---

# Subscription plans

Tables below shows general subscription plans for CloudAMQP for either `RabbitMQ` or `LavinMQ`, for full price list see [cloudamqp](https://www.cloudamqp.com/plans.html).

*Information can differ from your actually valid plans, e.g. your team have been given preview access to unreleased plans. The up to date collection of available plans can be retrieved with your team API access key.*

```shell
curl -u :xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx \
  https://customer.cloudamqp.com/api/plans
```

## Plans using RabbitMQ

`Lemur` and `Tiger` are shared instances and share underlying hardware with other instances. They are also limited to which CloudAMQP provider resources that can be used. Further information on availability on each resource page.

Name | Plan | Type | Nodes
---- | ---- | ---- | ----
Little lemur    | lemur   | shared
Tough Tiger     | tiger   | shared
Sassy Squirrel  | squirrel-1    | dedicated | 1
Big Bunny       | bunny-1,3     | dedicated | 1,3
Roaring Rabbit  | rabbit-1,3,5  | dedicated | 1,3,5
Power Panda     | panda-1,3,5   | dedicated | 1,3,5
Awesome Ape     | ape-1,3,5     | dedicated | 1,3,5
Heavy Hippo     | hippo-1,3,5   | dedicated | 1,3,5
Loud Lion       | lion-1,3,5    | dedicated | 1,3,5
Raging Rhino    | rhino-1       | dedicated | 1

```shell
# Filter out available plans for RabbitMQ
curl -u :xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx \
  https://customer.cloudamqp.com/api/plans?backend=rabbitmq
```

## Plans using LavinMQ

`Lemming`and `Ermine` are shared instances and share underlying hardware with other instances. They are also limited to which CloudAMQP provider resources that can be used. Further information on availability on each resource page.

Name | Plan | Type | Nodes
---- | ---- | ---- | ----
Loyal Lemming       | lemming  | shared
Elegant Ermine      | ermine   | shared
Passionate Puffin   | puffin-1    | dedicated | 1
Playful Penguin     | penguin-1   | dedicated | 1
Lively Lynx         | lynx-1      | dedicated | 1
Wild Wolverine      | wolverine-1 | dedicated | 1
Remarkable Reindeer | reindeer-1  | dedicated | 1
Brave Bear          | bear-1      | dedicated | 1
Outstanding Orca    | orca-1      | dedicated | 1

```shell
# Filter out available plans for LavinMQ
curl -u :xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx \
  https://customer.cloudamqp.com/api/plans?backend=lavinmq
```

<br>

# Legacy subscription plans

Table below shows deprecated subscription plans for CloudAMQP. Existing plans will still work, but there will not be possible to create new ones.

Name | Plan | Type
---- | ---- | ----
Little lemur    | lemur   | shared
Tough Tiger     | tiger   | shared
Big Bunny       | bunny   | dedicated
Roaring Rabbit  | rabbit  | dedicated
Power Panda     | panda   | dedicated
Awesome Ape     | ape     | dedicated
Heavy Hippo     | hippo   | dedicated
Loud Lion       | lion    | dedicated

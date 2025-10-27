---
layout: "cloudamqp"
page_title: "CloudAMQP: subscription plans"
subcategory: "info"
description: |-
  Available subscription plans for CloudAMQP.
---

# Subscription plans

Tables below shows general subscription plans for CloudAMQP for either [**LavinMQ**] or
[**RabbitMQ**]. For full price list see [CloudAMQP plans].

-> Information can differ from your actually valid plans, e.g. your team have been given preview
access to unreleased plans. To retrieve an up to date list check out [CloudAMQP API plans].

## Plans using LavinMQ

All plans running LavinMQ backend.

`Lemming`and `Ermine` are shared instances and share underlying hardware with other instances. They
are also limited to which CloudAMQP provider resources that can be used. Further information on
availability on each resource page.

-> With LavinMQ 2.0 version, [clustering] is supported to ensure high availability.

Name                | Plan           | Type      | Nodes
--------------------|----------------|-----------|------
Loyal Lemming       | lemming        | shared    | -
Elegant Ermine      | ermine         | shared    | -
Passionate Puffin   | puffin-1,3,5   | dedicated | 1,3,5
Playful Penguin     | penguin-1,3,5  | dedicated | 1,3,5
Lively Lynx         | lynx-1,3,5     | dedicated | 1,3,5
Wild Wolverin       | wolverin-1,3,5 | dedicated | 1,3,5
Remarkable Reindeer | reindeer-1,3,5 | dedicated | 1,3,5
Brave Bear          | bear-1,3,5     | dedicated | 1,3,5
Outstanding Orca    | orca-1,3,5     | dedicated | 1,3,5

## Plans using RabbitMQ

All plans running RabbitMQ backend.

`Lemur` and `Tiger` are shared instances and share underlying hardware with other instances. They
are also limited to which CloudAMQP provider resources that can be used. Further information on
availability on each resource page.

Name            | Plan          | Type      | Nodes
----------------|---------------|-----------|------
Little lemur    | lemur         | shared    | -
Tough Tiger     | tiger         | shared    | -
Sassy Squirrel  | squirrel-1    | dedicated | 1
Big Bunny       | bunny-1,3     | dedicated | 1,3
Happy Hare      | hare-1,3      | dedicated | 1,3
Roaring Rabbit  | rabbit-1,3,5  | dedicated | 1,3,5
Power Panda     | panda-1,3,5   | dedicated | 1,3,5
Awesome Ape     | ape-1,3,5     | dedicated | 1,3,5
Heavy Hippo     | hippo-1,3,5   | dedicated | 1,3,5
Loud Lion       | lion-1,3,5    | dedicated | 1,3,5
Raging Rhino    | rhino-1       | dedicated | 1

## Legacy subscription plans

All plans running RabbitMQ backend.

Table below shows deprecated subscription plans for CloudAMQP. Existing plans will still work, but
there will not be possible to create new ones.

Name            | Plan    | Type
----------------|---------| ----
Little lemur    | lemur   | shared
Tough Tiger     | tiger   | shared
Big Bunny       | bunny   | dedicated
Roaring Rabbit  | rabbit  | dedicated
Power Panda     | panda   | dedicated
Awesome Ape     | ape     | dedicated
Heavy Hippo     | hippo   | dedicated
Loud Lion       | lion    | dedicated

[CloudAMQP API plans]: https://docs.cloudamqp.com/index.html#tag/plans
[CloudAMQP plans]: https://www.cloudamqp.com/plans.html
[clustering]: https://lavinmq.com/documentation/clustering
[**LavinMQ**]: https://lavinmq.com/
[**RabbitMQ**]: https://www.rabbitmq.com/

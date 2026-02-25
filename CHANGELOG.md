## 1.43.0 (Unreleased)

NOTES:

* Major refactoring of the API client layer to use unified retry mechanism across all API methods, providing consistent
  error handling and improved reliability.

IMPROVEMENTS:

* ClientLibrary: Refactored all API methods to use unified `callWithRetry()` approach for consistent error handling and retry logic ([#454])
* CLientLibrary: Added rate limit retry with exponential backoff ([#455])
* ClientLibrary: Refactored the genric retry with centralized HTTP status code handling and isolated logic ([#456])
* Integration: Added support for `ap2` Datadog region for resources: ([#457])
  * `cloudamqp_integration_metrics_prometheus`
  * `cloudamqp_integration_log`

DEPENDENCIES:

* Bumped goreleaser/goreleaser-action from 6.4.0 to 7.0.0 ([#458])
* Bumped github.com/hashicorp/terraform-plugin-go from 0.29.0 to 0.30.0 ([#459])
* Bumped hashicorp/setup-terraform from 3 to 4 ([#460])

[#454]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/454
[#455]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/455
[#456]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/456
[#457]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/457
[#458]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/458
[#459]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/459
[#460]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/460

## 1.42.1 (11 Feb, 2026)

IMPROVEMENTS:

* Integration: Updated how metrics filter is created/updated for `cloudamqp_integration_metric_prometheus` ([#451])

[#451]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/451

## 1.42.0 (06 Feb, 2026)

FEATURES:

* Resource: Added support for file provider in `cloudamqp_trust_store` ([#446])

DEPENDENCIES:

* Bumped github.com/hashicorp/terraform-plugin-sdk/v2 from 2.38.1 to 2.38.2 ([#449])

[#446]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/446
[#449]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/449

## 1.41.1 (12 Jan, 2026)

IMPROVEMENTS:

* Resource: Added additional argument "key_id" to `cloudamqp_custom_certificate` ([#448])

[#448]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/448

## 1.41.0 (22 Dec, 2025)

NOTES:

* Added `Claude.md` and `AGENTS.md` files with AI agent instructions.

FEATURES:

* Resource: Added new `cloudamqp_trust_store` resource ([#443])
* Migrated `cloudamqp_node_actions` towards Terraform plugin framework ([#445])

IMPROVEMENTS:

* Resource: Added additional MQTT/SSL settings to `cloudamqp_rabbitmq_configuration` ([#440])
* Resource: Added message interceptos setting to `cloudamqp_rabbitmq_configuration` ([#442])
* Resource: Added cluster actions to `cloudamqp_node_actions` ([#445])

[#440]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/440
[#442]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/442
[#443]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/443
[#445]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/445

## 1.40.0 (11 Dec, 2025)

IMPROVEMENTS:

* ClientLibrary: Added retries for status code 423 and 503 ([#437])
* ClientLibrary: Added retry for custom domain requests ([#438])
* Resource: Added sleep/timeout for `cloudamqp_custom_domain` ([#438])

DEPENDENCIES:

* Bumped github.com/hashicorp/terraform-plugin-framework from 1.16.1 to 1.17.0 ([#435])

[#435]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/435
[#437]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/437
[#438]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/438

## 1.39.0 (02 Dec, 2025)

NOTES:

* Support added for custom certificate, "Bring your own certificate"

FEATURES:

* Resource: Added `cloudamqp_custom_certificate` resource ([#430])

[#430]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/430

## 1.38.3 (21 Nov, 2025)

NOTES:

* Changelog: Cleaned up and make more consistent ([#421])

IMPROVEMENTS:

* Integration: Added region argument to New Relic V3 metrics integration ([#422])
* VCR-test: Renamed fixture filenames containing `::` characters ([#429])

DEPENDENCIES:

* Bumped github.com/hashicorp/terraform-plugin-log from 0.9.0 to 0.10.0 ([#424])
* Bumped golang.org/x/crypto from 0.42.0 to 0.45.0 ([#426])
* Bumped actions/checkout from 5 to 6 ([#427])

[#421]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/421
[#422]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/422
[#424]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/424
[#426]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/426
[#427]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/427
[#429]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/429

## 1.38.2 (13 Nov, 2025)

BUG FIXES:

* VPC-peering: Handled more 400 responses with retry (GCE) ([#423])

[#423]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/423

## 1.38.1 (06 Nov, 2025)

IMPROVEMENTS:

* Docs: Fixed typos in documentation ([#417])
* CI: Avoid parallel test race condition in CI ([#418])
* Docs: Added missing header to Prometheus metrics integrations page ([#419])

DEPENDENCIES:

* Bumped actions/setup-go from 5 to 6 ([#369])

[#369]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/369
[#417]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/417
[#418]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/418
[#419]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/419

## 1.38.0 (27 Oct, 2025)

NOTES:

* More support for Prometheus metrics integrations

FEATURES:

* Integration: Added retention and tags to Cloudwatch logs ([#379])
* Resource: Added `preferred_az` to `cloudamqp_instance` ([#381])
* Resource: Added Oauth2 configuration resource ([#385])
* Integration: Added more support for Prometheus metrics integrations
  * Cloudwatch ([#398])
  * Dynatrace ([#388])
  * Splunk ([#386])
  * Stackdriver ([#410])
* Integration: Added `rabbitmq_dashboard_metrics_format` as option to Datadog prometheus integration ([#413])
* Integration: Added metrics filter settings for prometheus integrations ([#414])

IMPROVEMENTS:

* VCR-test: Speed up testing by enable parallelism ([#403])
* Docs: Updated CloudAMQP API documentation links ([#415])

BUG FIXES:

* Resource drift: Fixed plugins error if instance manually deleted ([#401])
* Resource drift: Fixed for multiple resources if instance or resource manually deleted ([#402])
  * `cloudamqp_maintenance_window`
  * `cloudamqp_privatelink_aws/azure`
  * `cloudamqp_security_firewall`
  * `cloudamqp_vpc_connect`

DEPENDENCIES:

* Bumped github.com/hashicorp/terraform-plugin-framework-validators from 0.18.0 to 0.19.0 ([#396])

[#379]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/379
[#381]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/381
[#385]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/385
[#386]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/386
[#388]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/388
[#396]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/396
[#398]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/398
[#401]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/401
[#402]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/402
[#403]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/403
[#410]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/410
[#413]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/413
[#414]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/414
[#415]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/415

## 1.37.0 (02 Oct, 2025)

FEATURES:

* Added new account action to enable VPC feature ([#371])
* Added support for Prometheus metrics integrations ([#380])
  * Azure monitor
  * Datadog v3
  * New Relic v3

IMPROVEMENTS:

* Client library: Added generic retry ([#362])
* Migrated `cloudamqp_integration_log` towards Terraform plugin framework ([#377])
* Migrated `cloudamqp_integration_metric` towards Terraform plugin framework ([#383])
* Migrated `cloudamqp_vpc` towards Terraform plugin framework ([#363])
* Migrated `cloudamqp_webhook` towards Terraform plugin framework ([#364])
* Handled resource drift `cloudamqp_rabbitmq_configuration` ([#366])
* Default values drift ([#372])
* Updated test cases for VPC ([#384])

DEPENDENCIES:

* Bumped github.com/hashicorp/terraform-plugin-mux from 0.20.0 to 0.21.0 ([#373])
* Bumped github.com/hashicorp/terraform-plugin-sdk/v2 from 2.37.0 to 2.38.1 ([#375], [#378])
* Bumped github.com/hashicorp/terraform-plugin-framework from 1.15.1 to 1.16.1 ([#376], [#382])

[#362]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/362
[#363]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/363
[#364]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/364
[#366]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/366
[#371]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/371
[#372]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/372
[#373]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/373
[#375]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/375
[#376]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/376
[#377]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/377
[#378]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/378
[#380]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/380
[#382]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/382
[#383]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/383
[#384]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/384

## 1.36.0 (29 Aug, 2025)

NOTES:

* Change of underlying maintenance for LavinMQ ([#351])

IMPROVEMENTS:

* Handled deleted resource drift
  * log/metric integration ([#359])
  * alarm, aws eventbridge, instance, notification, webhook ([#361])

DEPENDENCIES:

* Bumped github.com/hashicorp/terraform-plugin-framework from 1.15.0 to 1.15.1 ([#355])
* Bumped actions/checkout from 4 to 5 ([#357])
* Bumped goreleaser/goreleaser-action from 6.3.0 to 6.4.0 ([#358])

[#351]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/351
[#355]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/355
[#357]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/357
[#358]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/358
[#359]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/359
[#361]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/361

## 1.35.1 (Jun 16, 2025)

BUG FIXES:

* Fixed provider schema mismatch ([#349])

[#349]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/349

## 1.35.0 (Jun 16, 2025)

NOTES:

* Initial migration to Terraform Plugin Framework

FEATURES:

* Served the provider through a mux server ([#344])
* Implemented skeleton of provider in the framework plugin ([#344])
* Migrated `cloudamqp_integration_aws_eventbridge` towards Terraform plugin framework ([#344])
* Migrated `cloudamqp_account_actions` towards Terraform plugin framework ([#345])
* Migrated `cloudamqp_rabbitmq_configuration` towards Terraform plugin framework ([#346])

DEPENDENCIES:

* Bumped github.com/cloudflare/circl from 1.6.0 to 1.6.1 ([#342])

[#342]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/342
[#344]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/344
[#345]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/345
[#346]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/346

## 1.34.0 (May 2, 2025)

FEATURES:

* Added new data source `cloudamqp_alarms` to fetch all alarms ([#335])
* Added new data source `cloudamqp_notifications` to fetch all recipients ([#339])

DEPENDENCIES:

* Bumped golang.org/x/net from 0.36.0 to 0.38.0 ([#338])

[#335]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/335
[#338]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/338
[#339]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/339

## 1.33.0 (Apr 7, 2025)

NOTES:

* Changed logging from legacy log to tflog
* Documentation updates

FEATURES:

* Added support to managing maintenance window ([#326])

IMPROVEMENTS:

* Removed legacy log and use tflog package ([#329])
* Docs: LavinMQ supports MQTT and MQTTS ([#327])
* Updated firewall destroy behavior ([#330])
* Docs: Updated extra disk information for Azure ([#331])
* Docs: Updated information about LavinMQ ([#334])

DEPENDENCIES:

* Bumped Go version to 1.24
* Bumped golang.org/x/net from 0.34.0 to 0.36.0 ([#325])
* Bumped crazy-max/ghaction-import-gpg from 6.2.0 to 6.3.0 ([#332])
* Bumped goreleaser/goreleaser-action from 6.2.1 to 6.3.0 ([#333])

[#325]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/325
[#326]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/326
[#327]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/327
[#329]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/329
[#330]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/330
[#331]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/331
[#332]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/332
[#333]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/333
[#334]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/334

## 1.32.3 (Mar 3, 2025)

IMPROVEMENTS:

* Docs: Updated VPC connect resource information when used for Azure ([#311])
* Docs: Added guide page about Marketplace migration ([#312])
* Added include auto delete queues to metric integration ([#323])

BUG FIXES:

* Fixed throwing error when external VPC identifier cannot be found. ([#310])

[#310]: https://github.com/cloudamqp/terraform-provider-cloudamqp/issues/310
[#311]: https://github.com/cloudamqp/terraform-provider-cloudamqp/issues/311
[#312]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/312
[#323]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/323

## 1.32.2 (Dec 20, 2024)

IMPROVEMENTS:

* Enabled import of VPC peering resource for GCP ([#308])
* Enabled import of VPC peering resource for AWS ([#309])

[#308]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/308
[#309]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/309

## 1.32.1 (Oct 28, 2024)

BUG FIXES:

* Fixed incorrect schemas in plugin data sources ([#300])

[#300]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/300

## 1.32.0 (Sep 4, 2024)

FEATURES:

* Added support to upgrade LavinMQ instances ([#296])

[#296]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/296

## 1.31.0 (Aug 19, 2024)

FEATURES:

* Added support to specify RabbitMQ version when upgrading ([#295])
* Added support to use data source when upgrading RabbitMQ version ([#295])

[#295]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/295

## 1.30.1 (Jul 30, 2024)

IMPROVEMENTS:

* Docs: Added notification example for Slack ([#287])
* Added internal hostname information to nodes data source ([#289])
* Added availability zone information to nodes data source ([#290])

[#287]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/287
[#289]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/289
[#290]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/291

## 1.30.0 (Jun 10, 2024)

NOTES:

* Github CI workflow with Go VCR basic resource testing
* Go-API client library imported into provider and removed external dependency
* Terraform Plugin SDK v2

FEATURES:

* Added Go VCR basic resource testing that extends acceptance test with stored fixtures ([#257])
* Updated Terraform Plugin SDK to V2 ([#261])

IMRPOVEMENTS:

* Added support for updating webhook resource ([#268])
* Added configurable retries for webhook resource ([#268])
* Updated integration resource docs for Datadog tags ([#277])
* Imported Go-API client library with history ([#282])
* Posted Go-API import modifications ([#284])

[#257]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/257
[#261]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/261
[#268]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/268
[#277]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/277
[#282]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/282
[#284]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/284

## 1.29.5 (Apr 04, 2024)

IMPROVEMENTS:

* Fixed link to instance regions guide from instances page ([#263])
* Added information on how to use Message Broker HTTP API ([#264])
* Added handling of "creating/deleting" notice alarm ([#265])

[#263]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/263
[#264]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/264
[#265]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/265

## 1.29.4 (Feb 15, 2024)

IMPROVEMENTS:

* Added optional responders argument for OpsGenie recipient ([#258])

[#258]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/258

## 1.29.3 (Jan 26, 2024)

IMPROVEMENTS:

* Added support for Azure monitor log integration ([#254])
* Added support for signl4 alarms recipient ([#255])

[#254]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/254
[#255]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/255

## 1.29.2 (Jan 17, 2024)

IMPROVEMENTS:

* Added support for the Coralogix log integration ([#253])

[#253]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/253

## 1.29.1 (Dec 21, 2023)

BUG FIXES:

* Fixed PrivateLink/Private Service Connect import ([#250])

[#250]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/250

## 1.29.0 (Dec 18, 2023)

NOTES:

* Updated the API wrapper (go-api) dependency with newer version (1.15.0)

FEATURES:

* Added resource that invoke account actions. ([#231])
* Added new generic resource for VPC Connect ([#240])
  * Enables GCP Private Service Connect
  * Handles AWS PrivateLink
  * Handles Azure PrivateLink
* Added configurable retries for plugin resources ([#241])
* Added configurable retry when reading PrivateLink information ([#246])
* Added configurable retry for GCP VPC peering ([#247])

IMPROVEMENTS:

* Updated and cleaned up samples ([#235])
* Removed default RMQ version request when version left out ([#237])
* Handled gone VPC resource ([#238])

[#231]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/231
[#235]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/235
[#237]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/237
[#238]: https://github.com/cloudamqp/terraform-provider-cloudamqp/issues/238
[#240]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/240
[#241]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/241
[#246]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/246
[#247]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/247

## 1.28.0 (Sep 27, 2023)

NOTES:

* Updated the API wrapper (go-api) dependency with newer version (1.12.4)

FEATURES:

* Added support for "copy settings" from another instance when creating a new. ([#218])
* Configured wait on GCP Peering status. ([#228])

[#218]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/218
[#228]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/228

## 1.27.1 (Sep 08, 2023)

NOTES:

* Updated the API wrapper (go-api) dependency with newer version (1.12.3)

IMPROVEMENTS:

* Cleanup RabbitMQ configuration resource ([#215])
* Added ForceNew to resources with cloudamqp_instance dependency ([#222])

BUG FIXES:

* Indirect handle managed required plugins failing to be destroyed. ([#227])

[#215]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/215
[#222]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/222
[#227]: https://github.com/cloudamqp/terraform-provider-cloudamqp/issues/227

## 1.27.0 (Jun 12, 2023)

NOTES:

* New provider configuration option to enable faster instance destroy.

FEATURES:

* Added assume role authentication for CloudWatch metrics integration ([#208])
* Enable faster instance destroy options in provider configuration ([#209])

IMPROVEMENTS:

* Added missing `Happy Hare` plan to the docs ([#206])
* Updated firewall rules, PrivateLink and VPC Peering documentation ([#207])
* Allowed queue_index_embed_msgs_below to be set to 0 in RabbitMQ configuration.

[#206]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/206
[#207]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/207
[#208]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/208
[#209]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/209

## 1.26.2 (May 12, 2023)

NOTES:

* Updated the API wrapper (go-api) dependecy with newer version (1.12.2)

IMPROVEMENTS:

* Indirect improvements with retry when deleting firewall settings

## 1.26.1 (May 05, 2023)

NOTES:

* Updated the API wrapper (go-api) dependecy with newer version (1.12.1)

BUG FIXES:

* Fixed underlying issue with validation error response

## 1.26.0 (May 03, 2023)

NOTES:

* Updated Go version with newer version (1.20)
* Updated the API wrapper (go-api) dependecy with newer version (1.12.0)

FEATURES:

* Added support to use the API backend to validate plans and regions ([#201])

IMPROVEMENTS:

* Added support to configure cluster_partition_handling ([#200])

BUG FIXES:

* Added missing options attribute for notification data source ([#199])

[#199]: https://github.com/cloudamqp/terraform-provider-cloudamqp/issues/199
[#200]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/200
[#201]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/201

## 1.25.0 (Apr 03, 2023)

NOTES:

* Resize disk with extra disk resource supports more platforms (GCE, Azure)
* Updated the API wrapper (go-api) dependecy with newer version (1.11.1)

FEATURES:

* Resize disk with for more platforms and using new optional argument allow_downtime ([#194])

[#194]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/194

## 1.24.2 (Mar 30, 2023)

NOTES:

* Fixed issues introduced in previous version 1.24.1

BUG FIXES:

* Stackdriver optional arguments assignments ([#198])

[#198]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/198

## 1.24.1 (Mar 14, 2023)

BUG FIXES:

* Converted optional queue/vhost to correct JSON fields for metrics integration

## 1.24.0 (Mar 07, 2023)

NOTES:

* Updated the API wrapper (go-api) dependecy with newer version (1.11.0)

FEATURES:

* Added support for AWS EventBridge integration ([#192])

DEPENDENCIES:

* Bumped github.com/hashicorp/go-getter from 1.6.1 to 1.7.0 ([#187])
* Bumped golang.org/x/net from 0.0.0-20210326060303-6b1517762897 to 0.7.0 ([#190])
* Bumped golang.org/x/crypto from 0.0.0-20210921155107-089bfa567519 to 0.1.0 ([#191])

[#187]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/187
[#190]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/190
[#191]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/191
[#192]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/192

## 1.23.0 (Jan 26, 2023)

NOTES:

* Enabled creating shared subscription beta plan for [LavinMQ](https://lavinmq.com/).

IMPROVEMENTS:

* Added LavinMQ lemming ([#182])

[#182]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/182

## 1.22.1 (Jan 18, 2023)

IMPROVEMENTS:

* Updated subscription plan validation with new plans `hare-1` and `hare-3`.

## 1.22.0 (Jan 09, 2023)

NOTES:

* Alarm notification recipients options parameter.

IMPROVEMENTS:

* Optional options key-value pair argument for alarm notification/recipient ([#185])

[#185]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/185

## 1.21.0 (Dec 21, 2022)

NOTES:

* Stackdiver integrations (log & metric) to use raw Google Service Account key credentials.

IMPROVEMENTS:

* Updated Stackdriver integrations to use raw Google Service Account key credentials ([#184])

BUG FIXES:

* Excluded additional parameters (tags, queue_allowlist, vhost_allowlist) from integrations when not used ([#184])

[#184]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/184

## 1.20.2 (Dec 14, 2022)

NOTES:

* Updated the API wrapper (go-api) dependecy with newer version (1.10.2)

IMPROVEMENTS:

* Added configurable sleep and timeout for firewall configuration ([#183])

[#183]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/183

## 1.20.1 (Dec 07, 2022)

NOTES:

* Updated the API wrapper (go-api) dependecy with newer version (1.10.1)
* Extended response handling for read/update RabbitMQ configuration

IMPROVEMENTS:

* Added configurable sleep and timeout for RabbitNQ configuration

## 1.20.0 (Oct 24, 2022)

NOTES:

* Updated the API wrapper (go-api) dependecy with newer version (1.10.0)
* Added support for PrivateLink for AWS and Azure

FEATURES:

* Added support for PrivateLink for AWS and Azure ([#173])

BUG FIXES:

* Updated minimum value of heartbeat to 0 ([#176])
* Missing required splunk integration parameter ([#177])

[#173]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/173
[#176]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/176
[#177]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/177

## 1.19.3 (Oct 07, 2022)

NOTES:

* Updated the API wrapper (go-api) dependecy with newer version (1.9.2)
* Added support for retry VPC peering and wait for status

BUG FIXES:

* Added additional computed fields to plugins resources ([#170])

[#170]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/170

## 1.19.2 (Sep 14, 2022)

NOTE:

* Updated the API wrapper (go-api) dependecy with newer version (1.9.1).
* Now supports asynchronous request for plugin/community actions. Solved issues when enabling multiple plugins.

IMPROVEMENTS:

* Added CIDR address validation ([#168])
* Updated workflow for updating RabbitMQ configuration
([#166])

[#166]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/166
[#168]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/168

## 1.19.1 (Aug 04, 2022)

IMPROVEMENTS:

* Added support to disable consumer_timeout for RabbitMQ configuration.
* Excluded nodes argument when using shared instance plan.

## 1.19.0 (Jul 01, 2022)

NOTE:

* Updated the API wrapper (go-api) dependecy with newer version (1.9.0)
* Updated goutils dependecy with newer version (1.1.1)

FEATURES:

* Added support for resize disk ([#162])

IMPROVEMENTS:

* Updated nodes data source with original and additional disk sizes ([#162])

[#162]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/162

## 1.18.0 (Jun 08, 2022)

NOTE:

* Updated the API wrapper (go-api) dependency with newer version (1.8.1)
* Updated go-getter dependency with newer version (1.6.1)

FEATURES:

* Added support for updating RabbitMQ config ([#150])
* Added support for invoking node actions ([#150])

IMPROVEMENTS:

* Updated wrong information in documentation for VPC peering.

[#150]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/150

## 1.17.2 (May 27, 2022)

IMPROVEMENTS:

* Added `flow` as supported alarm type.

## 1.17.1 (May 24, 2022)

IMPROVEMENTS:

* Added `reminder_interval` schema argument for alarms.

## 1.17.0 (May 24, 2022)

NOTE:

* Updated the API wrapper (go-api) dependency with newer version (1.8.0)
* Configurable timeout/sleep for VPC peering, avoids firewall configuration blocking VPC peering.

FEATURES:

* Added support to upgrade to latest possible versions for RabbitMQ and Erlang ([#151])

IMPROVEMENTS:

* Added configurable timeout/sleep for accept/remove VPC peering. ([#153])

[#151]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/151
[#153]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/153

## 1.16.0 (May 09, 2022)

NOTE:

* Updated the API wrapper (go-api) dependency with newer version (1.6.0)
* Introducing managed VPC resource to decouple VPC from instance. ([#148])
* To avoid breaking changes
  * It's still possible to create VPC from instance with vpc_subnet, but is discouraged.
  * Default behaviour for instance is still to delete associated VPC.
  * To keep managed VPC, set attribute *keep_associated_vpc = true* on each instance resource. This will override the default behaviour when deleting an instance.

FEATURES:

* Added support for managed VPC resource.
* Added list on all available standalone VPC for an account.
* Added multiple attribute (vpc_id and instance_id) to fetch VPC information.
* Added multiple attribute (vpc_id and instance_id) to handle VPC peering.
* Added documentations for managed VPC resources and guide

IMPROVEMENTS:

* Added keep_associated_vpc attribute for instance resource

DEPRECATED:

* data_source/vpc_gcp_info, intance_id use vpc_id instead
* data_source/vpc_info, instance_id use vpc_id instead
* resource/instance, vpc_subnet create managed VPC instead
* resource/vpc_gcp_peering, intance_id use vpc_id instead
* resource/vpc_peering, intance_id use vpc_id instead

[#148]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/148

## 1.15.3 (Apr 06, 2022)

IMPROVEMENTS:

* Added support for Scalyr log integrations ([#147])

[#147]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/147

## 1.15.2 (Mar 29, 2022)

IMPROVEMENTS:

* Added new attribute, value_calculation, to alarms ([#138])
* Added support for `CLOUDAMQP_BASEURL` in provider, made testing easier ([#143])

BUG FIXES:

* Corrected validation for firewall rule attributes ([#141])

[#138]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/138
[#141]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/141
[#143]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/143

## 1.15.1 (Dec 21, 2021)

NOTE:

* Updated the API wrapper (go-api) dependency with newer version (1.5.4)

IMPROVEMENTS:

* Removed peer_subnet as schema attribute from VPC GCP peering
* Removed formatting response data for firewall rules
* Indirect multiple retry functionality to create and update firewall rules
* Updated VPC GCP peering documentation

## 1.15.0 (Dec 20, 2021)

NOTE:

* Updated the API wrapper (go-api) dependency with newer version (1.5.3)

FEATURES:

* Added VPC information for Google Cloud Platform ([#131])
* Added VPC peer configuration for Google Cloud Platform ([#74])

[#74]: https://github.com/cloudamqp/terraform-provider-cloudamqp/issues/74
[#131]: https://github.com/cloudamqp/terraform-provider-cloudamqp/issues/131

## 1.14.0 (Dec 3, 2021)

Note:

* Updated the API wrapper (go-api) dependency with newer version (1.5.2) ([#129])

IMPROVEMENTS:

* Added `STREAM`, `STREAM_SSL` as supported firewall services ([#128])

[#128]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/128
[#129]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/129

## 1.13.0 (Nov 15, 2021)

IMPROVEMENTS:

* Added attribute `host_internal` to instance resource ([#127])
* Made the attribute `host` always return the external hostname ([#127])
* Set `ForceNew` on `region` in instance resource ([#122]) (**Note**: when forcing a region change, the previous instance will be destroyed and a new one created)

[#127]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/127
[#122]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/122

## 1.12.0 (Oct 29, 2021)

Note:

* Created instance in existing VPC

FEATURES:

* Added support for creating instance in existing VPC
* Added config for auto generated release notes

IMPROVEMENTS:

* Updated typos in documentation
* Removed unused attributes from instance resource

## 1.11.0 (Oct 06, 2021)

Note:

* Updated the API wrapper (go-api) dependency with newer version (1.5.1)
* Updated go to version 1.17
* Updated Terraform Plugin SDK to version 1.17.2

FEATURES:

* Added resource for account
* Added resource for custom domain

IMPROVEMENTS:

* Updated internal handling of provider version number
* Updated handling of number of nodes
* Indirectly improved community plugin request that can fail due to backend being busy (go-api v1.5.0)

BUG FIXES:

* Added missing schema attributes for instance data source

## 1.10.0 (Sep, 20 2021)

Note: Updated the API wrapper (go-api) dependency with newer version (1.4.0)

IMPORVEMENTS:

* Indirectly improved common request that can fail due to backend being busy.

## 1.9.4 (Sep, 16 2021)

Note: Re-release 1.9.3 with missing information

## 1.9.3 (Sep, 16 2021)

IMPROVEMENTS:

* Validated schema attributes when reading response from API calls.
* Optional attributes changed to computed for data sources.
* Alarm: Populated alarm_id with correct identifier.
* Documentation: Added identifier attribute reference to all resources and data source.

BUG FIXES:

* Added configured attribute to nodes data source.
* Instance: Updated switch statement to get correct plan type.

## 1.9.2 (May 21, 2021)

IMPROVEMENTS:

* Added HTTPS as a supported firewall service
* Allowed MS Teams to be used as recipient type

BUG FIXES:

* Import of log intergration, corrected the identifiers needed to fetch integration
* Import of metrics integration, corrected the identifier needed to fetch integration
* Displayed hostname and internal hostname correctly

## 1.9.1 (Feb 5, 2021)

IMPROVEMENTS:

* Added validation of plan name before execution

## 1.9.0 (Feb 2, 2021)

NOTES:

* Enabled switching to new subscription plans. See documentation for more information.

## 1.8.6 (Dec 9, 2020)

BUG FIX

* Removed default values from attributes with computed/optional combination.

## 1.8.5 (Dec 8, 2020)

BUG FIXES:

* Failed to fetch default RMQ version from CloudAMQP API.

IMPROVEMENTS:

* Updated CHANGELOG with missing releases.
* Cleaned up OS/arch combinations, reverted back during initial publish to Terraform registry.

## 1.8.4 (Nov 18, 2020)

NOTES:

* Cleaned up language used, deprecated white-/blacklist

IMPROVEMENTS:

* Deprecated white-/blacklist, added allow-/blocklist.

## 1.8.3 (Nov 12, 2020)

IMPROVEMENTS:

* Removed some OS/arch combinations.

## 1.8.2 (Nov 12, 2020)

NOTES:

* Terraform Registry: New releases automatically updates registry with the help of GitHub actions.
* Webhook added already Oct 6, but no release until Nov 12.

FEATURES:

* Added support for webhook implementation.

IMPROVEMENTS:

* Using version 1.3.4 of wrapper API (go-api).
* Updated instance to wait until all nodes are finished configuring after update.

## 1.8.1 (Unreleased)

BUG FIXES:

* Removed invalid attribute validation, caused log integration to fail.

## 1.8.0 (Unreleased)

NOTES:

* Initial release for Terraform Provider Development Program

## 1.7.3 (Jul 7, 2020)

IMPROVEMENTS:

* README information about where to find instance info
* Firewall: Handled updates in wrapper API (go-api 1.3.3), waiting on firewall changes.
* Metrics integrations: Enabled contributed Stackdriver functionality.

## 1.7.2 (Jun 15, 2020)

IMPROVEMENTS:

* Updated install instructions.

BUG FIXES:

* Renamed GNUmakefile to correct naming, due to missing target for make.

## 1.7.1 (Jun 12, 2020)

IMPROVEMENTS:

* Addded `no_default_alarms` to `cloudamqp_instance`.
* Updated Terraform.io documents.

## 1.7.0 (Jun, 8, 2020)

NOTES:

* Resolved initial review feedback

IMPROVEMENTS:

* Naming convetion on data source and resource files
* Updated Makefile to GNUMakefile
* Added script for provider integration
* Re-enabled vendor folder
* Terraform.io website documentation
* Updated changelog to make it readable for release bot
* Updated samples with use of variables
* Double checked data sources and resource for required, optional, computed and sensitive properties.
* Triggered read resource information when the resource has been updated.
* Lint: naming convention
* Lint: error checks

BUG FIXES:

* Underlying error messages for shared instances.

## 1.6.0 (Apr 27, 2020)

FEATURES:

* **New Resource**: resource_integration_log - Log integration to third party service
* **New Resource**: resource_integration_metric - Metric integration to third party service
* Acceptence test for majority of data sources and resources

IMPROVEMENTS:

* Instance: Merged contribution to handle plan changes between shared and dedicated
* VPC Peering: Peering request status information

## 1.5.0 (Mar 10, 2020)

FEATURES:

* **New Data Source**: data_source_instance

IMPROVEMENTS:

* Added validation of message type attribute when populating alarm schema
* Message type key in create and update for alarm
* Added validation of attributes when populating instance schema

## 1.4.1 (Mar 6, 2020)

BUG FIXES:

* Missing required message type fields for queue alarms.

## 1.4.0 (Feb 27, 2020)

NOTES:

* Underlying API changes required updated payload for alarm and notification.

IMPROVEMENTS:

* Alarm: Additional schema attributes [enabled]
* Alarm: Rename of schema attribute notification_ids -> recipients
* Notifications: Additional schema attributes [name]

## 1.3.2 (Feb 18, 2020)

IMPROVEMENTS:

* Extract host and vhost information when creating new instance.

BUG FIXES:

* Updated go-api dependency with minor regex fix (second try).

## 1.3.1 (Feb 17, 2020)

BUG FIXES:

* Updated go-api dependency with minor regex fix.

## 1.3.0 (Jan 16, 2019)

FEATURES:

* Changed depedenacy mangement from package to modules.

IMPROVEMENTS:

* Additional information about security group added to data_source_vpc_info.go
* Added .exe extension on Windows release for cross-compile

## 1.2.3 (Dec 16, 2019) // Unreleased

IMPROVEMENTS:

* Changed package path for ldflags to get correct version

## 1.2.2 (Dec 13, 2019)

IMPROVEMENTS:

* Exposed computed reference of host and vhost for CloudAMQP instances

## 1.2.1 (Dec 9, 2019)

IMPROVEMENTS:

* Added debug logging through out data sources and resources
* Added validation of identifiers before internal assigning them
* Extended release support of cross compile GOOS and GOARCH

## 1.2.0 (Nov 26, 2019)

FEATURES:

* **New Data Source:** data_source_credentials - Extract credentials
* **New Resource:** resource_security_firewall.go - Firewall configuration
* **New Resource:** resource_plugins.go - Configurable Rabbit MQ plugins
* **New Resource:** resource_community_plugins.go - Configurable community plugins
* **New Resource:** resource_vpc_peering.go - Enable VPC support for AWS instances

IMPROVEMENTS:

* Restructured and moved data source, resource and provider etc. files into cloudamqp sub-folder
* Upgraded to Terraform version 0.12.9
* Versioning on compiled provider.
* Configurable Rabbit MQ version
* Added validation of functions (alarm types, notifications types, firewall settings and ports).

## 1.1.3(Unreleased)

IMPROVEMENTS:

* Updated installation part in Readme. (Merge pull request [#30](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/30))
* Makefile compile issue for MacBook. (Merge pull request [#29](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/29))

## 1.1.2 (Unreleased)

IMPROVEMENTS:

* Updated of documentation. (Merge pull request [#28](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/28))
* Install procedure for Linux. (Merge pull request [#27](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/27))

## 1.1.1 (Unreleased)

IMPROVEMENTS:

* Updated examples to match Terraform 0.12.* (Merge pull request [#24](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/24))

## 1.1.0 (Oct 8, 2019)

FEATURES:

* **New Resource**: resource_alarm.go - Configurable alarms for different metrics
* **New Resource**: resource_notifications.go - Configurable notifications endpoints and recipients

IMPROVEMENTS:

* Tags on instances
* Enabled Terraform import on resources
* Cross-compile provider release as make command

## 1.0.0 (Sep 17, 2019)

NOTES:

* Initial release

IMPROVEMENTS:

* Added tags to instance resource
* Updated tags type
* Updated go-api branch depedency

## 0.4.1 (Unreleased)

IMPROVEMENTS:

* Documentation updates. (Merge pull request [#10](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/10))

## 0.4.0 (Unreleased)

FEATURES:

* **New Resource**: resource_alarm.go - Configurable alarms to monitoring metrics
* **New Resource**: resource_notifications.go - Configurable notifications endpoint and recipients

IMPROVEMENTS:

* Updated API endpoints
* Documents and examples

## 0.3.2 (Unreleased)

IMPROVEMENTS:

* Vendor cleanup. (Merge pull request [#2](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/2))
* Support for Terraform 0.12. (Merge pull request [#16](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/16))

## 0.3.1 (Unreleased)

IMPROVEMENTS

* Additional .gitignore updates. (Merge pull request [#8](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/8))
* Lemur example for newer users
* Updated and extended documentation.

## 0.3.0 (Unreleased)

FEATURES:

* Readme.md

IMPROVEMENTS:

* Updated dependencies
* Generic API for resources

## 0.2.0 (Unreleased)

IMPROVEMENTS:

* Instance update and delete
* Added support for vpc_subnet, nodes and rmq_versions
* Url and apikey set as sensitive

## 0.1.0 (Unreleased)

NOTES:

* Initial commit

FEATURES:

* **New Resource**: resource_instance.go - Main resource
* Basic provider
* Makefile, logic for clean, dependenecy updates, build and install.

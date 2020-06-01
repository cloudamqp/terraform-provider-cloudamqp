## 1.7.0 (Unreleased)

NOTES:

* Resolved initial review feedback

IMPROVEMENTS:

* Naming convetion on data source and resource files
* Updated Makefile to GNUMakefile
* Added script for provider integration
* Re-enabled vendor folder
* Terraform.io website documentation
* Updated changelog to make it readble for release bot
* Updated samples with use of variables
* Double checked data sources and resource for required, optional, computed and sensitive properties.
* Trigger read resource information when the resource has been updated.

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

* Validating message type attribute when populating alarm schema
* Message type key in create and update for alarm
* Validating attributes when populating instance schema

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

* Expose computed reference of host and vhost for CloudAMQP instances

## 1.2.1 (Dec 9, 2019)

IMPROVEMENTS:

* Added debug logging through out data sources and resources
* Validation of identifiers before internal assigning them
* Extended release support of cross compile GOOS and GOARCH

## 1.2.0 (Nov 26, 2019)

FEATURES:

* **New Data Source:** data_source_credentials - Extract credentials
* **New Resource:** resource_security_firewall.go - Firewall configuration
* **New Resource:** resource_plugins.go - Configurable Rabbit MQ plugins
* **New Resource:** resource_community_plugins.go - Configurable community plugins
* **New Resource:** resource_vpc_peering.go - Enable VPC support for AWS instances

IMPROVEMENTS:

* Restructure and move data source, resource and provider etc. files into cloudamqp sub-folder
* Upgrade to Terraform version 0.12.9
* Versioning on compiled provider.
* Configurable Rabbit MQ version
* Validation functions (alarm types, notifications types, firewall settings and ports).

## 1.1.3(Unreleased)

IMPROVEMENTS:

* Updated installation part in Readme. (Merge pull request [#30](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/30))
* Makefile compile issue for MacBook. (Merge pull request [#29](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/29))

## 1.1.2 (Unreleased)

IMPROVEMENTS:

* Update of documentation. (Merge pull request [#28](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/28))
* Install procedure for Linux. (Merge pull request [#27](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/27))

## 1.1.1 (Unreleased)

IMPROVEMENTS:

* Update examples to match Terraform 0.12.* (Merge pull request [#24](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/24))

## 1.1.0 (Oct 8, 2019)

FEATURES:

* **New Resource**: resource_alarm.go - Configurable alarms for different metrics
* **New Resource**: resource_notifications.go - Configurable notifications endpoints and recipients

IMPROVEMENTS:

* Tags on instances
* Enable Terraform import on resources
* Cross-compile provider release as make command

## 1.0.0 (Sep 17, 2019)

NOTES:

* Initial release

IMPROVEMENTS:

* Added tags to instance resource
* Updated tags type
* Update go-api branch depedency

## 0.4.1 (Unreleased)

IMPROVEMENTS:

* Documentation updates. (Merge pull request [#10](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/10))

## 0.4.0 (Unreleased)

FEATURES:

* **New Resource**: resource_alarm.go - Configurable alarms to monitoring metrics
* **New Resource**: resource_notifications.go - Configurable notifications endpoint and recipients

IMPROVEMENTS:

* Update API endpoints
* Documents and examples

## 0.3.2 (Unreleased)

IMPROVEMENTS:

* Vendor cleanup. (Merge pull request [#2](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/2))
* Support for Terraform 0.12. (Merge pull request [#16](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/16))

## 0.3.1 (Unreleased)

IMPROVEMENTS

* Additional .gitignore updates. (Merge pull request [#8](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/8))
* Lemur example for newer users
* Update and extending documentation.

## 0.3.0 (Unreleased)

FEATURES:

* Readme.md

IMPROVEMENTS:

* Update dependencies
* Generic API for resources

## 0.2.0 (Unreleased)

IMPROVEMENTS:

* Instance update and delete
* Support for vpc_subnet, nodes and rmq_versions
* Url and apikey set as sensitive

## 0.1.0 (Unreleased)

NOTES:

* Initial commit

FEATURES:

* **New Resource**: resource_instance.go - Main resource
* Basic provider
* Makefile, logic for clean, dependenecy updates, build and install.

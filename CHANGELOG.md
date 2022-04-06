## 1.15.3 (Apr 06, 2022)

IMPROVEMENTS:

* Add support for Scalyr log integrations ([#147]
https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/147)

## 1.15.2 (Mar 29, 2022)

IMPROVEMENTS:

* Add new attribute, value_calculation, to alarms ([#138]https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/138)
* Add support for `CLOUDAMQP_BASEURL` in provider, make testing easier ([#143](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/143))

BUG FIXES:

* Correct validation for firewall rule attributes ([#141]https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/141)

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

* Added VPC information for Google Cloud Platform ([#131](https://github.com/cloudamqp/terraform-provider-cloudamqp/issues/131))
* Added VPC peer configuration for Google Cloud Platform ([#74](https://github.com/cloudamqp/terraform-provider-cloudamqp/issues/74))

## 1.14.0 (Dec 3, 2021)

Note:

* Updated the API wrapper (go-api) dependency with newer version (1.5.2) ([#129](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/129))

IMPROVEMENTS:

* Add `STREAM`, `STREAM_SSL` as supported firewall services ([#128](https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/128))

## 1.13.0 (Nov 15, 2021)

IMPROVEMENTS:

* Add attribute `host_internal` to instance resource ([#127])
* Make the attribute `host` always return the external hostname ([#127])
* Set `ForceNew` on `region` in instance resource ([#122]) (**Note**: when forcing a region change, the previous instance will be destroyed and a new one created)

[#127]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/127
[#122]: https://github.com/cloudamqp/terraform-provider-cloudamqp/pull/122

## 1.12.0 (Oct 29, 2021)

Note:

* Create instance in existing VPC

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
* Indirect improve community plugin request that can fail due to backend being busy (go-api v1.5.0)

BUG FIXES:

* Added missing schema attributes for instance data source

## 1.10.0 (Sep, 20 2021)

Note: Update the API wrapper (go-api) dependency with newer version (1.4.0)

IMPORVEMENTS:

* Indirect improve common request that can fail due to backend being busy.

## 1.9.4 (Sep, 16 2021)

Note: Re-release 1.9.3 with missing information

## 1.9.3 (Sep, 16 2021)

IMPROVEMENTS:

* Validate schema attributes when reading response from API calls.
* Optional attributes changed to computed for data sources.
* Alarm: Populate alarm_id with correct identifier.
* Documentation: Added identifier attribute reference to all resources and data source.

BUG FIXES:

* Added configured attribute to nodes data source.
* Instance: Updated switch statement to get correct plan type.

## 1.9.2 (May 21, 2021)

IMPROVEMENTS:

* Add HTTPS as a supported firewall service
* Allow MS Teams to be used as recipient type

BUG FIXES:

* Import of log intergration, corrected the identifiers needed to fetch integration
* Import of metrics integration, corrected the identifier needed to fetch integration
* Display of hostname and internal hostname

## 1.9.1 (Feb 5, 2021)

IMPROVEMENTS:

* Validation of plan name before execution

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

* Cleanup of language used, deprecate white-/blacklist

IMPROVEMENTS:

* Deprecated white-/blacklist, added allow-/blocklist.

## 1.8.3 (Nov 12, 2020)

IMPROVEMENTS:

* Remove some OS/arch combinations.

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
* Firewall: Handling updates in wrapper API (go-api 1.3.3), waiting on firewall changes.
* Metrics integrations: Enable contributed Stackdriver functionality.

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
* Updated changelog to make it readble for release bot
* Updated samples with use of variables
* Double checked data sources and resource for required, optional, computed and sensitive properties.
* Trigger read resource information when the resource has been updated.
* Lint: naming convetion
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

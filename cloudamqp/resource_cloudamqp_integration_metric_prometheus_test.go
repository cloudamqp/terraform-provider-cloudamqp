package cloudamqp

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccIntegrationMetricPrometheusNewRelicV3_Basic: Add NewRelic v3 prometheus metric integration and import.
func TestAccIntegrationMetricPrometheusNewRelicV3_Basic(t *testing.T) {
	var (
		fileNames                      = []string{"instance", "integrations/metrics/integration_metric_prometheus_newrelic_v3"}
		instanceResourceName           = "cloudamqp_instance.instance"
		prometheusNewRelicResourceName = "cloudamqp_integration_metric_prometheus.newrelic_v3"

		params = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusNewRelicV3_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"NewRelicApiKey": "NEWRELIC_APIKEY",
			"NewRelicTags":   "env=test,region=us1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusNewRelicResourceName, "newrelic_v3.#", "1"),
					resource.TestCheckResourceAttr(prometheusNewRelicResourceName, "newrelic_v3.0.tags", params["NewRelicTags"]),
				),
			},
			{
				ResourceName:      prometheusNewRelicResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, prometheusNewRelicResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusNewRelicV3_WithoutTags: Test NewRelic v3 prometheus integration without optional tags.
func TestAccIntegrationMetricPrometheusNewRelicV3_WithoutTags(t *testing.T) {
	var (
		fileNames                      = []string{"instance", "integrations/metrics/integration_metric_prometheus_newrelic_v3_notags"}
		instanceResourceName           = "cloudamqp_instance.instance"
		prometheusNewRelicResourceName = "cloudamqp_integration_metric_prometheus.newrelic_v3_notags"

		params = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusNewRelicV3_WithoutTags",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"NewRelicApiKey": "NEWRELIC_APIKEY",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusNewRelicResourceName, "newrelic_v3.#", "1"),
					resource.TestCheckResourceAttr(prometheusNewRelicResourceName, "newrelic_v3.0.tags", ""),
				),
			},
			{
				ResourceName:      prometheusNewRelicResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, prometheusNewRelicResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusNewRelicV3_Update: Test updating NewRelic v3 prometheus integration.
func TestAccIntegrationMetricPrometheusNewRelicV3_Update(t *testing.T) {
	var (
		fileNames                      = []string{"instance", "integrations/metrics/integration_metric_prometheus_newrelic_v3"}
		instanceResourceName           = "cloudamqp_instance.instance"
		prometheusNewRelicResourceName = "cloudamqp_integration_metric_prometheus.newrelic_v3"

		paramsCreate = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusNewRelicV3_Update",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"NewRelicApiKey": "NEWRELIC_APIKEY",
			"NewRelicTags":   "env=test,region=us1",
		}

		paramsUpdate = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusNewRelicV3_Update",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"NewRelicApiKey": "NEWRELIC_APIKEY",
			"NewRelicTags":   "env=prod,region=eu1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsCreate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsCreate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusNewRelicResourceName, "newrelic_v3.0.tags", paramsCreate["NewRelicTags"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsUpdate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusNewRelicResourceName, "newrelic_v3.0.tags", paramsUpdate["NewRelicTags"]),
				),
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusDatadogV3_Basic: Add Datadog v3 prometheus metric integration and import.
func TestAccIntegrationMetricPrometheusDatadogV3_Basic(t *testing.T) {
	var (
		fileNames                     = []string{"instance", "integrations/metrics/integration_metric_prometheus_datadog_v3"}
		instanceResourceName          = "cloudamqp_instance.instance"
		prometheusDatadogResourceName = "cloudamqp_integration_metric_prometheus.datadog_v3"

		params = map[string]string{
			"InstanceName":  "TestAccIntegrationMetricPrometheusDatadogV3_Basic",
			"InstanceID":    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":  "bunny-1",
			"DatadogApiKey": "DATADOG_APIKEY",
			"DatadogRegion": "us1",
			"DatadogTags":   "env=test,region=us1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.#", "1"),
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.0.region", params["DatadogRegion"]),
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.0.tags", params["DatadogTags"]),
				),
			},
			{
				ResourceName:      prometheusDatadogResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, prometheusDatadogResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusDatadogV3_WithoutTags: Test Datadog v3 prometheus integration without optional tags.
func TestAccIntegrationMetricPrometheusDatadogV3_WithoutTags(t *testing.T) {
	var (
		fileNames                     = []string{"instance", "integrations/metrics/integration_metric_prometheus_datadog_v3_notags"}
		instanceResourceName          = "cloudamqp_instance.instance"
		prometheusDatadogResourceName = "cloudamqp_integration_metric_prometheus.datadog_v3_notags"

		params = map[string]string{
			"InstanceName":  "TestAccIntegrationMetricPrometheusDatadogV3_WithoutTags",
			"InstanceID":    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":  "bunny-1",
			"DatadogApiKey": "DATADOG_APIKEY",
			"DatadogRegion": "us1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.#", "1"),
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.0.region", params["DatadogRegion"]),
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.0.tags", ""),
				),
			},
			{
				ResourceName:      prometheusDatadogResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, prometheusDatadogResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusDatadogV3_Update: Test updating Datadog v3 prometheus integration.
func TestAccIntegrationMetricPrometheusDatadogV3_Update(t *testing.T) {
	var (
		fileNames                     = []string{"instance", "integrations/metrics/integration_metric_prometheus_datadog_v3"}
		instanceResourceName          = "cloudamqp_instance.instance"
		prometheusDatadogResourceName = "cloudamqp_integration_metric_prometheus.datadog_v3"

		paramsCreate = map[string]string{
			"InstanceName":  "TestAccIntegrationMetricPrometheusDatadogV3_Update",
			"InstanceID":    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":  "bunny-1",
			"DatadogApiKey": "DATADOG_APIKEY",
			"DatadogRegion": "us1",
			"DatadogTags":   "env=test,region=us1",
		}

		paramsUpdate = map[string]string{
			"InstanceName":  "TestAccIntegrationMetricPrometheusDatadogV3_Update",
			"InstanceID":    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":  "bunny-1",
			"DatadogApiKey": os.Getenv("DATADOG_APIKEY"),
			"DatadogRegion": "us1",
			"DatadogTags":   "env=prod,region=us1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsCreate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsCreate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.0.region", paramsCreate["DatadogRegion"]),
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.0.tags", paramsCreate["DatadogTags"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsUpdate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.0.region", paramsUpdate["DatadogRegion"]),
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.0.tags", paramsUpdate["DatadogTags"]),
				),
			},
		},
	})
}

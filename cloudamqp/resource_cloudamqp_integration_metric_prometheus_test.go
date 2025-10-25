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
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "NEWRELIC_APIKEY"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("NEWRELIC_APIKEY")
	}

	var (
		fileNames                      = []string{"instance", "integrations/metrics/integration_metric_prometheus_newrelic_v3"}
		instanceResourceName           = "cloudamqp_instance.instance"
		prometheusNewRelicResourceName = "cloudamqp_integration_metric_prometheus.newrelic_v3"

		params = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusNewRelicV3_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"NewRelicApiKey": testApiKey,
			"NewRelicTags":   "key=value,key2=value2",
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
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "NEWRELIC_APIKEY"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("NEWRELIC_APIKEY")
	}

	var (
		fileNames                      = []string{"instance", "integrations/metrics/integration_metric_prometheus_newrelic_v3_notags"}
		instanceResourceName           = "cloudamqp_instance.instance"
		prometheusNewRelicResourceName = "cloudamqp_integration_metric_prometheus.newrelic_v3_notags"

		params = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusNewRelicV3_WithoutTags",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"NewRelicApiKey": testApiKey,
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
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "NEWRELIC_APIKEY"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("NEWRELIC_APIKEY")
	}

	var (
		fileNames                      = []string{"instance", "integrations/metrics/integration_metric_prometheus_newrelic_v3"}
		instanceResourceName           = "cloudamqp_instance.instance"
		prometheusNewRelicResourceName = "cloudamqp_integration_metric_prometheus.newrelic_v3"

		paramsCreate = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusNewRelicV3_Update",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"NewRelicApiKey": testApiKey,
			"NewRelicTags":   "key=value,key2=value2",
		}

		paramsUpdate = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusNewRelicV3_Update",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"NewRelicApiKey": testApiKey,
			"NewRelicTags":   "key=value2,key2=value3",
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
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "DATADOG_APIKEY"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("DATADOG_APIKEY")
	}

	var (
		fileNames                     = []string{"instance", "integrations/metrics/integration_metric_prometheus_datadog_v3"}
		instanceResourceName          = "cloudamqp_instance.instance"
		prometheusDatadogResourceName = "cloudamqp_integration_metric_prometheus.datadog_v3"

		params = map[string]string{
			"InstanceName":  "TestAccIntegrationMetricPrometheusDatadogV3_Basic",
			"InstanceID":    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":  "bunny-1",
			"DatadogApiKey": testApiKey,
			"DatadogRegion": "us1",
			"DatadogTags":   "key=value,key2=value2",
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
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "DATADOG_APIKEY"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("DATADOG_APIKEY")
	}

	var (
		fileNames                     = []string{"instance", "integrations/metrics/integration_metric_prometheus_datadog_v3_notags"}
		instanceResourceName          = "cloudamqp_instance.instance"
		prometheusDatadogResourceName = "cloudamqp_integration_metric_prometheus.datadog_v3_notags"

		params = map[string]string{
			"InstanceName":  "TestAccIntegrationMetricPrometheusDatadogV3_WithoutTags",
			"InstanceID":    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":  "bunny-1",
			"DatadogApiKey": testApiKey,
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
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "DATADOG_APIKEY"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("DATADOG_APIKEY")
	}

	var (
		fileNames                     = []string{"instance", "integrations/metrics/integration_metric_prometheus_datadog_v3"}
		instanceResourceName          = "cloudamqp_instance.instance"
		prometheusDatadogResourceName = "cloudamqp_integration_metric_prometheus.datadog_v3"

		paramsCreate = map[string]string{
			"InstanceName":                          "TestAccIntegrationMetricPrometheusDatadogV3_Update",
			"InstanceID":                            fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":                          "bunny-1",
			"DatadogApiKey":                         testApiKey,
			"DatadogRegion":                         "us1",
			"DatadogTags":                           "key=value,key2=value2",
			"DatadogRabbitmqDashboardMetricsFormat": "false",
		}

		paramsUpdate = map[string]string{
			"InstanceName":                          "TestAccIntegrationMetricPrometheusDatadogV3_Update",
			"InstanceID":                            fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":                          "bunny-1",
			"DatadogApiKey":                         testApiKey,
			"DatadogRegion":                         "us1",
			"DatadogTags":                           "key=value2,key2=value3",
			"DatadogRabbitmqDashboardMetricsFormat": "true",
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
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.0.rabbitmq_dashboard_metrics_format", "false"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsUpdate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.0.region", paramsUpdate["DatadogRegion"]),
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.0.tags", paramsUpdate["DatadogTags"]),
					resource.TestCheckResourceAttr(prometheusDatadogResourceName, "datadog_v3.0.rabbitmq_dashboard_metrics_format", "true"),
				),
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusAzureMonitor_Basic: Add Azure Monitor prometheus metric integration and import.
func TestAccIntegrationMetricPrometheusAzureMonitor_Basic(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "AZM_INSTRUMENTATION_KEY"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("AZM_INSTRUMENTATION_KEY")
	}

	var (
		fileNames                          = []string{"instance", "integrations/metrics/integration_metric_prometheus_azure_monitor"}
		instanceResourceName               = "cloudamqp_instance.instance"
		prometheusAzureMonitorResourceName = "cloudamqp_integration_metric_prometheus.azure_monitor"

		params = map[string]string{
			"InstanceName":                 "TestAccIntegrationMetricPrometheusAzureMonitor_Basic",
			"InstanceID":                   fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":                 "bunny-1",
			"AzureMonitorConnectionString": fmt.Sprintf("InstrumentationKey=%s;IngestionEndpoint=https://swedencentral-0.in.applicationinsights.azure.com/;LiveEndpoint=https://swedencentral.livediagnostics.monitor.azure.com/;ApplicationId=3c2ad7f7-65d0-4e39-ae82-8d2fd7b6f69f", testApiKey),
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusAzureMonitorResourceName, "azure_monitor.#", "1"),
				),
			},
			{
				ResourceName:      prometheusAzureMonitorResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, prometheusAzureMonitorResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusAzureMonitor_Update: Test updating Azure Monitor prometheus integration connection string.
func TestAccIntegrationMetricPrometheusAzureMonitor_Update(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "AZM_INSTRUMENTATION_KEY"
	testApiKey2 := "AZM_INSTRUMENTATION_KEY_2"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("AZM_INSTRUMENTATION_KEY")
		testApiKey2 = os.Getenv("AZM_INSTRUMENTATION_KEY_2")
	}

	var (
		fileNames                          = []string{"instance", "integrations/metrics/integration_metric_prometheus_azure_monitor"}
		instanceResourceName               = "cloudamqp_instance.instance"
		prometheusAzureMonitorResourceName = "cloudamqp_integration_metric_prometheus.azure_monitor"

		paramsCreate = map[string]string{
			"InstanceName":                 "TestAccIntegrationMetricPrometheusAzureMonitor_Update",
			"InstanceID":                   fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":                 "bunny-1",
			"AzureMonitorConnectionString": fmt.Sprintf("InstrumentationKey=%s;IngestionEndpoint=https://swedencentral-0.in.applicationinsights.azure.com/;LiveEndpoint=https://swedencentral.livediagnostics.monitor.azure.com/;ApplicationId=3c2ad7f7-65d0-4e39-ae82-8d2fd7b6f69f", testApiKey),
		}

		paramsUpdate = map[string]string{
			"InstanceName":                 "TestAccIntegrationMetricPrometheusAzureMonitor_Update",
			"InstanceID":                   fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":                 "bunny-1",
			"AzureMonitorConnectionString": fmt.Sprintf("InstrumentationKey=%s;IngestionEndpoint=https://swedencentral-0.in.applicationinsights.azure.com/;LiveEndpoint=https://swedencentral.livediagnostics.monitor.azure.com/;ApplicationId=3c2ad7f7-65d0-4e39-ae82-8d2fd7b6f69f", testApiKey2),
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsCreate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsCreate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusAzureMonitorResourceName, "azure_monitor.#", "1"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsUpdate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusAzureMonitorResourceName, "azure_monitor.#", "1"),
				),
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusSplunkV2_Basic: Add Splunk v2 prometheus metric integration and import.
func TestAccIntegrationMetricPrometheusSplunkV2_Basic(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "SPLUNK_TOKEN"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("SPLUNK_TOKEN")
	}

	var (
		fileNames                    = []string{"instance", "integrations/metrics/integration_metric_prometheus_splunk_v2"}
		instanceResourceName         = "cloudamqp_instance.instance"
		prometheusSplunkResourceName = "cloudamqp_integration_metric_prometheus.splunk_v2"

		params = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusSplunkV2_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"SplunkToken":    testApiKey,
			"SplunkEndpoint": "https://prd-p-abcde.splunkcloud.com:8088/services/collector",
			"SplunkTags":     "key=value,key2=value2",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.#", "1"),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.0.endpoint", params["SplunkEndpoint"]),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.0.tags", params["SplunkTags"]),
				),
			},
			{
				ResourceName:      prometheusSplunkResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, prometheusSplunkResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusSplunkV2_WithoutTags: Test Splunk v2 prometheus integration without optional tags.
func TestAccIntegrationMetricPrometheusSplunkV2_WithoutTags(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "SPLUNK_TOKEN"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("SPLUNK_TOKEN")
	}

	var (
		fileNames                    = []string{"instance", "integrations/metrics/integration_metric_prometheus_splunk_v2_notags"}
		instanceResourceName         = "cloudamqp_instance.instance"
		prometheusSplunkResourceName = "cloudamqp_integration_metric_prometheus.splunk_v2_notags"

		params = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusSplunkV2_WithoutTags",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"SplunkToken":    testApiKey,
			"SplunkEndpoint": "https://prd-p-abcde.splunkcloud.com:8088/services/collector",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.#", "1"),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.0.endpoint", params["SplunkEndpoint"]),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.0.tags", ""),
				),
			},
			{
				ResourceName:      prometheusSplunkResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, prometheusSplunkResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusSplunkV2_Update: Test updating Splunk v2 prometheus integration.
func TestAccIntegrationMetricPrometheusSplunkV2_Update(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "SPLUNK_TOKEN"
	testApiKey2 := "SPLUNK_TOKEN_2"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("SPLUNK_TOKEN")
		testApiKey2 = os.Getenv("SPLUNK_TOKEN_2")
	}

	var (
		fileNames                    = []string{"instance", "integrations/metrics/integration_metric_prometheus_splunk_v2"}
		instanceResourceName         = "cloudamqp_instance.instance"
		prometheusSplunkResourceName = "cloudamqp_integration_metric_prometheus.splunk_v2"

		paramsCreate = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusSplunkV2_Update",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"SplunkToken":    testApiKey,
			"SplunkEndpoint": "https://prd-p-abcde.splunkcloud.com:8088/services/collector",
			"SplunkTags":     "key=value,key2=value2",
		}

		paramsUpdate = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusSplunkV2_Update",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"SplunkToken":    testApiKey2,
			"SplunkEndpoint": "https://prd-p-fghij.splunkcloud.com:8088/services/collector",
			"SplunkTags":     "key=value2,key2=value3",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsCreate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsCreate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.0.endpoint", paramsCreate["SplunkEndpoint"]),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.0.tags", paramsCreate["SplunkTags"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsUpdate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.0.endpoint", paramsUpdate["SplunkEndpoint"]),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.0.tags", paramsUpdate["SplunkTags"]),
				),
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusDynatrace_Basic: Add Dynatrace prometheus metric integration and import.
func TestAccIntegrationMetricPrometheusDynatrace_Basic(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "DYNATRACE_TOKEN"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("DYNATRACE_TOKEN")
	}

	var (
		fileNames                       = []string{"instance", "integrations/metrics/integration_metric_prometheus_dynatrace"}
		instanceResourceName            = "cloudamqp_instance.instance"
		prometheusDynatraceResourceName = "cloudamqp_integration_metric_prometheus.dynatrace"

		params = map[string]string{
			"InstanceName":           "TestAccIntegrationMetricPrometheusDynatrace_Basic",
			"InstanceID":             fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":           "bunny-1",
			"DynatraceEnvironmentID": "abc12345",
			"DynatraceAccessToken":   testApiKey,
			"DynatraceTags":          "env=prod,service=rabbitmq",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusDynatraceResourceName, "dynatrace.#", "1"),
					resource.TestCheckResourceAttr(prometheusDynatraceResourceName, "dynatrace.0.environment_id", params["DynatraceEnvironmentID"]),
					resource.TestCheckResourceAttr(prometheusDynatraceResourceName, "dynatrace.0.tags", params["DynatraceTags"]),
				),
			},
			{
				ResourceName:      prometheusDynatraceResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, prometheusDynatraceResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusDynatrace_WithoutTags: Test Dynatrace prometheus integration without optional tags.
func TestAccIntegrationMetricPrometheusDynatrace_WithoutTags(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "DYNATRACE_TOKEN"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("DYNATRACE_TOKEN")
	}

	var (
		fileNames                       = []string{"instance", "integrations/metrics/integration_metric_prometheus_dynatrace_notags"}
		instanceResourceName            = "cloudamqp_instance.instance"
		prometheusDynatraceResourceName = "cloudamqp_integration_metric_prometheus.dynatrace_notags"

		params = map[string]string{
			"InstanceName":           "TestAccIntegrationMetricPrometheusDynatrace_WithoutTags",
			"InstanceID":             fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":           "bunny-1",
			"DynatraceEnvironmentID": "abc12345",
			"DynatraceAccessToken":   testApiKey,
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusDynatraceResourceName, "dynatrace.#", "1"),
					resource.TestCheckResourceAttr(prometheusDynatraceResourceName, "dynatrace.0.environment_id", params["DynatraceEnvironmentID"]),
					resource.TestCheckResourceAttr(prometheusDynatraceResourceName, "dynatrace.0.tags", ""),
				),
			},
			{
				ResourceName:      prometheusDynatraceResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, prometheusDynatraceResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusDynatrace_Update: Test updating Dynatrace prometheus integration.
func TestAccIntegrationMetricPrometheusDynatrace_Update(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "DYNATRACE_TOKEN"
	testApiKey_2 := "DYNATRACE_TOKEN_2"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("DYNATRACE_TOKEN")
		testApiKey_2 = os.Getenv("DYNATRACE_TOKEN_2")
	}

	var (
		fileNames                       = []string{"instance", "integrations/metrics/integration_metric_prometheus_dynatrace"}
		instanceResourceName            = "cloudamqp_instance.instance"
		prometheusDynatraceResourceName = "cloudamqp_integration_metric_prometheus.dynatrace"

		params = map[string]string{
			"InstanceName":           "TestAccIntegrationMetricPrometheusDynatrace_Update",
			"InstanceID":             fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":           "bunny-1",
			"DynatraceEnvironmentID": "abc12345",
			"DynatraceAccessToken":   testApiKey,
			"DynatraceTags":          "env=prod,service=rabbitmq",
		}

		paramsUpdate = map[string]string{
			"InstanceName":           "TestAccIntegrationMetricPrometheusDynatrace_Update",
			"InstanceID":             fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":           "bunny-1",
			"DynatraceEnvironmentID": "xyz67890",
			"DynatraceAccessToken":   testApiKey_2,
			"DynatraceTags":          "env=staging,service=messaging,team=platform",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusDynatraceResourceName, "dynatrace.0.environment_id", params["DynatraceEnvironmentID"]),
					resource.TestCheckResourceAttr(prometheusDynatraceResourceName, "dynatrace.0.tags", params["DynatraceTags"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsUpdate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusDynatraceResourceName, "dynatrace.0.environment_id", paramsUpdate["DynatraceEnvironmentID"]),
					resource.TestCheckResourceAttr(prometheusDynatraceResourceName, "dynatrace.0.tags", paramsUpdate["DynatraceTags"]),
				),
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusCloudwatchV3_Basic: Add CloudWatch v3 prometheus metric integration and import.
func TestAccIntegrationMetricPrometheusCloudwatchV3_Basic(t *testing.T) {
	t.Parallel()

	var (
		fileNames                        = []string{"instance", "integrations/metrics/integration_metric_prometheus_cloudwatch_v3"}
		instanceResourceName             = "cloudamqp_instance.instance"
		prometheusCloudwatchResourceName = "cloudamqp_integration_metric_prometheus.cloudwatch_v3"

		params = map[string]string{
			"InstanceName":            "TestAccIntegrationMetricPrometheusCloudwatchV3_Basic",
			"InstanceID":              fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":            "bunny-1",
			"CloudwatchIAMRole":       "arn:aws:iam::123456789012:role/cloudamqp-role",
			"CloudwatchIAMExternalID": "cloudamqp-external-id-123",
			"CloudwatchRegion":        "us-east-1",
			"CloudwatchTags":          "env=test,service=rabbitmq",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.#", "1"),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.iam_role", params["CloudwatchIAMRole"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.iam_external_id", params["CloudwatchIAMExternalID"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.region", params["CloudwatchRegion"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.tags", params["CloudwatchTags"]),
				),
			},
			{
				ResourceName:      prometheusCloudwatchResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, prometheusCloudwatchResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusCloudwatchV3_WithoutTags: Test CloudWatch v3 prometheus integration without optional tags.
func TestAccIntegrationMetricPrometheusCloudwatchV3_WithoutTags(t *testing.T) {
	t.Parallel()

	var (
		fileNames                        = []string{"instance", "integrations/metrics/integration_metric_prometheus_cloudwatch_v3_notags"}
		instanceResourceName             = "cloudamqp_instance.instance"
		prometheusCloudwatchResourceName = "cloudamqp_integration_metric_prometheus.cloudwatch_v3_notags"

		params = map[string]string{
			"InstanceName":            "TestAccIntegrationMetricPrometheusCloudwatchV3_WithoutTags",
			"InstanceID":              fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":            "bunny-1",
			"CloudwatchIAMRole":       "arn:aws:iam::123456789012:role/cloudamqp-role",
			"CloudwatchIAMExternalID": "cloudamqp-external-id-123",
			"CloudwatchRegion":        "us-east-1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.#", "1"),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.iam_role", params["CloudwatchIAMRole"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.iam_external_id", params["CloudwatchIAMExternalID"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.region", params["CloudwatchRegion"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.tags", ""),
				),
			},
			{
				ResourceName:      prometheusCloudwatchResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, prometheusCloudwatchResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusCloudwatchV3_Update: Test updating CloudWatch v3 prometheus integration.
func TestAccIntegrationMetricPrometheusCloudwatchV3_Update(t *testing.T) {
	t.Parallel()

	var (
		fileNames                        = []string{"instance", "integrations/metrics/integration_metric_prometheus_cloudwatch_v3"}
		instanceResourceName             = "cloudamqp_instance.instance"
		prometheusCloudwatchResourceName = "cloudamqp_integration_metric_prometheus.cloudwatch_v3"

		paramsCreate = map[string]string{
			"InstanceName":            "TestAccIntegrationMetricPrometheusCloudwatchV3_Update",
			"InstanceID":              fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":            "bunny-1",
			"CloudwatchIAMRole":       "arn:aws:iam::123456789012:role/cloudamqp-role",
			"CloudwatchIAMExternalID": "cloudamqp-external-id-123",
			"CloudwatchRegion":        "us-east-1",
			"CloudwatchTags":          "env=test,service=rabbitmq",
		}

		paramsUpdate = map[string]string{
			"InstanceName":            "TestAccIntegrationMetricPrometheusCloudwatchV3_Update",
			"InstanceID":              fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":            "bunny-1",
			"CloudwatchIAMRole":       "arn:aws:iam::987654321098:role/cloudamqp-role-updated",
			"CloudwatchIAMExternalID": "cloudamqp-external-id-456",
			"CloudwatchRegion":        "us-west-2",
			"CloudwatchTags":          "env=prod,service=messaging,team=platform",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsCreate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsCreate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.iam_role", paramsCreate["CloudwatchIAMRole"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.iam_external_id", paramsCreate["CloudwatchIAMExternalID"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.region", paramsCreate["CloudwatchRegion"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.tags", paramsCreate["CloudwatchTags"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsUpdate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.iam_role", paramsUpdate["CloudwatchIAMRole"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.iam_external_id", paramsUpdate["CloudwatchIAMExternalID"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.region", paramsUpdate["CloudwatchRegion"]),
					resource.TestCheckResourceAttr(prometheusCloudwatchResourceName, "cloudwatch_v3.0.tags", paramsUpdate["CloudwatchTags"]),
				),
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusMetricsFilter_Default: Test prometheus integration with default metrics (no metrics_filter specified).
func TestAccIntegrationMetricPrometheusMetricsFilter_Default(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "SPLUNK_TOKEN"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("SPLUNK_TOKEN")
	}

	var (
		fileNames                    = []string{"instance", "integrations/metrics/integration_metric_prometheus_splunk_v2"}
		instanceResourceName         = "cloudamqp_instance.instance"
		prometheusSplunkResourceName = "cloudamqp_integration_metric_prometheus.splunk_v2"

		params = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusMetricsFilter_Default",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"SplunkToken":    testApiKey,
			"SplunkEndpoint": "https://prd-p-abcde.splunkcloud.com:8088/services/collector",
			"SplunkTags":     "key=value,key2=value2",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.#", "1"),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.0.endpoint", params["SplunkEndpoint"]),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "metrics_filter.#", "37"),
				),
			},
			{
				ResourceName:      prometheusSplunkResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, prometheusSplunkResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetricPrometheusMetricsFilter_Custom: Test prometheus integration with custom metrics_filter.
func TestAccIntegrationMetricPrometheusMetricsFilter_Custom(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testApiKey := "SPLUNK_TOKEN"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testApiKey = os.Getenv("SPLUNK_TOKEN")
	}

	var (
		fileNames                    = []string{"instance", "integrations/metrics/integration_metric_prometheus_splunk_v2_metrics_filter"}
		instanceResourceName         = "cloudamqp_instance.instance"
		prometheusSplunkResourceName = "cloudamqp_integration_metric_prometheus.splunk_v2_metrics_filter"

		params = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheusMetricsFilter_Custom",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"SplunkToken":    testApiKey,
			"SplunkEndpoint": "https://prd-p-abcde.splunkcloud.com:8088/services/collector",
			"MetricsFilter":  `["rabbitmq_connections", "rabbitmq_detailed_queue_messages", "system_cpu_utilization_ratio"]`,
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.#", "1"),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "splunk_v2.0.endpoint", params["SplunkEndpoint"]),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "metrics_filter.#", "3"),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "metrics_filter.0", "rabbitmq_connections"),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "metrics_filter.1", "rabbitmq_detailed_queue_messages"),
					resource.TestCheckResourceAttr(prometheusSplunkResourceName, "metrics_filter.2", "system_cpu_utilization_ratio"),
				),
			},
			{
				ResourceName:      prometheusSplunkResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, prometheusSplunkResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

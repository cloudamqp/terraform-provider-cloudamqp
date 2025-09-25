package cloudamqp

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccIntegrationMetricPrometheus_Basic: Add prometheus metric integrations and import.
func TestAccIntegrationMetricPrometheus_Basic(t *testing.T) {
	var (
		fileNames                      = []string{"instance", "integrations/metrics/integration_metric_prometheus"}
		instanceResourceName           = "cloudamqp_instance.instance"
		prometheusNewRelicResourceName = "cloudamqp_integration_metric_prometheus.newrelic_v3"

		params = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheus_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"NewRelicApiKey": os.Getenv("NEWRELIC_V3_APIKEY"),
			"NewRelicTags":   "env=test,region=us1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             configuration.GetTemplatedConfig(t, fileNames, params),
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

// TestAccIntegrationMetricPrometheus_WithoutTags: Test prometheus integration without optional tags.
func TestAccIntegrationMetricPrometheus_WithoutTags(t *testing.T) {
	var (
		fileNames                      = []string{"instance", "integrations/metrics/integration_metric_prometheus_notags"}
		instanceResourceName           = "cloudamqp_instance.instance"
		prometheusNewRelicResourceName = "cloudamqp_integration_metric_prometheus.newrelic_v3_notags"

		params = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheus_WithoutTags",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"NewRelicApiKey": os.Getenv("NEWRELIC_V3_APIKEY"),
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             configuration.GetTemplatedConfig(t, fileNames, params),
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

// TestAccIntegrationMetricPrometheus_Update: Test updating prometheus integration.
func TestAccIntegrationMetricPrometheus_Update(t *testing.T) {
	var (
		fileNames                      = []string{"instance", "integrations/metrics/integration_metric_prometheus"}
		instanceResourceName           = "cloudamqp_instance.instance"
		prometheusNewRelicResourceName = "cloudamqp_integration_metric_prometheus.newrelic_v3"

		paramsCreate = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheus_Update",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"NewRelicApiKey": os.Getenv("NEWRELIC_V3_APIKEY"),
			"NewRelicTags":   "env=test,region=us1",
		}

		paramsUpdate = map[string]string{
			"InstanceName":   "TestAccIntegrationMetricPrometheus_Update",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"NewRelicApiKey": os.Getenv("NEWRELIC_V3_APIKEY"),
			"NewRelicTags":   "env=prod,region=eu1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             configuration.GetTemplatedConfig(t, fileNames, paramsCreate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsCreate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusNewRelicResourceName, "newrelic_v3.0.tags", paramsCreate["NewRelicTags"]),
				),
			},
			{
				ExpectNonEmptyPlan: true,
				Config:             configuration.GetTemplatedConfig(t, fileNames, paramsUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsUpdate["InstanceName"]),
					resource.TestCheckResourceAttr(prometheusNewRelicResourceName, "newrelic_v3.0.tags", paramsUpdate["NewRelicTags"]),
				),
			},
		},
	})
}

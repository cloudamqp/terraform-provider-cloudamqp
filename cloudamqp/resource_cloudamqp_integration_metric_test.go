package cloudamqp

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccIntegrationMetric_CloudWatch: Test CloudWatch v2 integration.
func TestAccIntegrationMetric_CloudWatch(t *testing.T) {
	var (
		fileNames              = []string{"instance", "integrations/metrics/cloudwatch_v2"}
		instanceResourceName   = "cloudamqp_instance.instance"
		cloudwatchResourceName = "cloudamqp_integration_metric.cloudwatch_v2"

		params = map[string]string{
			"InstanceName":              "TestAccIntegrationMetric_CloudWatch",
			"InstanceID":                fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":              "bunny-1",
			"CloudwatchAccessKeyId":     os.Getenv("CLOUDWATCH_ACCESS_KEY_ID"),
			"CloudwatchSecretAccessKey": os.Getenv("CLOUDWATCH_SECRET_ACCESS_KEY"),
			"CloudwatchRegion":          "us-east-1",
			"CloudwatchTags":            "env=test,region=us-east-1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:             configuration.GetTemplatedConfig(t, fileNames, params),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "name", "cloudwatch_v2"),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "region", params["CloudwatchRegion"]),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "tags", params["CloudwatchTags"]),
				),
			},
			{
				ResourceName:      cloudwatchResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, cloudwatchResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetric_DataDog: Test DataDog v2 integration.
func TestAccIntegrationMetric_DataDog(t *testing.T) {
	var (
		fileNames            = []string{"instance", "integrations/metrics/datadog_v2"}
		instanceResourceName = "cloudamqp_instance.instance"
		dataDogResourceName  = "cloudamqp_integration_metric.datadog_v2"

		params = map[string]string{
			"InstanceName":  "TestAccIntegrationMetric_DataDog",
			"InstanceID":    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":  "bunny-1",
			"DataDogRegion": "us1",
			"DataDogApiKey": os.Getenv("DATADOG_APIKEY"),
			"DataDogTags":   "env=test,region=us1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:             configuration.GetTemplatedConfig(t, fileNames, params),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(dataDogResourceName, "name", "datadog_v2"),
					resource.TestCheckResourceAttr(dataDogResourceName, "region", params["DataDogRegion"]),
					resource.TestCheckResourceAttr(dataDogResourceName, "tags", params["DataDogTags"]),
				),
			},
			{
				ResourceName:      dataDogResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, dataDogResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetric_Librato: Test Librato integration.
func TestAccIntegrationMetric_Librato(t *testing.T) {
	var (
		fileNames            = []string{"instance", "integrations/metrics/librato"}
		instanceResourceName = "cloudamqp_instance.instance"
		libratoResourceName  = "cloudamqp_integration_metric.librato"

		params = map[string]string{
			"InstanceName":  "TestAccIntegrationMetric_Librato",
			"InstanceID":    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":  "bunny-1",
			"LibratoEmail":  "test@example.com",
			"LibratoApiKey": os.Getenv("LIBRATO_APIKEY"),
			"LibratoTags":   "env=test",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:             configuration.GetTemplatedConfig(t, fileNames, params),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(libratoResourceName, "name", "librato"),
					resource.TestCheckResourceAttr(libratoResourceName, "email", params["LibratoEmail"]),
					resource.TestCheckResourceAttr(libratoResourceName, "tags", params["LibratoTags"]),
				),
			},
			{
				ResourceName:      libratoResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, libratoResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccIntegrationMetric_NewRelic: Test New Relic v2 integration.
func TestAccIntegrationMetric_NewRelic(t *testing.T) {
	var (
		fileNames            = []string{"instance", "integrations/metrics/newrelic_v2"}
		instanceResourceName = "cloudamqp_instance.instance"
		newrelicResourceName = "cloudamqp_integration_metric.newrelic_v2"

		params = map[string]string{
			"InstanceName":   "TestAccIntegrationMetric_NewRelic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"NewRelicApiKey": os.Getenv("NEWRELIC_APIKEY"),
			"NewRelicRegion": "us",
			"NewRelicTags":   "env=test,region=us1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:             configuration.GetTemplatedConfig(t, fileNames, params),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(newrelicResourceName, "name", "newrelic_v2"),
					resource.TestCheckResourceAttr(newrelicResourceName, "region", params["NewRelicRegion"]),
					resource.TestCheckResourceAttr(newrelicResourceName, "tags", params["NewRelicTags"]),
				),
			},
			{
				ResourceName:      newrelicResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, newrelicResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

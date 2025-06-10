package cloudamqp

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccIntegrationMetric_Basic: Add metric integrations and import.
func TestAccIntegrationMetric_Basic(t *testing.T) {
	var (
		fileNames              = []string{"instance", "integration_metric"}
		instanceResourceName   = "cloudamqp_instance.instance"
		cloudwatchResourceName = "cloudamqp_integration_metric.cloudwatch_v2"
		dataDogResourceName    = "cloudamqp_integration_metric.datadog_v2"
		libratoResourceName    = "cloudamqp_integration_metric.librato"
		newrelicResourceName   = "cloudamqp_integration_metric.newrelic_v2"

		params = map[string]string{
			"InstanceName":              "TestAccIntegrationMetric_Basic",
			"InstanceID":                fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":              "bunny-1",
			"CloudwatchAccessKeyId":     os.Getenv("CLOUDWATCH_ACCESS_KEY_ID"),
			"CloudwatchSecretAccessKey": os.Getenv("CLOUDWATCH_SECRET_ACCESS_KEY"),
			"CloudwatchRegion":          "us-east-1",
			"CloudwatchTags":            "env=test,region=us-east-1",
			"DataDogRegion":             "us1",
			"DataDogApiKey":             os.Getenv("DATADOG_APIKEY"),
			"DataDogTags":               "env=test,region=us1",
			"LibratoEmail":              "test@example.com",
			"LibratoApiKey":             os.Getenv("LIBRATO_APIKEY"),
			"LibratoTags":               "env=test",
			"NewRelicApiKey":            os.Getenv("NEWRELIC_APIKEY"),
			"NewRelicRegion":            "us",
			"NewRelicTags":              "env=test,region=us1",
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
					resource.TestCheckResourceAttr(cloudwatchResourceName, "name", "cloudwatch_v2"),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "access_key_id", "CLOUDWATCH_ACCESS_KEY_ID"),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "region", params["CloudwatchRegion"]),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "tags", params["CloudwatchTags"]),
					resource.TestCheckResourceAttr(dataDogResourceName, "name", "datadog_v2"),
					resource.TestCheckResourceAttr(dataDogResourceName, "region", params["DataDogRegion"]),
					resource.TestCheckResourceAttr(dataDogResourceName, "tags", params["DataDogTags"]),
					resource.TestCheckResourceAttr(libratoResourceName, "name", "librato"),
					resource.TestCheckResourceAttr(libratoResourceName, "email", params["LibratoEmail"]),
					resource.TestCheckResourceAttr(libratoResourceName, "tags", params["LibratoTags"]),
					resource.TestCheckResourceAttr(newrelicResourceName, "name", "newrelic_v2"),
					resource.TestCheckResourceAttr(newrelicResourceName, "region", params["NewRelicRegion"]),
					resource.TestCheckResourceAttr(newrelicResourceName, "tags", params["NewRelicTags"]),
				),
			},
			{
				ResourceName:      cloudwatchResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, cloudwatchResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      dataDogResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, dataDogResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      libratoResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, libratoResourceName),
				ImportState:       true,
				ImportStateVerify: true,
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

package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
			"CloudwatchAccessKeyId":     "AKIAI44QH8DHBEXAMPLE",                     // Example key id
			"CloudwatchSecretAccessKey": "je7MtGbClwBFd2Zp9Utkdh3yCo8nvbEXAMPLEKEY", // Example secret key
			"CloudwatchRegion":          "us-east-1",
			"CloudwatchTags":            "env=test,region=us-east-1",
			"DataDogRegion":             "us1",
			"DataDogApiKey":             "1af4f17471e98bcee88b6d9d6ba1626f", // Note: require real (temporary) key when recording.
			"DataDogTags":               "env=test,region=us1",
			"LibratoEmail":              "test@example.com",
			"LibratoApiKey":             "7b857ea2-b9d3-4268-955f-7e4b4abf877c", // Randomized token
			"LibratoTags":               "env=test",
			"NewRelicApiKey":            "9985ba19-f566-48fa-b90a-628474004067", // Randomized token
			"NewRelicRegion":            "us",
			"NewRelicTags":              "env=test,region=us1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "name", "cloudwatch_v2"),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "access_key_id", params["CloudwatchAccessKeyId"]),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "region", params["CloudwatchRegion"]),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "tags", params["CloudwatchTags"]),
					resource.TestCheckResourceAttr(dataDogResourceName, "name", "datadog_v2"),
					resource.TestCheckResourceAttr(dataDogResourceName, "region", params["DataDogRegion"]),
					resource.TestCheckResourceAttr(dataDogResourceName, "tags", params["DataDogTags"]),
					resource.TestCheckResourceAttr(libratoResourceName, "name", "librato"),
					resource.TestCheckResourceAttr(libratoResourceName, "email", params["LibratoEmail"]),
					resource.TestCheckResourceAttr(libratoResourceName, "api_key", params["LibratoApiKey"]),
					resource.TestCheckResourceAttr(libratoResourceName, "tags", params["LibratoTags"]),
					resource.TestCheckResourceAttr(newrelicResourceName, "name", "newrelic_v2"),
					resource.TestCheckResourceAttr(newrelicResourceName, "api_key", params["NewRelicApiKey"]),
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

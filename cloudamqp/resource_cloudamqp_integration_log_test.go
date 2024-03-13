package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccIntegrationLog_Basic: Add log integrations and import.
func TestAccIntegrationLog_Basic(t *testing.T) {
	var (
		fileNames              = []string{"instance", "integration_log"}
		instanceResourceName   = "cloudamqp_instance.instance"
		azmResourceName        = "cloudamqp_integration_log.azure_monitor"
		cloudwatchResourceName = "cloudamqp_integration_log.cloudwatch"
		coralogixResourceName  = "cloudamqp_integration_log.coralogix"
		dataDogResourceName    = "cloudamqp_integration_log.datadog"
		logentriesResourceName = "cloudamqp_integration_log.logentries"
		logglyResourceName     = "cloudamqp_integration_log.loggly"
		papertrailResourceName = "cloudamqp_integration_log.papertrail"
		scalyrResourceName     = "cloudamqp_integration_log.scalyr"
		splunkResourceName     = "cloudamqp_integration_log.splunk"

		params = map[string]string{
			"InstanceName":              "TestAccIntegrationLog_Basic",
			"InstanceID":                fmt.Sprintf("%s.id", instanceResourceName),
			"InstanceHost":              fmt.Sprintf("%s.host", instanceResourceName),
			"AzmTentantId":              "71e89a32-14f3-4458-b136-7395bb6d1969",     // Radnomized token
			"AzmApplicationId":          "3e303e72-4024-494c-b5f6-f5ffbe8139de",     // Radnomized token
			"AzmApplicationSecret":      "DA10F~FSqsdjnW3nHFWwXdeW1zdvqIQhdSTfVdes", // Radnomized token
			"AzmDcrId":                  "dcr-7cae904d070344d7ace2b8b33b743c84",
			"AzmDceUri":                 "https://cloudamqp-log-integration.australiasoutheast-1.ingest.monitor.azure.com",
			"AzmTable":                  "cloudamqp_CL",
			"CloudwatchAccessKeyId":     "AKIAI44QH8DHBEXAMPLE",                     // Example key id
			"CloudwatchSecretAccessKey": "je7MtGbClwBFd2Zp9Utkdh3yCo8nvbEXAMPLEKEY", // Example secret key
			"CloudwatchRegion":          "us-east-1",
			"CoralogixSendDataKey":      "ca755454-823b-46e9-9f7e-996baa35249b", // Radnomized token
			"CoralogixEndpoint":         "syslog.cx498.coralogix.com:6514",
			"CoralogixApplication":      "playground",
			"DataDogRegion":             "us1",
			"DataDogApiKey":             "4b08474cdead14fb57a1099ba2b32ee6", // Note: require real (temporary) key when recording.
			"DataDogTags":               "env=test,region=us1",
			"LogEntriesToken":           "10de0c0c-6a65-4070-9501-177d46a3f8f0", // Radnomized token
			"LogglyToken":               "dec02f13-7b54-4874-b339-d80ecb02299b", // Randomized token
			"PapertrailUrl":             "logs.papertrailapp.com:11111",
			"ScalyrToken":               "3dUM/LLdkodksksDKK2lsjkd9kdkd/2djjdJdi8ejsld-", // Randomized token
			"ScalyrHost":                "app.scalyr.com",
			"SplunkToken":               "53f96e41-857d-4fa0-a609-8bb7f2776737", // Randomized token
			"SplunkHostPort":            "logs.splunk.com:11111",
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
					resource.TestCheckResourceAttr(azmResourceName, "name", "azure_monitor"),
					resource.TestCheckResourceAttr(azmResourceName, "table", params["AzmTable"]),
					resource.TestCheckResourceAttr(azmResourceName, "dcr_id", params["AzmDcrId"]),
					resource.TestCheckResourceAttr(azmResourceName, "dce_uri", params["AzmDceUri"]),
					resource.TestCheckResourceAttr(azmResourceName, "tenant_id", params["AzmTentantId"]),
					resource.TestCheckResourceAttr(azmResourceName, "application_id", params["AzmApplicationId"]),
					resource.TestCheckResourceAttr(azmResourceName, "application_secret", params["AzmApplicationSecret"]),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "name", "cloudwatchlog"),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "access_key_id", params["CloudwatchAccessKeyId"]),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "region", params["CloudwatchRegion"]),
					resource.TestCheckResourceAttr(coralogixResourceName, "name", "coralogix"),
					resource.TestCheckResourceAttr(coralogixResourceName, "private_key", params["CoralogixSendDataKey"]),
					resource.TestCheckResourceAttr(coralogixResourceName, "endpoint", params["CoralogixEndpoint"]),
					resource.TestCheckResourceAttr(coralogixResourceName, "application", params["CoralogixApplication"]),
					resource.TestCheckResourceAttr(dataDogResourceName, "name", "datadog"),
					resource.TestCheckResourceAttr(dataDogResourceName, "region", params["DataDogRegion"]),
					resource.TestCheckResourceAttr(logentriesResourceName, "name", "logentries"),
					resource.TestCheckResourceAttr(logentriesResourceName, "token", params["LogEntriesToken"]),
					resource.TestCheckResourceAttr(logglyResourceName, "name", "loggly"),
					resource.TestCheckResourceAttr(logglyResourceName, "token", params["LogglyToken"]),
					resource.TestCheckResourceAttr(papertrailResourceName, "name", "papertrail"),
					resource.TestCheckResourceAttr(papertrailResourceName, "url", params["PapertrailUrl"]),
					resource.TestCheckResourceAttr(scalyrResourceName, "name", "scalyr"),
					resource.TestCheckResourceAttr(scalyrResourceName, "token", params["ScalyrToken"]),
					resource.TestCheckResourceAttr(scalyrResourceName, "host", params["ScalyrHost"]),
					resource.TestCheckResourceAttr(splunkResourceName, "name", "splunk"),
					resource.TestCheckResourceAttr(splunkResourceName, "token", params["SplunkToken"]),
					resource.TestCheckResourceAttr(splunkResourceName, "host_port", params["SplunkHostPort"]),
				),
			},
			{
				ResourceName:      azmResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, azmResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      cloudwatchResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, cloudwatchResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      coralogixResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, coralogixResourceName),
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
				ResourceName:      logentriesResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, logentriesResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      logglyResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, logglyResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      papertrailResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, papertrailResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      splunkResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, splunkResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

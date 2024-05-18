package cloudamqp

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
			"InstancePlan":              "bunny-1",
			"InstanceHost":              fmt.Sprintf("%s.host", instanceResourceName),
			"AzmTentantId":              "71e89a32-14f3-4458-b136-7395bb6d1969", // Randomized token
			"AzmApplicationId":          "3e303e72-4024-494c-b5f6-f5ffbe8139de", // Randomized token
			"AzmApplicationSecret":      os.Getenv("AZM_APPLICATION_SECRET"),
			"AzmDcrId":                  "dcr-7cae904d070344d7ace2b8b33b743c84",
			"AzmDceUri":                 "https://cloudamqp-log-integration.australiasoutheast-1.ingest.monitor.azure.com",
			"AzmTable":                  "cloudamqp_CL",
			"CloudwatchAccessKeyId":     os.Getenv("CLOUDWATCH_ACCESS_KEY_ID"),
			"CloudwatchSecretAccessKey": os.Getenv("CLOUDWATCH_SECRET_ACCESS_KEY"),
			"CloudwatchRegion":          "us-east-1",
			"CoralogixSendDataKey":      os.Getenv("CORALOGIX_SEND_DATA_KEY"),
			"CoralogixEndpoint":         "syslog.cx498.coralogix.com:6514",
			"CoralogixApplication":      "playground",
			"DataDogRegion":             "us1",
			"DataDogApiKey":             os.Getenv("DATADOG_APIKEY"),
			"DataDogTags":               "env=test,region=us1",
			"LogEntriesToken":           os.Getenv("LOGENTIRES_TOKEN"),
			"LogglyToken":               os.Getenv("LOGGLY_TOKEN"),
			"PapertrailUrl":             "logs.papertrailapp.com:11111",
			"ScalyrToken":               os.Getenv("SCALYR_TOKEN"),
			"ScalyrHost":                "app.scalyr.com",
			"SplunkToken":               os.Getenv("SPLUNK_TOKEN"),
			"SplunkHostPort":            "logs.splunk.com:11111",
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
					resource.TestCheckResourceAttr(azmResourceName, "name", "azure_monitor"),
					resource.TestCheckResourceAttr(azmResourceName, "table", params["AzmTable"]),
					resource.TestCheckResourceAttr(azmResourceName, "dcr_id", params["AzmDcrId"]),
					resource.TestCheckResourceAttr(azmResourceName, "dce_uri", params["AzmDceUri"]),
					resource.TestCheckResourceAttr(azmResourceName, "tenant_id", params["AzmTentantId"]),
					resource.TestCheckResourceAttr(azmResourceName, "application_id", params["AzmApplicationId"]),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "name", "cloudwatchlog"),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "access_key_id", "CLOUDWATCH_ACCESS_KEY_ID"),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "region", params["CloudwatchRegion"]),
					resource.TestCheckResourceAttr(coralogixResourceName, "name", "coralogix"),
					resource.TestCheckResourceAttr(coralogixResourceName, "endpoint", params["CoralogixEndpoint"]),
					resource.TestCheckResourceAttr(coralogixResourceName, "application", params["CoralogixApplication"]),
					resource.TestCheckResourceAttr(dataDogResourceName, "name", "datadog"),
					resource.TestCheckResourceAttr(dataDogResourceName, "region", params["DataDogRegion"]),
					resource.TestCheckResourceAttr(logentriesResourceName, "name", "logentries"),
					resource.TestCheckResourceAttr(logglyResourceName, "name", "loggly"),
					resource.TestCheckResourceAttr(papertrailResourceName, "name", "papertrail"),
					resource.TestCheckResourceAttr(papertrailResourceName, "url", params["PapertrailUrl"]),
					resource.TestCheckResourceAttr(scalyrResourceName, "name", "scalyr"),
					resource.TestCheckResourceAttr(scalyrResourceName, "host", params["ScalyrHost"]),
					resource.TestCheckResourceAttr(splunkResourceName, "name", "splunk"),
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

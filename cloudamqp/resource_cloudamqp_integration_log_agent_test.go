package cloudamqp

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIntegrationLogAgent_Cloudwatch_Basic(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testIAMRole := "CLOUDWATCH_IAM_ROLE"
	testIAMExternalID := "CLOUDWATCH_IAM_EXTERNAL_ID"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testIAMRole = os.Getenv("CLOUDWATCH_IAM_ROLE")
		testIAMExternalID = os.Getenv("CLOUDWATCH_IAM_EXTERNAL_ID")
	}

	var (
		instanceResourceName   = "cloudamqp_instance.instance"
		cloudwatchResourceName = "cloudamqp_integration_log_agent.cloudwatch"
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
					  name   = "TestAccIntegrationLogAgent_Cloudwatch_Basic"
					  plan   = "penguin-1"
					  region = "amazon-web-services::eu-central-1"
					  tags   = ["vcr-test"]
					}

					resource "cloudamqp_integration_log_agent" "cloudwatch" {
					  instance_id = cloudamqp_instance.instance.id
					  cloudwatch {
					    iam_role        = "%s"
					    iam_external_id = "%s"
					    region          = "eu-central-1"
					    log_stream      = cloudamqp_instance.instance.cluster_name
					  }
					}
				`, testIAMRole, testIAMExternalID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccIntegrationLogAgent_Cloudwatch_Basic"),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "cloudwatch.iam_role", testIAMRole),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "cloudwatch.iam_external_id", testIAMExternalID),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "cloudwatch.region", "eu-central-1"),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "cloudwatch.log_group", "CloudAMQP"),
					resource.TestCheckResourceAttrPair(cloudwatchResourceName, "cloudwatch.log_stream", instanceResourceName, "cluster_name"),
				),
			},
			{
				ResourceName:      cloudwatchResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, cloudwatchResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
					  name   = "TestAccIntegrationLogAgent_Cloudwatch_Basic"
					  plan   = "penguin-1"
					  region = "amazon-web-services::eu-central-1"
					  tags   = ["vcr-test"]
					}

					resource "cloudamqp_integration_log_agent" "cloudwatch" {
					  instance_id = cloudamqp_instance.instance.id
					  cloudwatch {
					    iam_role        = "%s"
					    iam_external_id = "%s"
					    region          = "eu-central-1"
					    log_group       = "MyLogGroup"
					    log_stream      = "MyLogStream"
					  }
					}
				`, testIAMRole, testIAMExternalID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(cloudwatchResourceName, "cloudwatch.region", "eu-central-1"),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "cloudwatch.log_group", "MyLogGroup"),
					resource.TestCheckResourceAttr(cloudwatchResourceName, "cloudwatch.log_stream", "MyLogStream"),
				),
			},
		},
	})
}

// TestAccIntegrationLogAgent_Coralogix_Basic: Create Coralogix log agent integration, import (ignoring write-only private_key), and update by incrementing private_key_version.
func TestAccIntegrationLogAgent_Coralogix_Basic(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testPrivateKey := "CORALOGIX_SEND_DATA_KEY"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testPrivateKey = os.Getenv("CORALOGIX_SEND_DATA_KEY")
	}

	var (
		instanceResourceName  = "cloudamqp_instance.instance"
		coralogixResourceName = "cloudamqp_integration_log_agent.coralogix"
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
					  name   = "TestAccIntegrationLogAgent_Coralogix_Basic"
					  plan   = "penguin-1"
					  region = "amazon-web-services::eu-central-1"
					  tags   = ["vcr-test"]
					}

					resource "cloudamqp_integration_log_agent" "coralogix" {
					  instance_id = cloudamqp_instance.instance.id
					  coralogix {
					    private_key = "%s"
					    region      = "eu2"
					    application = "cloudamqp"
					    subsystem   = cloudamqp_instance.instance.host
					  }
					}
				`, testPrivateKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccIntegrationLogAgent_Coralogix_Basic"),
					resource.TestCheckResourceAttr(coralogixResourceName, "coralogix.region", "eu2"),
					resource.TestCheckResourceAttr(coralogixResourceName, "coralogix.application", "cloudamqp"),
					resource.TestCheckResourceAttr(coralogixResourceName, "coralogix.private_key_version", "1"),
					resource.TestCheckResourceAttrPair(coralogixResourceName, "coralogix.subsystem", instanceResourceName, "host"),
				),
			},
			{
				ResourceName:            coralogixResourceName,
				ImportStateIdFunc:       testAccImportCombinedStateIdFunc(instanceResourceName, coralogixResourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"coralogix.private_key"},
			},
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
					  name   = "TestAccIntegrationLogAgent_Coralogix_Basic"
					  plan   = "penguin-1"
					  region = "amazon-web-services::eu-central-1"
					  tags   = ["vcr-test"]
					}

					resource "cloudamqp_integration_log_agent" "coralogix" {
					  instance_id = cloudamqp_instance.instance.id
					  coralogix {
					    private_key         = "%s"
					    private_key_version = 2
					    region              = "eu2"
					    application         = "cloudamqp"
					    subsystem           = cloudamqp_instance.instance.host
					  }
					}
				`, testPrivateKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(coralogixResourceName, "coralogix.region", "eu2"),
					resource.TestCheckResourceAttr(coralogixResourceName, "coralogix.application", "cloudamqp"),
					resource.TestCheckResourceAttr(coralogixResourceName, "coralogix.private_key_version", "2"),
				),
			},
		},
	})
}

// TestAccIntegrationLogAgent_Grafana_Basic: Create Grafana Cloud log agent integration, import (ignoring write-only api_token), and update by incrementing api_token_version.
func TestAccIntegrationLogAgent_Grafana_Basic(t *testing.T) {
	t.Parallel()

	// Set sanitized values for playback and use real values for recording
	testAPIToken := "GRAFANA_API_TOKEN"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testAPIToken = os.Getenv("GRAFANA_API_TOKEN")
	}

	var (
		instanceResourceName = "cloudamqp_instance.instance"
		grafanaResourceName  = "cloudamqp_integration_log_agent.grafana"
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
					  name   = "TestAccIntegrationLogAgent_Grafana_Basic"
					  plan   = "penguin-1"
					  region = "amazon-web-services::eu-central-1"
					  tags   = ["vcr-test"]
					}

					resource "cloudamqp_integration_log_agent" "grafana" {
					  instance_id = cloudamqp_instance.instance.id
					  grafana {
					    endpoint            = "https://otlp-gateway-prod-us-central-0.grafana.net/otlp"
					    grafana_instance_id = "123456"
					    api_token           = "%s"
					  }
					}
				`, testAPIToken),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccIntegrationLogAgent_Grafana_Basic"),
					resource.TestCheckResourceAttr(grafanaResourceName, "grafana.endpoint", "https://otlp-gateway-prod-us-central-0.grafana.net/otlp"),
					resource.TestCheckResourceAttr(grafanaResourceName, "grafana.grafana_instance_id", "123456"),
					resource.TestCheckResourceAttr(grafanaResourceName, "grafana.api_token_version", "1"),
				),
			},
			{
				ResourceName:            grafanaResourceName,
				ImportStateIdFunc:       testAccImportCombinedStateIdFunc(instanceResourceName, grafanaResourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"grafana.api_token"},
			},
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
					  name   = "TestAccIntegrationLogAgent_Grafana_Basic"
					  plan   = "penguin-1"
					  region = "amazon-web-services::eu-central-1"
					  tags   = ["vcr-test"]
					}

					resource "cloudamqp_integration_log_agent" "grafana" {
					  instance_id = cloudamqp_instance.instance.id
					  grafana {
					    endpoint            = "https://otlp-gateway-prod-us-central-0.grafana.net/otlp"
					    grafana_instance_id = "123456"
					    api_token           = "%s"
					    api_token_version   = 2
					  }
					}
				`, testAPIToken),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(grafanaResourceName, "grafana.endpoint", "https://otlp-gateway-prod-us-central-0.grafana.net/otlp"),
					resource.TestCheckResourceAttr(grafanaResourceName, "grafana.grafana_instance_id", "123456"),
					resource.TestCheckResourceAttr(grafanaResourceName, "grafana.api_token_version", "2"),
				),
			},
		},
	})
}

// TestAccIntegrationLogAgent_Splunk_Basic: Create Splunk HEC log agent integration, import (ignoring write-only token), and update by incrementing token_version.
func TestAccIntegrationLogAgent_Splunk_Basic(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testToken := "SPLUNK_TOKEN"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testToken = os.Getenv("SPLUNK_TOKEN")
	}

	var (
		instanceResourceName = "cloudamqp_instance.instance"
		splunkResourceName   = "cloudamqp_integration_log_agent.splunk"
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
					  name   = "TestAccIntegrationLogAgent_Splunk_Basic"
					  plan   = "penguin-1"
					  region = "amazon-web-services::eu-central-1"
					  tags   = ["vcr-test"]
					}

					resource "cloudamqp_integration_log_agent" "splunk" {
					  instance_id = cloudamqp_instance.instance.id
					  splunk {
					    endpoint = "https://my-instance.splunkcloud.com:443/services/collector"
					    token        = "%s"
					    source_type  = "cloudamqp"
					  }
					}
				`, testToken),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccIntegrationLogAgent_Splunk_Basic"),
					resource.TestCheckResourceAttr(splunkResourceName, "splunk.endpoint", "https://my-instance.splunkcloud.com:443/services/collector"),
					resource.TestCheckResourceAttr(splunkResourceName, "splunk.source_type", "cloudamqp"),
					resource.TestCheckResourceAttr(splunkResourceName, "splunk.token_version", "1"),
				),
			},
			{
				ResourceName:            splunkResourceName,
				ImportStateIdFunc:       testAccImportCombinedStateIdFunc(instanceResourceName, splunkResourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"splunk.token"},
			},
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
					  name   = "TestAccIntegrationLogAgent_Splunk_Basic"
					  plan   = "penguin-1"
					  region = "amazon-web-services::eu-central-1"
					  tags   = ["vcr-test"]
					}

					resource "cloudamqp_integration_log_agent" "splunk" {
					  instance_id = cloudamqp_instance.instance.id
					  splunk {
					    endpoint   = "https://my-instance.splunkcloud.com:443/services/collector"
					    token          = "%s"
					    source_type    = "cloudamqp"
					    token_version  = 2
					  }
					}
				`, testToken),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(splunkResourceName, "splunk.endpoint", "https://my-instance.splunkcloud.com:443/services/collector"),
					resource.TestCheckResourceAttr(splunkResourceName, "splunk.source_type", "cloudamqp"),
					resource.TestCheckResourceAttr(splunkResourceName, "splunk.token_version", "2"),
				),
			},
		},
	})
}

// TestAccIntegrationLogAgent_Uptrace_Basic: Create Uptrace log agent integration, import (ignoring write-only dsn), and update by incrementing dsn_version.
func TestAccIntegrationLogAgent_Uptrace_Basic(t *testing.T) {
	t.Parallel()

	// Set sanitized value for playback and use real value for recording
	testDSN := "UPTRACE_DSN"
	if os.Getenv("CLOUDAMQP_RECORD") != "" {
		testDSN = os.Getenv("UPTRACE_DSN")
	}

	var (
		instanceResourceName = "cloudamqp_instance.instance"
		uptraceResourceName  = "cloudamqp_integration_log_agent.uptrace"
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
					  name   = "TestAccIntegrationLogAgent_Uptrace_Basic"
					  plan   = "penguin-1"
					  region = "amazon-web-services::eu-central-1"
					  tags   = ["vcr-test"]
					}

					resource "cloudamqp_integration_log_agent" "uptrace" {
					  instance_id = cloudamqp_instance.instance.id
					  uptrace {
					    dsn = "%s"
					  }
					}
				`, testDSN),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccIntegrationLogAgent_Uptrace_Basic"),
					resource.TestCheckResourceAttr(uptraceResourceName, "uptrace.dsn_version", "1"),
				),
			},
			{
				ResourceName:            uptraceResourceName,
				ImportStateIdFunc:       testAccImportCombinedStateIdFunc(instanceResourceName, uptraceResourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"uptrace.dsn"},
			},
			{
				Config: fmt.Sprintf(`
					resource "cloudamqp_instance" "instance" {
					  name   = "TestAccIntegrationLogAgent_Uptrace_Basic"
					  plan   = "penguin-1"
					  region = "amazon-web-services::eu-central-1"
					  tags   = ["vcr-test"]
					}

					resource "cloudamqp_integration_log_agent" "uptrace" {
					  instance_id = cloudamqp_instance.instance.id
					  uptrace {
					    dsn         = "%s"
					    dsn_version = 2
					  }
					}
				`, testDSN),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(uptraceResourceName, "uptrace.dsn_version", "2"),
				),
			},
		},
	})
}

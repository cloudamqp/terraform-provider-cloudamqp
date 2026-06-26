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

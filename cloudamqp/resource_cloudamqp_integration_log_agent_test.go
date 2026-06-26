package cloudamqp

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccIntegrationLogAgent_Cloudwatch_Basic: Create CloudWatch log agent integration, import, and update log_group/log_stream.
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

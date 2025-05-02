package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccIntegrationAwsEventbridge_Basic: Enabled AWS eventbridge integration and import.
func TestAccIntegrationAwsEventbridge_Basic(t *testing.T) {
	var (
		fileNames               = []string{"instance", "integration_aws_eventbridge"}
		instanceResourceName    = "cloudamqp_instance.instance"
		eventbridgeResourceName = "cloudamqp_integration_aws_eventbridge.aws_eventbridge"

		params = map[string]string{
			"InstanceName":            "TestAccIntegrationAwsEventbridge_Basic",
			"InstanceID":              fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":            "bunny-1",
			"AwsEventbridgeVhost":     "myvhost",
			"AwsEventbridgeQueue":     "myqueue",
			"AwsEventbridgeAccountID": "012345678910",
			"AwsEventbridgeRegion":    "us-east-1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(eventbridgeResourceName, "vhost", params["AwsEventbridgeVhost"]),
					resource.TestCheckResourceAttr(eventbridgeResourceName, "queue", params["AwsEventbridgeQueue"]),
					resource.TestCheckResourceAttr(eventbridgeResourceName, "aws_account_id", params["AwsEventbridgeAccountID"]),
					resource.TestCheckResourceAttr(eventbridgeResourceName, "aws_region", params["AwsEventbridgeRegion"]),
					resource.TestCheckResourceAttr(eventbridgeResourceName, "with_headers", "true"),
				),
			},
			{
				ResourceName:      eventbridgeResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, eventbridgeResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

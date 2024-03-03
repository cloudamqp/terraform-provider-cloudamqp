package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccIntegrationAwsEventbridge_Basic: Enabled AWS eventbridge integration and import.
func TestAccIntegrationAwsEventbridge_Basic(t *testing.T) {
	var (
		fileNames    = []string{"instance", "integration_aws_eventbridge"}
		instanceName = "cloudamqp_instance.instance"
		resourceName = "cloudamqp_integration_aws_eventbridge.aws_eventbridge"

		params = map[string]string{
			"InstanceName":            "TestAccIntegrationAwsEventbridge_Basic",
			"InstanceID":              fmt.Sprintf("%s.id", instanceName),
			"AwsEventbridgeVhost":     "myvhost",
			"AwsEventbridgeQueue":     "myqueue",
			"AwsEventbridgeAccountID": "012345678910",
			"AwsEventbridgeRegion":    "us-east-1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(resourceName, "vhost", params["AwsEventbridgeVhost"]),
					resource.TestCheckResourceAttr(resourceName, "queue", params["AwsEventbridgeQueue"]),
					resource.TestCheckResourceAttr(resourceName, "aws_account_id", params["AwsEventbridgeAccountID"]),
					resource.TestCheckResourceAttr(resourceName, "aws_region", params["AwsEventbridgeRegion"]),
					resource.TestCheckResourceAttr(resourceName, "with_headers", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceName, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

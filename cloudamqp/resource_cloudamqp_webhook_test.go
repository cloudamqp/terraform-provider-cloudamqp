package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccWebhook_Basic: Creating dedicated AWS instance, enable webhook integration and import.
func TestAccWebhook_Basic(t *testing.T) {
	var (
		fileNames    = []string{"instance", "webhook"}
		instanceName = "cloudamqp_instance.instance"
		resourceName = "cloudamqp_webhook.webhook_queue"

		params = map[string]string{
			"InstanceName": "TestAccWebhook_Basic",
			"InstanceID":   fmt.Sprintf("%s.id", instanceName),
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
					resource.TestCheckResourceAttr(resourceName, "vhost", "myvhost"),
					resource.TestCheckResourceAttr(resourceName, "queue", "myqueue"),
					resource.TestCheckResourceAttr(resourceName, "webhook_uri", "https://example.com/webhook?key=secret"),
					resource.TestCheckResourceAttr(resourceName, "retry_interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "concurrency", "5"),
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

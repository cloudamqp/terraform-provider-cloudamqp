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
		fileNames            = []string{"instance", "webhook"}
		instanceResourceName = "cloudamqp_instance.instance"
		webhookResourceName  = "cloudamqp_webhook.webhook_queue"

		params = map[string]string{
			"InstanceName": "TestAccWebhook_Basic",
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
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
					resource.TestCheckResourceAttr(webhookResourceName, "vhost", "myvhost"),
					resource.TestCheckResourceAttr(webhookResourceName, "queue", "myqueue"),
					resource.TestCheckResourceAttr(webhookResourceName, "webhook_uri", "https://example.com/webhook?key=secret"),
					resource.TestCheckResourceAttr(webhookResourceName, "retry_interval", "0"),
					resource.TestCheckResourceAttr(webhookResourceName, "concurrency", "5"),
				),
			},
			{
				ResourceName:      webhookResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, webhookResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

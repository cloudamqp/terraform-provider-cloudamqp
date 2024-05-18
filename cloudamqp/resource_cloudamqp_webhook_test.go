package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccWebhook_Basic: Creating dedicated AWS instance, enable webhook integration and import.
func TestAccWebhook_Basic(t *testing.T) {
	var (
		fileNames            = []string{"instance", "webhook"}
		instanceResourceName = "cloudamqp_instance.instance"
		webhookResourceName  = "cloudamqp_webhook.webhook_queue"

		params = map[string]string{
			"InstanceName":       "TestAccWebhook_Basic",
			"InstanceID":         fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":       "bunny-1",
			"WebhookVhost":       fmt.Sprintf("%s.vhost", instanceResourceName),
			"WebhookQueue":       "myqueue",
			"WebhookURI":         "https://example.com/webhook?key=secret",
			"WebhookConcurrency": "1",
		}

		paramsUpdated = map[string]string{
			"InstanceName":       "TestAccWebhook_Basic",
			"InstanceID":         fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":       "bunny-1",
			"WebhookVhost":       fmt.Sprintf("%s.vhost", instanceResourceName),
			"WebhookQueue":       "myqueue_02",
			"WebhookURI":         "https://example.com/webhook?key=secret",
			"WebhookConcurrency": "1",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(webhookResourceName, "queue", params["WebhookQueue"]),
					resource.TestCheckResourceAttr(webhookResourceName, "webhook_uri", params["WebhookURI"]),
					resource.TestCheckResourceAttr(webhookResourceName, "concurrency", params["WebhookConcurrency"]),
				),
			},
			{
				ResourceName:      webhookResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, webhookResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsUpdated["InstanceName"]),
					resource.TestCheckResourceAttr(webhookResourceName, "queue", paramsUpdated["WebhookQueue"]),
					resource.TestCheckResourceAttr(webhookResourceName, "webhook_uri", paramsUpdated["WebhookURI"]),
					resource.TestCheckResourceAttr(webhookResourceName, "concurrency", paramsUpdated["WebhookConcurrency"]),
				),
			},
		},
	})
}

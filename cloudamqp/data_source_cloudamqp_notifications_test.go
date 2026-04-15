package cloudamqp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNotificationsDataSource_Basic: Create notification recipient and verify it is returned in the data source.
func TestAccNotificationsDataSource_Basic(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
          resource "cloudamqp_notification" "email_recipient" {
            instance_id = 1085
            type        = "email"
            value       = "alarm@example.com"
            name        = "alarm"
          }`,
			},
			{
				Config: `
          resource "cloudamqp_notification" "email_recipient" {
            instance_id = 1085
            type        = "email"
            value       = "alarm@example.com"
            name        = "alarm"
          }

          data "cloudamqp_notifications" "notifications" {
            instance_id = 1085
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudamqp_notifications.notifications", "recipients.#", "2"),
					resource.TestCheckResourceAttr("data.cloudamqp_notifications.notifications", "recipients.0.value", "default@example.com"),
					resource.TestCheckResourceAttr("data.cloudamqp_notifications.notifications", "recipients.0.name", "Default"),
					resource.TestCheckResourceAttr("data.cloudamqp_notifications.notifications", "recipients.0.type", "email"),
					resource.TestCheckResourceAttr("data.cloudamqp_notifications.notifications", "recipients.1.value", "alarm@example.com"),
					resource.TestCheckResourceAttr("data.cloudamqp_notifications.notifications", "recipients.1.name", "alarm"),
					resource.TestCheckResourceAttr("data.cloudamqp_notifications.notifications", "recipients.1.type", "email"),
				),
			},
		},
	})
}

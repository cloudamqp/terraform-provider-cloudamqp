package cloudamqp

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNotificationDataSource_Default: Read default notification recipient of an instance.
func TestAccNotificationDataSource_Default(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
          data "cloudamqp_notification" "default_notification" {
            instance_id = 1085
            name        = "Default"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudamqp_notification.default_notification", "type", "email"),
					resource.TestCheckResourceAttr("data.cloudamqp_notification.default_notification", "value", "default@example.com"),
					resource.TestCheckResourceAttr("data.cloudamqp_notification.default_notification", "name", "Default"),
				),
			},
		},
	})
}

// TestAccNotificationDataSource_Identifier: Read notification recipient by identifier.
func TestAccNotificationDataSource_Identifier(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
          data "cloudamqp_notification" "default_notification" {
            instance_id  = 1085
            recipient_id = 830
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.cloudamqp_notification.default_notification", "type", "email"),
					resource.TestCheckResourceAttr("data.cloudamqp_notification.default_notification", "value", "default@example.com"),
					resource.TestCheckResourceAttr("data.cloudamqp_notification.default_notification", "name", "Default"),
				),
			},
		},
	})
}

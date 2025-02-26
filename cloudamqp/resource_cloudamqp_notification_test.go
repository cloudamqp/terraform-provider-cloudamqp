package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNotification_Basic: Create CPU alarm, import and change values.
func TestAccNotification_Basic(t *testing.T) {
	var (
		fileNames                = []string{"instance", "notification"}
		instanceResourceName     = "cloudamqp_instance.instance"
		notificationResourceName = "cloudamqp_notification.recipient"

		params = map[string]string{
			"InstanceName":   "TestAccNotification_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"RecipientType":  "email",
			"RecipientValue": "notification@example.com",
			"RecipientName":  "notification",
		}

		paramsUpdated = map[string]string{
			"InstanceName":   "TestAccNotification_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":   "bunny-1",
			"RecipientType":  "email",
			"RecipientValue": "test@example.com",
			"RecipientName":  "test",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(notificationResourceName, "type", params["RecipientType"]),
					resource.TestCheckResourceAttr(notificationResourceName, "value", params["RecipientValue"]),
					resource.TestCheckResourceAttr(notificationResourceName, "name", params["RecipientName"]),
				),
			},
			{
				ResourceName:      notificationResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, notificationResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(notificationResourceName, "type", paramsUpdated["RecipientType"]),
					resource.TestCheckResourceAttr(notificationResourceName, "value", paramsUpdated["RecipientValue"]),
					resource.TestCheckResourceAttr(notificationResourceName, "name", paramsUpdated["RecipientName"]),
				),
			},
		},
	})
}

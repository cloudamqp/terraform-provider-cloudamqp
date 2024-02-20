package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccNotification_Basic: Create CPU alarm, import and change values.
func TestAccNotification_Basic(t *testing.T) {
	var (
		fileNames    = []string{"instance", "notification"}
		instanceName = "cloudamqp_instance.instance"
		resourceName = "cloudamqp_notification.recipient"

		params = map[string]string{
			"InstanceName":   "TestAccNotification_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceName),
			"RecipientType":  "email",
			"RecipientValue": "notification@example.com",
			"RecipientName":  "notification",
		}

		paramsUpdated = map[string]string{
			"InstanceName":   "TestAccNotification_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceName),
			"RecipientType":  "email",
			"RecipientValue": "test@example.com",
			"RecipientName":  "test",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: loadTemplatedConfig(t, fileNames, params),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(instanceName, "name", params["InstanceName"]),
					resource.TestCheckResourceAttr(resourceName, "type", params["RecipientType"]),
					resource.TestCheckResourceAttr(resourceName, "value", params["RecipientValue"]),
					resource.TestCheckResourceAttr(resourceName, "name", params["RecipientName"]),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccImportStateIdFunc(instanceName, resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: loadTemplatedConfig(t, fileNames, paramsUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", paramsUpdated["RecipientType"]),
					resource.TestCheckResourceAttr(resourceName, "value", paramsUpdated["RecipientValue"]),
					resource.TestCheckResourceAttr(resourceName, "name", paramsUpdated["RecipientName"]),
				),
			},
		},
	})
}

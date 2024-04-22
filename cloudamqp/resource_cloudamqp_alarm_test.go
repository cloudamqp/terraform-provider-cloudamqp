package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// TestAccAlarm_Basic: Create CPU alarm, import and change values.
func TestAccAlarm_Basic(t *testing.T) {
	var (
		fileNames                    = []string{"instance", "notification", "data_source/notification_default", "alarm"}
		instanceResourceName         = "cloudamqp_instance.instance"
		defaultRecipientResourceName = "data.cloudamqp_notification.default_recipient"
		recipientResourceName        = "cloudamqp_notification.recipient"
		noticeAlarmResourceName      = "cloudamqp_alarm.notice"
		cpuAlarmResourceName         = "cloudamqp_alarm.cpu"
		dataCpuAlarmResourceName     = "data.cloudamqp_alarm.cpu"

		params = map[string]string{
			"InstanceName":   "TestAccAlarm_Basic",
			"InstanceID":     fmt.Sprintf("%s.id", instanceResourceName),
			"RecipientName":  "test",
			"RecipientType":  "email",
			"RecipientValue": "test@example.com",
			"NoticeRecipients": fmt.Sprintf("%s, %s",
				fmt.Sprintf("%s.id", defaultRecipientResourceName),
				fmt.Sprintf("%s.id", recipientResourceName),
			),
			"CPUAlarmEnabled":          "true",
			"CPUAlarmTimeThreshold":    "600",
			"CPUAlarmValueThreshold":   "90",
			"CPUAlarmReminderInterval": "0",
			"CPUAlarmRecipients":       fmt.Sprintf("%s.id", recipientResourceName),
		}

		fileNamesUpdated = []string{"instance", "notification", "notification_data", "alarm", "alarm_data"}
		paramsUpdated    = map[string]string{
			"InstanceName":             "TestAccAlarm_Basic",
			"InstanceID":               fmt.Sprintf("%s.id", instanceResourceName),
			"RecipientName":            "test",
			"RecipientType":            "email",
			"RecipientValue":           "test@example.com",
			"NoticeRecipients":         fmt.Sprintf("%s.id", defaultRecipientResourceName),
			"CPUAlarmType":             "cpu",
			"CPUAlarmEnabled":          "true",
			"CPUAlarmTimeThreshold":    "450",
			"CPUAlarmValueThreshold":   "50",
			"CPUAlarmReminderInterval": "0",
			"CPUAlarmRecipients":       fmt.Sprintf("%s.id", recipientResourceName),
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
					resource.TestCheckResourceAttr(recipientResourceName, "name", params["RecipientName"]),
					resource.TestCheckResourceAttr(recipientResourceName, "type", params["RecipientType"]),
					resource.TestCheckResourceAttr(recipientResourceName, "value", params["RecipientValue"]),
					resource.TestCheckResourceAttr(noticeAlarmResourceName, "type", "notice"),
					resource.TestCheckResourceAttr(noticeAlarmResourceName, "recipients.#", "2"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "type", "cpu"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "enabled", params["CPUAlarmEnabled"]),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "time_threshold", params["CPUAlarmTimeThreshold"]),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "value_threshold", params["CPUAlarmValueThreshold"]),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "reminder_interval", params["CPUAlarmReminderInterval"]),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "recipients.#", "1"),
				),
			},
			{
				ResourceName:      cpuAlarmResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, cpuAlarmResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesUpdated, paramsUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(noticeAlarmResourceName, "type", "notice"),
					resource.TestCheckResourceAttr(noticeAlarmResourceName, "recipients.#", "1"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "type", "cpu"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "enabled", paramsUpdated["CPUAlarmEnabled"]),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "time_threshold", paramsUpdated["CPUAlarmTimeThreshold"]),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "value_threshold", paramsUpdated["CPUAlarmValueThreshold"]),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "reminder_interval", paramsUpdated["CPUAlarmReminderInterval"]),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "recipients.#", "1"),
					// validate data sources
					resource.TestCheckResourceAttr(defaultRecipientResourceName, "name", "Default"),
					resource.TestCheckResourceAttr(defaultRecipientResourceName, "type", "email"),
					resource.TestCheckResourceAttr(dataCpuAlarmResourceName, "type", "cpu"),
					resource.TestCheckResourceAttr(dataCpuAlarmResourceName, "enabled", params["CPUAlarmEnabled"]),
					resource.TestCheckResourceAttr(dataCpuAlarmResourceName, "time_threshold", params["CPUAlarmTimeThreshold"]),
					resource.TestCheckResourceAttr(dataCpuAlarmResourceName, "value_threshold", params["CPUAlarmValueThreshold"]),
					resource.TestCheckResourceAttr(dataCpuAlarmResourceName, "reminder_interval", params["CPUAlarmReminderInterval"]),
					resource.TestCheckResourceAttr(dataCpuAlarmResourceName, "recipients.#", "1"),
				),
			},
		},
	})
}

func testAccImportCombinedStateIdFunc(instanceName, resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[instanceName]
		if !ok {
			return "", fmt.Errorf("Resource %s not found", instanceName)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("No resource id set")
		}
		instanceID := rs.Primary.ID

		rs, ok = state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("No resource id set")
		}
		resourceID := rs.Primary.ID
		return fmt.Sprintf("%s,%v", resourceID, instanceID), nil
	}
}

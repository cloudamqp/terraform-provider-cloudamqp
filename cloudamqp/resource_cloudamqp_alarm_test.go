package cloudamqp

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccAlarm_Basic: Create CPU alarm, import and change values.
func TestAccAlarm_Basic(t *testing.T) {
	t.Parallel()

	var (
		instanceResourceName       = "cloudamqp_instance.instance"
		notificationDataSourceName = "data.cloudamqp_notification.default_recipient"
		notificationResourceName   = "cloudamqp_notification.recipient"
		noticeAlarmResourceName    = "cloudamqp_alarm.notice"
		cpuAlarmResourceName       = "cloudamqp_alarm.cpu"
		cpuAlarmDataSourceName     = "data.cloudamqp_alarm.cpu"
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
				  resource "cloudamqp_instance" "instance" {
				    name              = "TestAccAlarm_Basic"
						region            = "amazon-web-services::us-east-1"
				    plan              = "penguin-1"
						tags              = ["vcr-test"]
				  }

					data "cloudamqp_notification" "default_recipient" {
					  instance_id = cloudamqp_instance.instance.id
					  name        = "Default"
					}

				  resource "cloudamqp_notification" "recipient" {
				    instance_id = cloudamqp_instance.instance.id
				    type        = "email"
				    value       = "test@example.com"
						name        = "test"
					}

					resource "cloudamqp_alarm" "notice" {
					  instance_id = cloudamqp_instance.instance.id
					  type        = "notice"
					  enabled     = true
					  recipients  = [data.cloudamqp_notification.default_recipient.id]
					}

					resource "cloudamqp_alarm" "cpu" {
					  instance_id       = cloudamqp_instance.instance.id
					  type              = "cpu"
					  enabled           = true
					  time_threshold    = 600
					  value_threshold   = 90
					  reminder_interval = 0
					  recipients        = [cloudamqp_notification.recipient.id]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccAlarm_Basic"),
					resource.TestCheckResourceAttr(notificationResourceName, "name", "test"),
					resource.TestCheckResourceAttr(notificationResourceName, "type", "email"),
					resource.TestCheckResourceAttr(notificationResourceName, "value", "test@example.com"),
					resource.TestCheckResourceAttr(noticeAlarmResourceName, "type", "notice"),
					resource.TestCheckResourceAttr(noticeAlarmResourceName, "recipients.#", "1"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "type", "cpu"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "time_threshold", "600"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "value_threshold", "90"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "reminder_interval", "0"),
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
				Config: `
				resource "cloudamqp_instance" "instance" {
				    name              = "TestAccAlarm_Basic"
						region            = "amazon-web-services::us-east-1"
				    plan              = "penguin-1"
						tags              = ["vcr-test"]
				  }

					data "cloudamqp_notification" "default_recipient" {
					  instance_id = cloudamqp_instance.instance.id
					  name        = "Default"
					}

				  resource "cloudamqp_notification" "recipient" {
				    instance_id = cloudamqp_instance.instance.id
				    type        = "email"
				    value       = "test@example.com"
						name        = "test"
					}

					resource "cloudamqp_alarm" "notice" {
					  instance_id = cloudamqp_instance.instance.id
					  type        = "notice"
					  enabled     = true
					  recipients  = [data.cloudamqp_notification.default_recipient.id]
					}

					resource "cloudamqp_alarm" "cpu" {
					  instance_id       = cloudamqp_instance.instance.id
					  type              = "cpu"
					  enabled           = true
					  time_threshold    = 450
					  value_threshold   = 50
					  reminder_interval = 0
					  recipients        = [cloudamqp_notification.recipient.id]
					}
						
					data "cloudamqp_alarm" "cpu" {
					  instance_id = cloudamqp_instance.instance.id
					  type        = "cpu"

						depends_on = [
							cloudamqp_alarm.cpu,
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(noticeAlarmResourceName, "type", "notice"),
					resource.TestCheckResourceAttr(noticeAlarmResourceName, "recipients.#", "1"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "type", "cpu"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "time_threshold", "450"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "value_threshold", "50"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "reminder_interval", "0"),
					resource.TestCheckResourceAttr(cpuAlarmResourceName, "recipients.#", "1"),
					// validate data sources
					resource.TestCheckResourceAttr(notificationDataSourceName, "name", "Default"),
					resource.TestCheckResourceAttr(notificationDataSourceName, "type", "email"),
					resource.TestCheckResourceAttr(cpuAlarmDataSourceName, "type", "cpu"),
					resource.TestCheckResourceAttr(cpuAlarmDataSourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(cpuAlarmDataSourceName, "time_threshold", "450"),
					resource.TestCheckResourceAttr(cpuAlarmDataSourceName, "value_threshold", "50"),
					resource.TestCheckResourceAttr(cpuAlarmDataSourceName, "reminder_interval", "0"),
					resource.TestCheckResourceAttr(cpuAlarmDataSourceName, "recipients.#", "1"),
				),
			},
		},
	})
}

func TestAccAlarm_Notice(t *testing.T) {
	t.Parallel()

	var (
		instanceResourceName       = "cloudamqp_instance.instance"
		notificationDataSourceName = "data.cloudamqp_notification.default_recipient"
		noticeAlarmResourceName    = "cloudamqp_alarm.notice"
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
				  resource "cloudamqp_instance" "instance" {
				    name              = "TestAccAlarm_Basic"
						region            = "amazon-web-services::us-east-1"
				    plan              = "penguin-1"
						tags              = ["vcr-test"]
				  }

					data "cloudamqp_notification" "default_recipient" {
					  instance_id = cloudamqp_instance.instance.id
					  name        = "Default"
					}

					resource "cloudamqp_alarm" "notice" {
					  instance_id = cloudamqp_instance.instance.id
					  type        = "notice"
					  enabled     = true
					  recipients  = [data.cloudamqp_notification.default_recipient.id]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", "TestAccAlarm_Basic"),
					resource.TestCheckResourceAttr(notificationDataSourceName, "name", "Default"),
					resource.TestCheckResourceAttr(notificationDataSourceName, "type", "email"),
					resource.TestCheckResourceAttr(noticeAlarmResourceName, "type", "notice"),
					resource.TestCheckResourceAttr(noticeAlarmResourceName, "recipients.#", "1"),
				),
			},
			{
				ResourceName:      noticeAlarmResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, noticeAlarmResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccAlarm_DiskAutoResize: Create a disk_auto_resize alarm and import it.
func TestAccAlarm_DiskAutoResize(t *testing.T) {
	t.Parallel()

	var (
		instanceResourceName      = "cloudamqp_instance.instance"
		resizeAlarmResourceName   = "cloudamqp_alarm.disk_auto_resize"
		resizeAlarmDataSourceName = "data.cloudamqp_alarm.disk_auto_resize"
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
				  resource "cloudamqp_instance" "instance" {
				    name              = "TestAccAlarm_DiskAutoResize"
						region            = "amazon-web-services::us-east-1"
				    plan              = "penguin-1"
						tags              = ["vcr-test"]
				  }

				  resource "cloudamqp_notification" "recipient" {
				    instance_id = cloudamqp_instance.instance.id
				    type        = "email"
				    value       = "test@example.com"
						name        = "test"
					}

					resource "cloudamqp_alarm" "disk_auto_resize" {
					  instance_id       = cloudamqp_instance.instance.id
					  type              = "disk_auto_resize"
					  enabled           = true
					  time_threshold    = 600
					  value_threshold   = 5
					  value_calculation = "percentage"
					  allow_downtime    = true
					  recipients        = [cloudamqp_notification.recipient.id]
					}

					data "cloudamqp_alarm" "disk_auto_resize" {
					  instance_id = cloudamqp_instance.instance.id
					  type        = "disk_auto_resize"

						depends_on = [
							cloudamqp_alarm.disk_auto_resize,
						]
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resizeAlarmResourceName, "type", "disk_auto_resize"),
					resource.TestCheckResourceAttr(resizeAlarmResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resizeAlarmResourceName, "time_threshold", "600"),
					resource.TestCheckResourceAttr(resizeAlarmResourceName, "value_threshold", "5"),
					resource.TestCheckResourceAttr(resizeAlarmResourceName, "value_calculation", "percentage"),
					resource.TestCheckResourceAttr(resizeAlarmResourceName, "allow_downtime", "true"),
					resource.TestCheckResourceAttr(resizeAlarmResourceName, "recipients.#", "1"),
					// validate data source
					resource.TestCheckResourceAttr(resizeAlarmDataSourceName, "type", "disk_auto_resize"),
					resource.TestCheckResourceAttr(resizeAlarmDataSourceName, "value_calculation", "percentage"),
					resource.TestCheckResourceAttr(resizeAlarmDataSourceName, "allow_downtime", "true"),
				),
			},
			{
				ResourceName:      resizeAlarmResourceName,
				ImportStateIdFunc: testAccImportCombinedStateIdFunc(instanceResourceName, resizeAlarmResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccAlarm_DiskFieldsRejectedOnWrongType: disk specific arguments are
// rejected at plan time on alarm types that do not support them.
func TestAccAlarm_DiskFieldsRejectedOnWrongType(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
					resource "cloudamqp_alarm" "invalid" {
					  instance_id     = 123
					  type            = "cpu"
					  enabled         = true
					  time_threshold  = 600
					  value_threshold = 90
					  allow_downtime  = true
					  recipients      = [1]
					}
				`,
				ExpectError: regexp.MustCompile(`allow_downtime can only be set when type is "disk_auto_resize"`),
			},
			{
				Config: `
					resource "cloudamqp_alarm" "invalid" {
					  instance_id       = 123
					  type              = "cpu"
					  enabled           = true
					  time_threshold    = 600
					  value_threshold   = 90
					  value_calculation = "percentage"
					  recipients        = [1]
					}
				`,
				ExpectError: regexp.MustCompile(`value_calculation can only be set when type is "disk" or "disk_auto_resize"`),
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

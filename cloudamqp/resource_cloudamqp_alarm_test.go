package cloudamqp

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"testing"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlarm_Basic(t *testing.T) {
	instanceName := "cloudamqp_instance.instance"
	resourceName := "cloudamqp_alarm.connection_01"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlarmDestroy(instanceName, resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccAlarmConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlarmExist(instanceName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "connection"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "value_threshold", "0"),
					resource.TestCheckResourceAttr(resourceName, "time_threshold", "60"),
				),
			},
			{
				Config: testAccAlarmConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlarmExist(instanceName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "type", "connection"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "value_threshold", "25"),
					resource.TestCheckResourceAttr(resourceName, "time_threshold", "120"),
				),
			},
			{
				Config: testAccAlarmConfigDisable(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlarmExist(instanceName, resourceName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlarmExist(instanceName, resourceName string) resource.TestCheckFunc {
	log.Printf("[DEBUG] resource_alarm::testAccCheckAlarmExist resource: %s", resourceName)
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		alarmID := rs.Primary.ID

		rs, ok = state.RootModule().Resources[instanceName]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set for instance")
		}
		instanceID, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		_, err := api.ReadAlarm(instanceID, alarmID)
		if err != nil {
			return fmt.Errorf("Error fetching item with resource %s. %s", resourceName, err)
		}
		return nil
	}
}

func testAccCheckAlarmDestroy(instanceName, resourceName string) resource.TestCheckFunc {
	log.Printf("[DEBUG] resource_alarm::testAccCheckAlarmDestroy resource: %s", resourceName)
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No record ID is set for the resource")
		}
		alarmID := rs.Primary.ID

		rs, ok = state.RootModule().Resources[instanceName]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No record id is set for the instance")
		}
		instanceID, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		_, err := api.ReadAlarm(instanceID, alarmID)
		if err == nil {
			return fmt.Errorf("Alert still exists")
		}
		notFoundErr := "Invalid ID"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
		return nil
	}
}

func testAccAlarmConfigBasic() string {
	return `
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-alarm-test"
			nodes 			= 1
			plan  			= "bunny-1"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		data "cloudamqp_notification" "default_recipient" {
			instance_id = cloudamqp_instance.instance.id
			name 				= "Default"
		}

		resource "cloudamqp_alarm" "connection_01" {
			instance_id 			= cloudamqp_instance.instance.id
			type 							= "connection"
			enabled						=  true
			value_threshold 	= 0
			time_threshold 		= 60
			recipients = [data.cloudamqp_notification.default_recipient.id]
		}
		`
}

func testAccAlarmConfigUpdate() string {
	return `
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-alarm-test"
			nodes 			= 1
			plan  			= "bunny-1"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		data "cloudamqp_notification" "default_recipient" {
			instance_id = cloudamqp_instance.instance.id
			name 				= "Default"
		}

		resource "cloudamqp_alarm" "connection_01" {
			instance_id 			= cloudamqp_instance.instance.id
			type 							= "connection"
			enabled 					= true
			value_threshold 	= 25
			time_threshold 		= 120
			recipients = [data.cloudamqp_notification.default_recipient.id]
		}
		`
}

func testAccAlarmConfigDisable() string {
	return `
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-alarm-test"
			nodes 			= 1
			plan  			= "bunny-1"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		data "cloudamqp_notification" "default_recipient" {
			instance_id = cloudamqp_instance.instance.id
			name 				= "Default"
		}

		resource "cloudamqp_alarm" "connection_01" {
			instance_id 			= cloudamqp_instance.instance.id
			type 							= "connection"
			enabled 					= false
			value_threshold 	= 25
			time_threshold 		= 120
			recipients = [data.cloudamqp_notification.default_recipient.id]
		}
		`
}

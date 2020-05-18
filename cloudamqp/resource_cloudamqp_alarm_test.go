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
	instance_name := "cloudamqp_instance.instance"
	resource_name := "cloudamqp_alarm.connection_01"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlarmDestroy(instance_name, resource_name),
		Steps: []resource.TestStep{
			{
				Config: testAccAlarmConfig_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlarmExist(instance_name, resource_name),
					resource.TestCheckResourceAttr(resource_name, "type", "connection"),
					resource.TestCheckResourceAttr(resource_name, "enabled", "true"),
					resource.TestCheckResourceAttr(resource_name, "value_threshold", "0"),
					resource.TestCheckResourceAttr(resource_name, "time_threshold", "60"),
				),
			},
			{
				Config: testAccAlarmConfig_Update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlarmExist(instance_name, resource_name),
					resource.TestCheckResourceAttr(resource_name, "type", "connection"),
					resource.TestCheckResourceAttr(resource_name, "enabled", "true"),
					resource.TestCheckResourceAttr(resource_name, "value_threshold", "25"),
					resource.TestCheckResourceAttr(resource_name, "time_threshold", "120"),
				),
			},
			{
				Config: testAccAlarmConfig_Disable(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlarmExist(instance_name, resource_name),
					resource.TestCheckResourceAttr(resource_name, "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlarmExist(instance_name, resource_name string) resource.TestCheckFunc {
	log.Printf("[DEBUG] resource_alarm::testAccCheckAlarmExist resource: %s", resource_name)
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resource_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		id := rs.Primary.ID

		rs, ok = state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set for instance")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		_, err := api.ReadAlarm(instance_id, id)
		if err != nil {
			return fmt.Errorf("Error fetching item with resource %s. %s", resource_name, err)
		}
		return nil
	}
}

func testAccCheckAlarmDestroy(instance_name, resource_name string) resource.TestCheckFunc {
	log.Printf("[DEBUG] resource_alarm::testAccCheckAlarmDestroy resource: %s", resource_name)
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource: %s not found", resource_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No record ID is set for the resource")
		}
		resource_id := rs.Primary.ID

		rs, ok = state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Instance resource not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No record id is set for the instance")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		api := testAccProvider.Meta().(*api.API)
		_, err := api.ReadAlarm(instance_id, resource_id)
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

func testAccAlarmConfig_Basic() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-alarm-test"
			nodes 			= 1
			plan  			= "bunny"
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
		`)
}

func testAccAlarmConfig_Update() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-alarm-test"
			nodes 			= 1
			plan  			= "bunny"
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
		`)
}

func testAccAlarmConfig_Disable() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-alarm-test"
			nodes 			= 1
			plan  			= "bunny"
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
		`)
}

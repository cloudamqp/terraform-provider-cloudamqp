package cloudamqp

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceAlarmDefault_Basic(t *testing.T) {
	instanceName := "cloudamqp_instance.instance"
	cpuResource := "data.cloudamqp_alarm.default_cpu"
	memoryResource := "data.cloudamqp_alarm.default_memory"
	diskResource := "data.cloudamqp_alarm.default_disk"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ResourceName:  instanceName,
				ImportStateId: "412",
				ImportState:   true,
				//ImportStateVerify: 	true
			},
			{
				Config: testAccAlarmDefaultDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlarmDataSourceExists(instanceName, cpuResource),
					resource.TestCheckResourceAttr(cpuResource, "type", "cpu"),
					resource.TestCheckResourceAttr(cpuResource, "time_threshold", "600"),
					resource.TestCheckResourceAttr(cpuResource, "value_threshold", "90"),
				),
			},
			{
				Config: testAccAlarmDefaultDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlarmDataSourceExists(instanceName, memoryResource),
					resource.TestCheckResourceAttr(memoryResource, "type", "memory"),
					resource.TestCheckResourceAttr(memoryResource, "time_threshold", "600"),
					resource.TestCheckResourceAttr(memoryResource, "value_threshold", "90"),
				),
			},
			{
				Config: testAccAlarmDefaultDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlarmDataSourceExists(instanceName, diskResource),
					resource.TestCheckResourceAttr(diskResource, "type", "disk"),
					resource.TestCheckResourceAttr(diskResource, "time_threshold", "600"),
					resource.TestCheckResourceAttr(diskResource, "value_threshold", "5"),
				),
			},
		},
	})
}

func testAccCheckAlarmDataSourceExists(instanceName, resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[instanceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", instanceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No resource id set")
		}
		instanceID, _ := strconv.Atoi(rs.Primary.ID)

		rs, ok = state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource %s not found", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No resource id set")
		}
		alarmID := rs.Primary.ID

		api := testAccProvider.Meta().(*api.API)
		_, err := api.ReadAlarm(instanceID, alarmID)
		if err != nil {
			return fmt.Errorf("Failed to fetch instance: %v", err)
		}

		return nil
	}
}

func testAccAlarmDefaultDataSourceConfigBasic() string {
	return `
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-alarm-ds-test"
			nodes 			= 1
			plan  			= "bunny-1"
			region 			= "amazon-web-services::eu-north-1"
			rmq_version = "3.8.2"
			tags 				= ["terraform"]
		}

		data "cloudamqp_alarm" "default_cpu" {
			instance_id = cloudamqp_instance.instance.id
			type 				= "cpu"
		}

		data "cloudamqp_alarm" "default_memory" {
			instance_id = cloudamqp_instance.instance.id
			type 				= "memory"
		}

		data "cloudamqp_alarm" "default_disk" {
			instance_id = cloudamqp_instance.instance.id
			type 				= "disk"
		}
	`
}

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
	instance_name := "cloudamqp_instance.instance"
	cpu_resource := "data.cloudamqp_alarm.default_cpu"
	memory_resource := "data.cloudamqp_alarm.default_memory"
	disk_resource := "data.cloudamqp_alarm.default_disk"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlarmDefaultDataSourceConfig_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlarmDataSourceExists(instance_name, cpu_resource),
					resource.TestCheckResourceAttr(cpu_resource, "type", "cpu"),
					resource.TestCheckResourceAttr(cpu_resource, "time_threshold", "600"),
					resource.TestCheckResourceAttr(cpu_resource, "value_threshold", "90"),
				),
			},
			{
				Config: testAccAlarmDefaultDataSourceConfig_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlarmDataSourceExists(instance_name, memory_resource),
					resource.TestCheckResourceAttr(memory_resource, "type", "memory"),
					resource.TestCheckResourceAttr(memory_resource, "time_threshold", "600"),
					resource.TestCheckResourceAttr(memory_resource, "value_threshold", "90"),
				),
			},
			{
				Config: testAccAlarmDefaultDataSourceConfig_Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlarmDataSourceExists(instance_name, disk_resource),
					resource.TestCheckResourceAttr(disk_resource, "type", "disk"),
					resource.TestCheckResourceAttr(disk_resource, "time_threshold", "600"),
					resource.TestCheckResourceAttr(disk_resource, "value_threshold", "5"),
				),
			},
		},
	})
}

func testAccCheckAlarmDataSourceExists(instance_name, resource_name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[instance_name]
		if !ok {
			return fmt.Errorf("Resource %s not found", instance_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No resource id set")
		}
		instance_id, _ := strconv.Atoi(rs.Primary.ID)

		rs, ok = state.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("Resource %s not found", resource_name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No resource id set")
		}
		alarm_id := rs.Primary.ID

		api := testAccProvider.Meta().(*api.API)
		_, err := api.ReadAlarm(instance_id, alarm_id)
		if err != nil {
			return fmt.Errorf("Failed to fetch instance: %v", err)
		}

		return nil
	}
}

func testAccAlarmDefaultDataSourceConfig_Basic() string {
	return fmt.Sprintf(`
		resource "cloudamqp_instance" "instance" {
			name 				= "terraform-alarm-ds-test"
			nodes 			= 1
			plan  			= "bunny"
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

		data "cloudamqp_alarm" "default_memory" {
			instance_id = cloudamqp_instance.instance.id
			type 				= "disk"
		}
	`)
}

package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMaintenanceWindow_LavinMQ(t *testing.T) {
	t.Parallel()

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: `
          resource "cloudamqp_instance" "instance" {
            name   = "TestAccMaintenanceWindow_LavinMQ"
            region = "amazon-web-services::us-east-1"
            plan   = "penguin-1"
            tags   = ["vcr-test"]
				  }
						
          resource "cloudamqp_maintenance_window" "this" {
            instance_id    = cloudamqp_instance.instance.id
            preferred_day  = "Monday"
            preferred_time = "01:00"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_maintenance_window.this", "preferred_day", "Monday"),
					resource.TestCheckResourceAttr("cloudamqp_maintenance_window.this", "preferred_time", "01:00"),
					resource.TestCheckResourceAttr("cloudamqp_maintenance_window.this", "automatic_updates", "on"),
				),
			},
			{
				Config: `
          resource "cloudamqp_instance" "instance" {
            name   = "TestAccMaintenanceWindow_LavinMQ"
            region = "amazon-web-services::us-east-1"
            plan   = "penguin-1"
            tags   = ["vcr-test"]
				  }
						
          resource "cloudamqp_maintenance_window" "this" {
            instance_id   = cloudamqp_instance.instance.id
            preferred_day = "Tuesday"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_maintenance_window.this", "preferred_day", "Tuesday"),
					resource.TestCheckResourceAttr("cloudamqp_maintenance_window.this", "preferred_time", ""),
					resource.TestCheckResourceAttr("cloudamqp_maintenance_window.this", "automatic_updates", "on"),
				),
			},
			{
				Config: `
          resource "cloudamqp_instance" "instance" {
            name   = "TestAccMaintenanceWindow_LavinMQ"
            region = "amazon-web-services::us-east-1"
            plan   = "penguin-1"
            tags   = ["vcr-test"]
          }
	
          resource "cloudamqp_maintenance_window" "this" {
            instance_id    = cloudamqp_instance.instance.id
            preferred_time = "02:00"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_maintenance_window.this", "preferred_day", ""),
					resource.TestCheckResourceAttr("cloudamqp_maintenance_window.this", "preferred_time", "02:00"),
					resource.TestCheckResourceAttr("cloudamqp_maintenance_window.this", "automatic_updates", "on"),
				),
			},
			{
				Config: `
          resource "cloudamqp_instance" "instance" {
            name   = "TestAccMaintenanceWindow_LavinMQ"
            region = "amazon-web-services::us-east-1"
            plan   = "penguin-1"
            tags   = ["vcr-test"]
				  }
						
          resource "cloudamqp_maintenance_window" "this" {
            instance_id       = cloudamqp_instance.instance.id
            preferred_day     = "Monday"
            preferred_time    = "01:00"
            automatic_updates = "off"
          }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cloudamqp_maintenance_window.this", "preferred_day", "Monday"),
					resource.TestCheckResourceAttr("cloudamqp_maintenance_window.this", "preferred_time", "01:00"),
					resource.TestCheckResourceAttr("cloudamqp_maintenance_window.this", "automatic_updates", "off"),
				),
			},
		},
	})
}

func TestAccMaintenanceWindow_RabbitMQ(t *testing.T) {
	t.Parallel()

	var (
		instanceResourceName          = "cloudamqp_instance.instance"
		maintenanceWindowResourceName = "cloudamqp_maintenance_window.this"

		name = "TestAccMaintenanceWindow_RabbitMQ"
		plan = "bunny-1"

		fileNamesDayTime = []string{"instance", "maintenance/set_day_and_time"}
		paramsDayTime    = map[string]string{
			"InstanceName":  name,
			"InstanceID":    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":  plan,
			"PreferredDay":  "Monday",
			"PreferredTime": "01:00",
		}

		fileNamesDay = []string{"instance", "maintenance/set_only_day"}
		paramsDay    = map[string]string{
			"InstanceName": name,
			"InstanceID":   fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan": plan,
			"PreferredDay": "Tuesday",
		}

		fileNamesTime = []string{"instance", "maintenance/set_only_time"}
		paramsTime    = map[string]string{
			"InstanceName":  name,
			"InstanceID":    fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":  plan,
			"PreferredTime": "02:00",
		}
	)

	cloudamqpResourceTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesDayTime, paramsDayTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsDayTime["InstanceName"]),
					resource.TestCheckResourceAttr(maintenanceWindowResourceName, "preferred_day",
						paramsDayTime["PreferredDay"]),
					resource.TestCheckResourceAttr(maintenanceWindowResourceName, "preferred_time",
						paramsDayTime["PreferredTime"]),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesDay, paramsDay),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsDay["InstanceName"]),
					resource.TestCheckResourceAttr(maintenanceWindowResourceName, "preferred_day",
						paramsDay["PreferredDay"]),
					resource.TestCheckResourceAttr(maintenanceWindowResourceName, "preferred_time", ""),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesTime, paramsTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsTime["InstanceName"]),
					resource.TestCheckResourceAttr(maintenanceWindowResourceName, "preferred_day", ""),
					resource.TestCheckResourceAttr(maintenanceWindowResourceName, "preferred_time",
						paramsTime["PreferredTime"]),
				),
			},
		},
	})
}

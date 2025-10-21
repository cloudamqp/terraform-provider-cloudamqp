package cloudamqp

import (
	"fmt"
	"testing"

	"github.com/cloudamqp/terraform-provider-cloudamqp/cloudamqp/vcr-testing/configuration"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMaintenanceWindow_LavinMQ(t *testing.T) {
	t.Parallel()

	var (
		instanceResourceName          = "cloudamqp_instance.instance"
		maintenanceWindowResourceName = "cloudamqp_maintenance_window.this"

		name = "TestAccMaintenanceWindow_LavinMQ"
		plan = "penguin-1"

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

		fileNamesAutomaticUpdates = []string{"instance", "maintenance/set_maintenance"}
		paramsAutomaticUpdates    = map[string]string{
			"InstanceName":     name,
			"InstanceID":       fmt.Sprintf("%s.id", instanceResourceName),
			"InstancePlan":     plan,
			"PreferredDay":     "Monday",
			"PreferredTime":    "01:00",
			"AutomaticUpdates": "on",
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
					resource.TestCheckResourceAttr(maintenanceWindowResourceName, "automatic_updates", "off"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesTime, paramsTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", paramsTime["InstanceName"]),
					resource.TestCheckResourceAttr(maintenanceWindowResourceName, "preferred_day", ""),
					resource.TestCheckResourceAttr(maintenanceWindowResourceName, "preferred_time",
						paramsTime["PreferredTime"]),
					resource.TestCheckResourceAttr(maintenanceWindowResourceName, "automatic_updates", "off"),
				),
			},
			{
				Config: configuration.GetTemplatedConfig(t, fileNamesAutomaticUpdates, paramsAutomaticUpdates),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name",
						paramsAutomaticUpdates["InstanceName"]),
					resource.TestCheckResourceAttr(maintenanceWindowResourceName, "preferred_day",
						paramsAutomaticUpdates["PreferredDay"]),
					resource.TestCheckResourceAttr(maintenanceWindowResourceName, "preferred_time",
						paramsAutomaticUpdates["PreferredTime"]),
					resource.TestCheckResourceAttr(maintenanceWindowResourceName, "automatic_updates",
						paramsAutomaticUpdates["AutomaticUpdates"]),
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

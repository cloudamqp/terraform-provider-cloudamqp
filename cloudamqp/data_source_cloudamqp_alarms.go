package cloudamqp

import (
	"context"
	"fmt"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlarms() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlarmsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Type of the alarm",
				ValidateDiagFunc: validateAlarmType(),
			},
			"alarms": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of alarms",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Alarm identifier",
						},
						"type": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      "Type of the alarm",
							ValidateDiagFunc: validateAlarmType(),
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable or disable an alarm",
						},
						"reminder_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The reminder interval (in seconds) to resend the alarm if not resolved. Set to 0 for no reminders",
						},
						"value_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "What value to trigger the alarm for",
						},
						"value_calculation": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Disk value threshold calculation. Fixed or percentage of disk space remaining",
						},
						"time_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "For how long (in seconds) the value_threshold should be active before trigger alarm",
						},
						"vhost_regex": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Regex for which vhost the queues are in",
						},
						"queue_regex": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Regex for which queues to check",
						},
						"message_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message types (total, unacked, ready) of the queue to trigger the alarm",
						},
						"recipients": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Description: "Identifiers for recipients to be notified.",
						},
					},
				},
			},
		},
	}
}

func dataSourceAlarmsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		data       []map[string]any
		instanceID = d.Get("instance_id").(int)
		err        error
	)

	api := meta.(*api.API)
	data, err = api.ListAlarms(ctx, instanceID)

	if err != nil {
		return diag.FromErr(err)
	}

	if d.Get("type") != "" {
		d.SetId(fmt.Sprintf("%v.%v.alarms", instanceID, d.Get("type")))
		filteredAlarms := make([]map[string]any, 0)
		for _, alarm := range data {
			if alarm["type"] == d.Get("type") {
				filteredAlarms = append(filteredAlarms, alarm)
			}
		}
		data = filteredAlarms
	} else {
		d.SetId(fmt.Sprintf("%v.alarms", instanceID))
	}

	alarms := make([]map[string]any, len(data))
	for k, v := range data {
		alarms[k] = readAlarm(v)
	}

	if err = d.Set("alarms", alarms); err != nil {
		return diag.Errorf("error setting alarms for resource %s: %s", d.Id(), err)
	}

	return diag.Diagnostics{}
}

func readAlarm(data map[string]any) map[string]any {
	alarm := make(map[string]any)
	for k, v := range data {
		if validateAlarmSchemaAttribute(k) {
			alarm[k] = v
		}
	}
	return alarm
}

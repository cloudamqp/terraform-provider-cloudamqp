package cloudamqp

import (
	"context"
	"fmt"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/monitoring"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the alarm",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
					"cpu", "memory", "disk", "queue", "connection",
					"flow", "consumer", "netsplit", "ssh", "notice", "server_unreachable",
				}, false)),
			},
			"alarms": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of alarms",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Alarm identifier",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the alarm",
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
							Computed:    true,
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
		instanceID = int64(d.Get("instance_id").(int))
	)

	client := meta.(*api.API)
	data, err := client.ListAlarms(ctx, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	filtered := data
	if alarmType := d.Get("type").(string); alarmType != "" {
		d.SetId(fmt.Sprintf("%d.%s.alarms", instanceID, alarmType))
		filtered = make([]model.AlarmResponse, 0)
		for _, alarm := range data {
			if alarm.Type == alarmType {
				filtered = append(filtered, alarm)
			}
		}
	} else {
		d.SetId(fmt.Sprintf("%d.alarms", instanceID))
	}

	alarms := make([]map[string]any, len(filtered))
	for k, v := range filtered {
		alarms[k] = readAlarm(v)
	}

	if err = d.Set("alarms", alarms); err != nil {
		return diag.Errorf("error setting alarms for resource %s: %s", d.Id(), err)
	}

	return diag.Diagnostics{}
}

func readAlarm(data model.AlarmResponse) map[string]any {
	alarm := map[string]any{
		"alarm_id": data.ID,
		"type":     data.Type,
		"enabled":  data.Enabled,
	}

	if data.ReminderInterval != nil {
		alarm["reminder_interval"] = *data.ReminderInterval
	} else {
		alarm["reminder_interval"] = int64(0)
	}

	if data.ValueThreshold != nil {
		alarm["value_threshold"] = *data.ValueThreshold
	} else {
		alarm["value_threshold"] = int64(0)
	}

	if data.ValueCalculation != nil {
		alarm["value_calculation"] = *data.ValueCalculation
	} else {
		alarm["value_calculation"] = ""
	}

	if data.TimeThreshold != nil {
		alarm["time_threshold"] = *data.TimeThreshold
	} else {
		alarm["time_threshold"] = int64(0)
	}

	if data.VhostRegex != nil {
		alarm["vhost_regex"] = *data.VhostRegex
	} else {
		alarm["vhost_regex"] = ""
	}

	if data.QueueRegex != nil {
		alarm["queue_regex"] = *data.QueueRegex
	} else {
		alarm["queue_regex"] = ""
	}

	if data.MessageType != nil {
		alarm["message_type"] = *data.MessageType
	} else {
		alarm["message_type"] = ""
	}

	if data.Recipients != nil {
		alarm["recipients"] = *data.Recipients
	} else {
		alarm["recipients"] = []int64{}
	}

	return alarm
}

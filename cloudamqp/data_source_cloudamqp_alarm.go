package cloudamqp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlarm() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlarmRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
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
	}
}

func dataSourceAlarmRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		data       map[string]any
		instanceID = d.Get("instance_id").(int)
		err        error
	)

	// Multiple purpose read. To be used when using data source either by declaring alarm id or type.
	if d.Get("alarm_id") != 0 {
		data, err = dataSourceAlarmIDRead(ctx, instanceID, d.Get("alarm_id").(int), meta)
	} else if d.Get("type") != "" {
		data, err = dataSourceAlarmTypeRead(ctx, instanceID, d.Get("type").(string), meta)
	}

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%v", data["id"]))
	d.Set("alarm_id", data["id"])
	for k, v := range data {
		if validateAlarmSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	return diag.Diagnostics{}
}

func dataSourceAlarmIDRead(ctx context.Context, instanceID int, alarmID int, meta any) (
	map[string]any, error) {

	api := meta.(*api.API)
	id := strconv.Itoa(alarmID)
	alarm, err := api.ReadAlarm(ctx, instanceID, id)
	return alarm, err
}

func dataSourceAlarmTypeRead(ctx context.Context, instanceID int, alarmType string,
	meta any) (map[string]any, error) {

	api := meta.(*api.API)
	alarms, err := api.ListAlarms(ctx, instanceID)

	if err != nil {
		return nil, err
	}
	for _, alarm := range alarms {
		if alarm["type"] == alarmType {
			return alarm, nil
		}
	}
	return nil, nil
}

package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlarm() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmCreate,
		ReadContext:   resourceAlarmRead,
		UpdateContext: resourceAlarmUpdate,
		DeleteContext: resourceAlarmDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Instance identifier",
			},
			"type": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Type of the alarm, valid options are: cpu, memory, disk_usage, queue_length, connection_count, consumers_count, net_split",
				ValidateDiagFunc: validateAlarmType(),
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable or disable an alarm",
			},
			"reminder_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The reminder interval (in seconds) to resend the alarm if not resolved. Set to 0 for no reminders. The Default is 0.",
			},
			"value_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "What value to trigger the alarm for",
			},
			"value_calculation": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Disk value threshold calculation. Fixed or percentage of disk space remaining",
				ValidateDiagFunc: validateValueCalculation(),
			},
			"time_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "For how long (in seconds) the value_threshold should be active before trigger alarm",
			},
			"vhost_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Regex for which vhost the queues are in",
			},
			"queue_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Regex for which queues to check",
			},
			"message_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Message types (total, unacked, ready) of the queue to trigger the alarm",
				ValidateDiagFunc: validateMessageType(),
			},
			"recipients": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Identifiers for recipients to be notified.",
			},
		},
	}
}

func resourceAlarmCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	keys := alarmAttributeKeys()
	params := make(map[string]any)
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	if d.Get("type") == "notice" {
		tflog.Info(ctx, "alarm type is 'notice', skip creation, retrieve existing alarm and update")
		alarms, err := api.ListAlarms(ctx, d.Get("instance_id").(int))

		if err != nil {
			return diag.FromErr(err)
		}

		for _, alarm := range alarms {
			if alarm["type"] == "notice" {
				d.SetId(strconv.FormatFloat(alarm["id"].(float64), 'f', 0, 64))
				tflog.Debug(ctx, fmt.Sprintf("retrieve existing 'notice' alarm with identifier: %s and invoke an update",
					d.Id()))
				return resourceAlarmUpdate(ctx, d, meta)
			}
		}

		return diag.Errorf("couldn't find notice alarm for instance_id: %s", d.Get("instance_id"))
	} else {
		data, err := api.CreateAlarm(ctx, d.Get("instance_id").(int), params)
		if err != nil {
			return diag.FromErr(err)
		}
		if data["id"] != nil {
			d.SetId(data["id"].(string))
		}
	}

	return resourceAlarmRead(ctx, d, meta)
}

func resourceAlarmRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if strings.Contains(d.Id(), ",") {
		tflog.Debug(ctx, fmt.Sprintf("import alarm from input identifier: %s", d.Id()))
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		instanceID, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
	}
	if d.Get("instance_id").(int) == 0 {
		return diag.Errorf("missing input identifier for import: {resource_id},{instance_id}")
	}

	api := meta.(*api.API)
	data, err := api.ReadAlarm(ctx, d.Get("instance_id").(int), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range data {
		if validateAlarmSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	return nil
}

func resourceAlarmUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	keys := alarmAttributeKeys()
	params := make(map[string]any)
	params["id"] = d.Id()

	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	if err := api.UpdateAlarm(ctx, d.Get("instance_id").(int), params); err != nil {
		return diag.FromErr(err)
	}

	return resourceAlarmRead(ctx, d, meta)
}

func resourceAlarmDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if d.Get("type") == "notice" {
		tflog.Debug(ctx, "alarm type is 'notice', skip deletion and just remove from state")
		return diag.Diagnostics{}
	}

	api := meta.(*api.API)
	if err := api.DeleteAlarm(ctx, d.Get("instance_id").(int), d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func validateAlarmType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"cpu",
		"memory",
		"disk",
		"queue",
		"connection",
		"flow",
		"consumer",
		"netsplit",
		"ssh",
		"notice",
		"server_unreachable",
	}, true))
}

func validateAlarmSchemaAttribute(key string) bool {
	switch key {
	case "type",
		"enabled",
		"reminder_interval",
		"value_threshold",
		"value_calculation",
		"time_threshold",
		"vhost_regex",
		"queue_regex",
		"message_type",
		"recipients":
		return true
	}
	return false
}

func validateMessageType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"total",
		"unacked",
		"ready",
	}, true))
}

func validateValueCalculation() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"fixed",
		"percentage",
	}, true))
}

func alarmAttributeKeys() []string {
	return []string{
		"type",
		"enabled",
		"reminder_interval",
		"value_threshold",
		"value_calculation",
		"time_threshold",
		"vhost_regex",
		"queue_regex",
		"message_type",
		"recipients",
	}
}

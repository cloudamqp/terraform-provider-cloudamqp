package cloudamqp

import (
	"fmt"
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlarm() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlarmRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"alarm_id": {
				Type:        schema.TypeInt,
				Required:    true,
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
			"value_threshold": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "What value to trigger the alarm for",
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

func dataSourceAlarmRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	alarm_id := strconv.Itoa(d.Get("alarm_id").(int))
	data, err := api.ReadAlarm(d.Get("instance_id").(int), alarm_id)

	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%v", data["id"]))
	for k, v := range data {
		if validateAlarmSchemaAttribute(k) {
			d.Set(k, v)
		}
	}

	return nil
}

package main

import (
	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
)

func resourceAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlarmCreate,
		Read:   resourceAlarmRead,
		Update: resourceAlarmUpdate,
		Delete: resourceAlarmDelete,
		Schema: map[string]*schema.Schema{
			"instance_id" : {
				Type:				schema.TypeInt,
				Required: 	true,
				Description: "Instance identifier",
			},
			"alarm_id": {
				Type:				schema.TypeInt,
				Optional: 	true,
				Description: "Alarm identifier",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the alarm, valid options are: cpu, memory, disk_usage, queue_length, connection_count, consumers_count, net_split",
			},
			"value_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "What value to trigger the alarm for",
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
			"notifications": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
          Type:         schema.TypeString,
        },
				Description: "Identifiers for recipients to be notified. Leave empty to notifiy all recipients.",
			},
		},
	}
}

func resourceAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"type", "value_threshold", "time_threshold", "vhost_regex", "queue_regex", "notifications"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	data, err := api.CreateAlarm(d.Get("instance_id").(int), params)
	if err != nil {
		return err
	}

	if data["id"] != nil {
		d.SetId(data["id"].(string))
	}
	for k, v := range data {
		if k == "id" {
			continue
		}
		d.Set(k, v)
	}
	return nil
}

func resourceAlarmRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	fmt.Println("resourceAlarmRead: " + d.Id())
	data, err := api.ReadAlarm(d.Get("instance_id").(int), d.Id())
	if err != nil {
		return err
	}
	for k, v := range data {
		d.Set(k, v)
	}
	return nil
}

func resourceAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"type", "value_threshold", "time_threshold", "vhost_regex", "queue_regex", "notifications"}
	params := make(map[string]interface{})
	params["alarm_id"] = d.Id()
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}
	return api.UpdateAlarm(d.Get("instance_id").(int), params)
}

func resourceAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	params := make(map[string]interface{})
	params["alarm_id"] = d.Id()
	return api.DeleteAlarm(d.Get("instance_id").(int), params)
}

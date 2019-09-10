package main

import (
	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNotification() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationCreate,
		Read:   resourceNotificationRead,
		Update: resourceNotificationUpdate,
		Delete: resourceNotificationDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type: 			schema.TypeInt,
				Required: 	true,
				Description: "Instance identifier",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the alarm, valid options are: cpu, memory, disk_usage, queue_length, connection_count, consumers_count, net_split",
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "What value to trigger the alarm for",
			},
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "For how long (in seconds) the value_threshold should be active before trigger alarm",
			},
			"last_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Regex for which vhost the queues are in",
			},
		},
	}
}

func resourceNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"type", "value"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}
	data, err := api.CreateNotification(d.Get("instance_id").(int), params)
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

func resourceNotificationRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data, err := api.ReadNotification(d.Get("instance_id").(int), d.Id())
	if err != nil {
		return err
	}
	for k, v := range data {
		d.Set(k, v)
	}
	return nil
}

func resourceNotificationUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"type", "value"}
	params := make(map[string]interface{})
	params["id"] = d.Id()
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}
	return api.UpdateNotification(d.Get("instance_id").(int), params)
}

func resourceNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	params := make(map[string]interface{})
	params["id"] = d.Id()
	return api.DeleteNotification(d.Get("instance_id").(int), params)
}

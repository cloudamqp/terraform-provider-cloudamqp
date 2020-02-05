package cloudamqp

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlarmCreate,
		Read:   resourceAlarmRead,
		Update: resourceAlarmUpdate,
		Delete: resourceAlarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Type of the alarm, valid options are: cpu, memory, disk_usage, queue_length, connection_count, consumers_count, net_split",
				ValidateFunc: validateAlarmType(),
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
			"notification_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Identifiers for recipients to be notified. Leave empty to notifiy all recipients.",
			},
		},
	}
}

func resourceAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"type", "value_threshold", "time_threshold", "vhost_regex", "queue_regex", "notification_ids"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}
	log.Printf("[DEBUG] cloudamqp::resource::alarm::create params: %v", params)

	var notificationIDs []int
	if attr := d.Get("notification_ids").([]interface{}); len(attr) > 0 {
		for _, v := range attr {
			val, ok := v.(int)
			if ok {
				notificationIDs = append(notificationIDs, val)
			}
		}
	}
	params["notifications"] = notificationIDs
	delete(params, "notification_ids")

	data, err := api.CreateAlarm(d.Get("instance_id").(int), params)
	log.Printf("[DEBUG] cloudamqp::resource::alarm::create data: %v", data)

	if err != nil {
		return err
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
		log.Printf("[DEBUG] cloudamqp::resource::alarm::create id set: %v", d.Id())
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
	if strings.Contains(d.Id(), ",") {
		log.Printf("[DEBUG] cloudamqp::resource::alarm::read id contains ,: %v", d.Id())
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		instance_id, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instance_id)
	}
	if d.Get("instance_id").(int) == 0 {
		return errors.New("Missing instance identifier: {resource_id},{instance_id}")
	}

	api := meta.(*api.API)
	log.Printf("[DEBUG] cloudamqp::resource::alarm::read instance id: %v", d.Get("instance_id"))
	data, err := api.ReadAlarm(d.Get("instance_id").(int), d.Id())
	log.Printf("[DEBUG] cloudamqp::resource::alarm::read data: %v", data)

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
	keys := []string{"type", "value_threshold", "time_threshold", "vhost_regex", "queue_regex", "notification_ids"}
	params := make(map[string]interface{})
	params["alarm_id"] = d.Id()
	log.Printf("[DEBUG] cloudamqp::resource::alarm::update params: %v", params)

	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	var notificationIDs []int
	if attr := d.Get("notification_ids").([]interface{}); len(attr) > 0 {
		for _, v := range attr {
			val, ok := v.(int)
			if ok {
				notificationIDs = append(notificationIDs, val)
			}
		}
	}
	params["notifications"] = notificationIDs
	delete(params, "notification_ids")
	log.Printf("[DEBUG] cloudamqp::resource::alarm::create updated params: %v", params)

	return api.UpdateAlarm(d.Get("instance_id").(int), params)
}

func resourceAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	params := make(map[string]interface{})
	params["alarm_id"] = d.Id()
	log.Printf("[DEBUG] cloudamqp::resource::alarm::delete params: %v", params)
	return api.DeleteAlarm(d.Get("instance_id").(int), params)
}

func validateAlarmType() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"cpu",
		"memory",
		"disk",
		"queue",
		"connection",
		"consumer",
		"netsplit",
		"ssh",
	}, true)
}

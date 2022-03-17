package cloudamqp

import (
	"errors"
	"fmt"
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
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable or disable an alarm",
			},
			"value_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "What value to trigger the alarm for",
			},
			"value_calculation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Disk value threshold calculation. Fixed or percentage of disk space remaining",
				ValidateFunc: validateValueCalculation(),
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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Message types (total, unacked, ready) of the queue to trigger the alarm",
				ValidateFunc: validateMessageType(),
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

func resourceAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := alarmAttributeKeys()
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}
	log.Printf("[DEBUG] cloudamqp::resource::alarm::create params: %v", params)

	data, err := api.CreateAlarm(d.Get("instance_id").(int), params)
	log.Printf("[DEBUG] cloudamqp::resource::alarm::create data: %v", data)

	if err != nil {
		return err
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
		log.Printf("[DEBUG] cloudamqp::resource::alarm::create id set: %v", d.Id())
	}

	return resourceAlarmRead(d, meta)
}

func resourceAlarmRead(d *schema.ResourceData, meta interface{}) error {
	if strings.Contains(d.Id(), ",") {
		log.Printf("[DEBUG] cloudamqp::resource::alarm::read id contains ,: %v", d.Id())
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		instanceID, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
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
		if validateAlarmSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	return nil
}

func resourceAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := alarmAttributeKeys()
	params := make(map[string]interface{})
	params["id"] = d.Id()
	log.Printf("[DEBUG] cloudamqp::resource::alarm::update params: %v", params)

	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	if err := api.UpdateAlarm(d.Get("instance_id").(int), params); err != nil {
		return err
	}

	return resourceAlarmRead(d, meta)
}

func resourceAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	params := make(map[string]interface{})
	params["id"] = d.Id()
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
		"notice",
		"server_unreachable",
	}, true)
}

func validateAlarmSchemaAttribute(key string) bool {
	switch key {
	case "type",
		"enabled",
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

func validateMessageType() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"total",
		"unacked",
		"ready",
	}, true)
}

func validateValueCalculation() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"fixed",
		"percentage",
	}, true)
}

func alarmAttributeKeys() []string {
	return []string {
		"type",
		"enabled",
		"value_threshold",
		"value_calculation",
		"time_threshold",
		"vhost_regex",
		"queue_regex",
		"message_type",
		"recipients",
	}
}

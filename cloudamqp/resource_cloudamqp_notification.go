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

func resourceNotification() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationCreate,
		Read:   resourceNotificationRead,
		Update: resourceNotificationUpdate,
		Delete: resourceNotificationDelete,
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
				Description:  "Type of the notification, valid options are: email, webhook, pagerduty, victorops, opsgenie, opsgenie-eu, slack",
				ValidateFunc: validateNotificationType(),
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Notification endpoint, where to send the notifcation",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional display name of the recipient",
			},
		},
	}
}

func resourceNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"type", "value", "name"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}
	log.Printf("[DEBUG] cloudamqp::resource::notification::create params: %v", params)

	data, err := api.CreateNotification(d.Get("instance_id").(int), params)
	log.Printf("[DEBUG] cloudamqp::resource::notification::create data: %v", data)

	if err != nil {
		return err
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
		log.Printf("[DEBUG] cloudamqp::resource::notification::create id set: %v", d.Id())
	}

	for k, v := range data {
		if validateRecipientAttribute(k) {
			d.Set(k, v)
		}
	}
	return nil
}

func resourceNotificationRead(d *schema.ResourceData, meta interface{}) error {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		instance_id, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instance_id)
	}
	if d.Get("instance_id").(int) == 0 {
		return errors.New("Missing instance identifier: {resource_id},{instance_id}")
	}

	log.Printf("[DEBUG] cloudamqp::resource::notification::read instance id: %v, id: %v", d.Get("instance_id"), d.Id())
	api := meta.(*api.API)
	data, err := api.ReadNotification(d.Get("instance_id").(int), d.Id())
	log.Printf("[DEBUG] cloudamqp::resource::notification::read data: %v", data)

	if err != nil {
		return err
	}

	for k, v := range data {
		if validateRecipientAttribute(k) {
			d.Set(k, v)
		}
	}
	return nil
}

func resourceNotificationUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"type", "value", "name"}
	params := make(map[string]interface{})
	params["id"] = d.Id()
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}
	log.Printf("[DEBUG] cloudamqp::resource::notification::update params: %v", params)
	return api.UpdateNotification(d.Get("instance_id").(int), params)
}

func resourceNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	params := make(map[string]interface{})
	params["id"] = d.Id()
	log.Printf("[DEBUG] cloudamqp::resource::notification::delete instance_id: %v, params: %v", d.Get("instance_id"), params)
	return api.DeleteNotification(d.Get("instance_id").(int), params)
}

func validateNotificationType() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"email",
		"webhook",
		"pagerduty",
		"victorops",
		"opsgenie",
		"opsgenie-eu",
		"slack",
	}, true)
}

func validateRecipientAttribute(key string) bool {
	switch key {
	case "type",
		"value",
		"name":
		return true
	}
	return false
}

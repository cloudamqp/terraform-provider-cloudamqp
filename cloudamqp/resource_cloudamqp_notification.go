package cloudamqp

import (
	"errors"
	"fmt"
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
				Description:  "Type of the notification, valid options are: email, webhook, pagerduty, victorops, opsgenie, opsgenie-eu, slack, teams",
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
			"routing_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Routing key used by VictorOps recipients",
			},
		},
	}
}

func resourceNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"type", "value", "name", "routing_key"}
	params := make(map[string]interface{})
	for _, k := range keys {
		v := d.Get(k)
		if k == "routing_key" {
			if v == nil {
				params["options"] = make(map[string]interface{})
			} else if d.Get("type") == "victorops" {
				rk := make(map[string]interface{})
				rk["rk"] = v
				params["options"] = rk
			}
		} else {
			if v != nil {
				params[k] = v
			}
		}
	}

	data, err := api.CreateNotification(d.Get("instance_id").(int), params)

	if err != nil {
		return err
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
	}

	return resourceNotificationRead(d, meta)
}

func resourceNotificationRead(d *schema.ResourceData, meta interface{}) error {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		instanceID, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
	}
	if d.Get("instance_id").(int) == 0 {
		return errors.New("Missing instance identifier: {resource_id},{instance_id}")
	}

	api := meta.(*api.API)
	data, err := api.ReadNotification(d.Get("instance_id").(int), d.Id())

	if err != nil {
		return err
	}

	for k, v := range data {
		if validateRecipientAttribute(k) {
			if k == "options" {
				if v != nil && v != "" && len(v.(map[string]interface{})) > 0 {
					k = "routing_key"
					v = v.(map[string]interface{})["rk"].(string)
				} else {
					continue
				}
			}
			if err = d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	return nil
}

func resourceNotificationUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"type", "value", "name", "routing_key"}
	params := make(map[string]interface{})
	params["id"] = d.Id()
	for _, k := range keys {
		v := d.Get(k)
		if k == "routing_key" {
			if v == nil {
				params["options"] = make(map[string]interface{})
			} else if d.Get("type") == "victorops" {
				rk := make(map[string]interface{})
				rk["rk"] = v
				params["options"] = rk
			}
		} else {
			if v != nil {
				params[k] = v
			}
		}
	}

	err := api.UpdateNotification(d.Get("instance_id").(int), params)
	if err != nil {
		return err
	}

	return resourceNotificationRead(d, meta)
}

func resourceNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	params := make(map[string]interface{})
	params["id"] = d.Id()
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
		"teams",
	}, true)
}

func validateRecipientAttribute(key string) bool {
	switch key {
	case "type",
		"value",
		"name",
		"options":
		return true
	}
	return false
}

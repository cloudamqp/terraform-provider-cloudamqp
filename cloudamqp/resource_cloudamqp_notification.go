package cloudamqp

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNotification() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationCreate,
		Read:   resourceNotificationRead,
		Update: resourceNotificationUpdate,
		Delete: resourceNotificationDelete,
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
				Type:     schema.TypeString,
				Required: true,
				Description: "Type of the notification, valid options are: email, opsgenie, opsgenie-eu," +
					"pagerduty, slack, signl4, teams, victorops, webhook",
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
			"options": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Optional key-value pair options parameters",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"responders": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Responder type, valid options are: team, user, escalation, schedule",
							ValidateFunc: validateOpsgenieRespondersType(),
						},
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
							Description:  "Responder ID",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Responder name",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Responder username",
						},
					},
				},
				Description: "Responders for OpsGenie alarms",
			},
		},
	}
}

func resourceNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api    = meta.(*api.API)
		keys   = []string{"type", "value", "name", "options", "responders"}
		params = make(map[string]interface{})
	)

	for _, k := range keys {
		if v := d.Get(k); v != nil && v != "" {
			switch k {
			case "responders":
				if len(v.(*schema.Set).List()) == 0 {
					continue
				}
				params["options"] = opsGenieRespondersParameter(v.(*schema.Set).List())
			default:
				params[k] = v
			}
		}
	}

	log.Printf("[DEBUG] resourceNotificationCreate params %v", params)
	data, err := api.CreateNotification(d.Get("instance_id").(int), params)
	if err != nil {
		return err
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
	}

	return nil
}

func resourceNotificationRead(d *schema.ResourceData, meta interface{}) error {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		instanceID, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
	}
	if d.Get("instance_id").(int) == 0 {
		return errors.New("missing instance identifier: {resource_id},{instance_id}")
	}

	api := meta.(*api.API)
	data, err := api.ReadNotification(d.Get("instance_id").(int), d.Id())
	if err != nil {
		return err
	}

	for k, v := range data {
		if !validateRecipientAttribute(k) {
			continue
		}
		if v == nil || v == "" {
			continue
		}

		switch k {
		case "options":
			for key, value := range data[k].(map[string]interface{}) {
				if key == "responders" {
					d.Set("responders", value.([]interface{}))
				} else {
					d.Set(k, v)
				}
			}
		default:
			d.Set(k, v)
		}
	}

	return nil
}

func resourceNotificationUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		api         = meta.(*api.API)
		keys        = []string{"type", "value", "name", "options", "responders"}
		params      = make(map[string]interface{})
		instanceID  = d.Get("instance_id").(int)
		recipientID = d.Id()
	)

	for _, k := range keys {
		if v := d.Get(k); v != nil && v != "" {
			switch k {
			case "responders":
				if len(v.(*schema.Set).List()) == 0 {
					continue
				}
				params["options"] = opsGenieRespondersParameter(v.(*schema.Set).List())
			default:
				params[k] = v
			}
		}
	}

	log.Printf("[DEBUG] resourceNotificationUpdate params %v", params)
	return api.UpdateNotification(instanceID, recipientID, params)
}

func resourceNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		api         = meta.(*api.API)
		instanceID  = d.Get("instance_id").(int)
		recipientID = d.Id()
	)

	return api.DeleteNotification(instanceID, recipientID)
}

func validateNotificationType() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"email",
		"opsgenie",
		"opsgenie-eu",
		"pagerduty",
		"signl4",
		"slack",
		"teams",
		"victorops",
		"webhook",
	}, true)
}

func validateRecipientAttribute(key string) bool {
	switch key {
	case "type",
		"value",
		"name",
		"options",
		"responders":
		return true
	}
	return false
}

func validateOpsgenieRespondersType() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"escalation",
		"schedule",
		"team",
		"user",
	}, true)
}

func opsGenieRespondersParameter(responders []interface{}) map[string]interface{} {
	responderParams := make(map[string]interface{})
	params := make([]map[string]interface{}, len(responders))
	for index, responder := range responders {
		param := make(map[string]interface{})
		for k, v := range responder.(map[string]interface{}) {
			if v != nil && v != "" {
				param[k] = v
			}
		}
		params[index] = param
	}
	responderParams["responders"] = params
	return responderParams
}

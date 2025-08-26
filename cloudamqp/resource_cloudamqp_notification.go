package cloudamqp

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationCreate,
		ReadContext:   resourceNotificationRead,
		UpdateContext: resourceNotificationUpdate,
		DeleteContext: resourceNotificationDelete,
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
				ValidateDiagFunc: validateNotificationType(),
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
							Type:             schema.TypeString,
							Required:         true,
							Description:      "Responder type, valid options are: team, user, escalation, schedule",
							ValidateDiagFunc: validateOpsgenieRespondersType(),
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

func resourceNotificationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api    = meta.(*api.API)
		keys   = []string{"type", "value", "name", "options", "responders"}
		params = make(map[string]any)
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
	data, err := api.CreateNotification(ctx, d.Get("instance_id").(int), params)
	if err != nil {
		return diag.FromErr(err)
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
	}

	return diag.Diagnostics{}
}

func resourceNotificationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if strings.Contains(d.Id(), ",") {
		tflog.Info(ctx, fmt.Sprintf("import of resource with identifier: %s", d.Id()))
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		instanceID, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
	}
	if d.Get("instance_id").(int) == 0 {
		return diag.Errorf("missing instance identifier: {resource_id},{instance_id}")
	}

	api := meta.(*api.API)
	data, err := api.ReadNotification(ctx, d.Get("instance_id").(int), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle resource drift and trigger re-creation if resource been deleted outside the provider
	if data == nil {
		d.SetId("")
		return nil
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
			for key, value := range data[k].(map[string]any) {
				if key == "responders" {
					d.Set("responders", value.([]any))
				} else {
					d.Set(k, v)
				}
			}
		default:
			d.Set(k, v)
		}
	}

	return diag.Diagnostics{}
}

func resourceNotificationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api         = meta.(*api.API)
		keys        = []string{"type", "value", "name", "options", "responders"}
		params      = make(map[string]any)
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

	if err := api.UpdateNotification(ctx, instanceID, recipientID, params); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func resourceNotificationDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api         = meta.(*api.API)
		instanceID  = d.Get("instance_id").(int)
		recipientID = d.Id()
	)

	if err := api.DeleteNotification(ctx, instanceID, recipientID); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func validateNotificationType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"email",
		"opsgenie",
		"opsgenie-eu",
		"pagerduty",
		"signl4",
		"slack",
		"teams",
		"victorops",
		"webhook",
	}, true))
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

func validateOpsgenieRespondersType() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"escalation",
		"schedule",
		"team",
		"user",
	}, true))
}

func opsGenieRespondersParameter(responders []any) map[string]any {
	responderParams := make(map[string]any)
	params := make([]map[string]any, len(responders))
	for index, responder := range responders {
		param := make(map[string]any)
		for k, v := range responder.(map[string]any) {
			if v != nil && v != "" {
				param[k] = v
			}
		}
		params[index] = param
	}
	responderParams["responders"] = params
	return responderParams
}

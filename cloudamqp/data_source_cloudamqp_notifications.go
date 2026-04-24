package cloudamqp

import (
	"context"
	"fmt"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	model "github.com/cloudamqp/terraform-provider-cloudamqp/api/models/monitoring"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNotifications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationsRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"recipients": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of notification recipients",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recipient_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Recipient identifier",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the notification",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Notification endpoint, where to send the notifcation",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name of the recipient",
						},
						"options": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Key-value pair options parameters",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"responders": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of responders (opsgenie/opsgenie-eu only)",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Responder type, valid options are: team, user, escalation, schedule",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Responder ID",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Responder name",
									},
									"username": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Responder username",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNotificationsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	instanceID := int64(d.Get("instance_id").(int))

	client := meta.(*api.API)
	data, err := client.ListNotifications(ctx, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d.notifications", instanceID))

	recipients := make([]map[string]any, len(data))
	for k, v := range data {
		recipients[k] = readNotification(v)
	}

	if err = d.Set("recipients", recipients); err != nil {
		return diag.Errorf("error setting recipients for resource %s: %s", d.Id(), err)
	}
	return diag.Diagnostics{}
}

func readNotification(data model.RecipientResponse) map[string]any {
	notification := map[string]any{
		"recipient_id": data.ID,
		"type":         data.Type,
		"value":        data.Value,
		"name":         data.Name,
		"options":      map[string]string{},
		"responders":   []map[string]any{},
	}

	switch data.Type {
	case "opsgenie", "opsgenie-eu":
		if data.Options != nil && data.Options.Responders != nil && len(*data.Options.Responders) > 0 {
			responders := make([]map[string]any, len(*data.Options.Responders))
			for i, r := range *data.Options.Responders {
				responder := map[string]any{
					"type":     r.Type,
					"id":       "",
					"name":     "",
					"username": "",
				}
				if r.ID != nil {
					responder["id"] = *r.ID
				}
				if r.Name != nil {
					responder["name"] = *r.Name
				}
				if r.Username != nil {
					responder["username"] = *r.Username
				}
				responders[i] = responder
			}
			notification["responders"] = responders
		}
	case "pagerduty", "victorops":
		if data.Options != nil {
			opts := map[string]string{}
			if data.Options.DedupKey != nil {
				opts["dedupkey"] = *data.Options.DedupKey
			}
			if data.Options.RK != nil {
				opts["rk"] = *data.Options.RK
			}
			notification["options"] = opts
		}
	}

	return notification
}

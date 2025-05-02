package cloudamqp

import (
	"context"
	"fmt"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
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
					},
				},
			},
		},
	}
}

func dataSourceNotificationsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		instanceID = d.Get("instance_id").(int)
		data       []map[string]any
		err        error
	)

	api := meta.(*api.API)
	data, err = api.ListNotifications(ctx, instanceID)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%v.notifications", instanceID))

	recipients := make([]map[string]any, len(data))
	for k, v := range data {
		recipients[k] = readNotification(v)
	}

	if err = d.Set("recipients", recipients); err != nil {
		return diag.Errorf("error setting recipients for resource %s: %s", d.Id(), err)
	}
	return diag.Diagnostics{}
}

func readNotification(data map[string]any) map[string]any {
	notification := make(map[string]any)
	for k, v := range data {
		if k == "id" {
			notification["recipient_id"] = int(v.(float64))
		} else if validateRecipientAttribute(k) {
			notification[k] = v
		}
	}
	return notification
}

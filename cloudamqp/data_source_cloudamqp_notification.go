package cloudamqp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNotification() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"recipient_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Recipient identifier",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the notification, valid options are: email, webhook, pagerduty, victorops, opsgenie, opsgenie-eu, slack",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
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
		},
	}
}

func dataSourceNotificationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		instanceID = d.Get("instance_id").(int)
		data       map[string]any
		err        error
	)

	// Multiple purpose read. To be used when using data source either by declaring recipient id or type.
	if d.Get("recipient_id") != 0 {
		data, err = dataSourceNotificationIDRead(ctx, instanceID, d.Get("recipient_id").(int), meta)
	} else if d.Get("name") != "" {
		data, err = dataSourceNotificationTypeRead(ctx, instanceID, d.Get("name").(string), meta)
	}

	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%v", data["id"]))
	for k, v := range data {
		if validateRecipientAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return diag.Diagnostics{}
}

func dataSourceNotificationIDRead(ctx context.Context, instanceID int, alarmID int,
	meta any) (map[string]any, error) {

	api := meta.(*api.API)
	id := strconv.Itoa(alarmID)
	recipient, err := api.ReadNotification(ctx, instanceID, id)
	return recipient, err
}

func dataSourceNotificationTypeRead(ctx context.Context, instanceID int, name string,
	meta any) (map[string]any, error) {

	api := meta.(*api.API)
	recipients, err := api.ListNotifications(ctx, instanceID)

	if err != nil {
		return nil, err
	}
	for _, recipient := range recipients {
		if recipient["name"] == name {
			return recipient, nil
		}
	}
	return nil, nil
}

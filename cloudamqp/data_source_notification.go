package cloudamqp

import (
	"fmt"
	"log"
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceNotification() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNotificationRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"recipient_id": {
				Type:        schema.TypeInt,
				Required:    true,
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
				Computed:    true,
				Description: "Optional display name of the recipient",
			},
		},
	}
}

func dataSourceNotificationRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	recipient_id := strconv.Itoa(d.Get("recipient_id").(int))
	data, err := api.ReadNotification(d.Get("instance_id").(int), recipient_id)
	log.Printf("[DEBUG] cloudamqp::data_source::notification::read data: %v", data)

	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%v", data["id"]))
	for k, v := range data {
		if validateRecipientAttribute(k) {
			d.Set(k, v)
		}
	}
	return nil
}

package cloudamqp

import (
	"fmt"
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
		},
	}
}

func dataSourceNotificationRead(d *schema.ResourceData, meta interface{}) error {
	var data map[string]interface{}
	var err error

	// Multiple purpose read. To be used when using data source either by declaring recipient id or type.
	if d.Get("recipient_id") != 0 {
		data, err = dataSourceNotificationIdRead(d.Get("instance_id").(int), d.Get("recipient_id").(int), meta)
	} else if d.Get("name") != "" {
		data, err = dataSourceNotificationTypeRead(d.Get("instance_id").(int), d.Get("name").(string), meta)
	}

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

func dataSourceNotificationIdRead(instance_id int, alarm_id int, meta interface{}) (map[string]interface{}, error) {
	api := meta.(*api.API)
	id := strconv.Itoa(alarm_id)
	recipient, err := api.ReadNotification(instance_id, id)
	return recipient, err
}

func dataSourceNotificationTypeRead(instance_id int, name string, meta interface{}) (map[string]interface{}, error) {
	api := meta.(*api.API)
	recipients, err := api.ReadNotifications(instance_id)

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

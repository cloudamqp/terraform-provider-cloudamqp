package cloudamqp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceWebhook() *schema.Resource {
	return &schema.Resource{
		Create: resourceWebhookCreate,
		Read:   resourceWebhookRead,
		Delete: resourceWebhookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Instance identifier",
			},
			"vhost": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the virtual host",
			},
			"queue": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The queue that should be forwarded, must be a durable queue!",
			},
			"webhook_uri": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A POST request will be made for each message in the queue to this endpoint",
			},
			"retry_interval": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "How often push of a message will retry if the previous call fails. In seconds",
			},
			"concurrency": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "How many times the request will be made if previous call fails",
			},
		},
	}
}

func resourceWebhookCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"vhost", "queue", "webhook_uri", "retry_interval", "concurrency"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	data, err := api.CreateWebhook(d.Get("instance_id").(int), params)
	if err != nil {
		return err
	}

	d.SetId(data["id"].(string))
	return resourceWebhookRead(d, meta)
}

func resourceWebhookRead(d *schema.ResourceData, meta interface{}) error {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		d.Set("name", s[0])
		instanceID, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
	}
	if d.Get("instance_id").(int) == 0 {
		return errors.New("Missing instance identifier: {resource_id},{instance_id}")
	}

	api := meta.(*api.API)
	data, err := api.ReadWebhook(d.Get("instance_id").(int), d.Id())

	if err != nil {
		return err
	}

	for k, v := range data {
		if validateWebhookSchemaAttribute(k) {
			err = d.Set(k, v)

			if err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	return nil
}

func resourceWebhookDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	return api.DeleteWebhook(d.Get("instance_id").(int), d.Id())
}

func validateWebhookSchemaAttribute(key string) bool {
	switch key {
	case "vhost",
		"queue",
		"webhook_uri",
		"retry_interval",
		"concurrency":
		return true
	}
	return false
}

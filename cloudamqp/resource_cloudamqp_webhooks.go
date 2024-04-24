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
		Update: resourceWebhookUpdate,
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
				Description: "The name of the virtual host",
			},
			"queue": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The queue that should be forwarded, must be a durable queue!",
			},
			"webhook_uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A POST request will be made for each message in the queue to this endpoint",
			},
			"concurrency": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "How many times the request will be made if previous call fails",
			},
			"sleep": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Configurable sleep time in seconds between retries for webhook",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1800,
				Description: "Configurable timeout time in seconds for webhook",
			},
		},
	}
}

func resourceWebhookCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		keys       = []string{"vhost", "queue", "webhook_uri", "concurrency"}
		params     = make(map[string]interface{})
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	data, err := api.CreateWebhook(instanceID, params, sleep, timeout)
	if err != nil {
		return err
	}

	d.SetId(data["id"].(string))
	return nil
}

func resourceWebhookRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		d.Set("name", s[0])
		instanceID, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
	}
	if d.Get("instance_id").(int) == 0 {
		return errors.New("missing instance identifier: {resource_id},{instance_id}")
	}

	data, err := api.ReadWebhook(instanceID, d.Id(), sleep, timeout)
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

func resourceWebhookUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		keys       = []string{"vhost", "queue", "webhook_uri", "concurrency"}
		params     = make(map[string]interface{})
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	return api.UpdateWebhook(instanceID, d.Id(), params, sleep, timeout)
}

func resourceWebhookDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	return api.DeleteWebhook(instanceID, d.Id(), sleep, timeout)
}

func validateWebhookSchemaAttribute(key string) bool {
	switch key {
	case "vhost",
		"queue",
		"webhook_uri",
		"concurrency":
		return true
	}
	return false
}

package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceWebhook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebhookCreate,
		ReadContext:   resourceWebhookRead,
		UpdateContext: resourceWebhookUpdate,
		DeleteContext: resourceWebhookDelete,
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

func resourceWebhookCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		keys       = []string{"vhost", "queue", "webhook_uri", "concurrency"}
		params     = make(map[string]any)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	data, err := api.CreateWebhook(ctx, instanceID, params, sleep, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(data["id"].(string))
	return diag.Diagnostics{}
}

func resourceWebhookRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	if strings.Contains(d.Id(), ",") {
		tflog.Info(ctx, fmt.Sprintf("import of resource with identifier: %s", d.Id()))
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		instanceID, _ = strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
		// Set default values for optional arguments
		d.Set("sleep", 10)
		d.Set("timeout", 1800)
	}
	if d.Get("instance_id").(int) == 0 {
		return diag.Errorf("missing instance identifier: {resource_id},{instance_id}")
	}

	data, err := api.ReadWebhook(ctx, instanceID, d.Id(), sleep, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range data {
		if validateWebhookSchemaAttribute(k) {
			err = d.Set(k, v)

			if err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	return nil
}

func resourceWebhookUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		keys       = []string{"vhost", "queue", "webhook_uri", "concurrency"}
		params     = make(map[string]any)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	params["webhook_id"] = d.Id()
	if err := api.UpdateWebhook(ctx, instanceID, d.Id(), params, sleep, timeout); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func resourceWebhookDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	if err := api.DeleteWebhook(ctx, instanceID, d.Id(), sleep, timeout); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
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

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

func resourcePluginCommunity() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginCommunityCreate,
		ReadContext:   resourcePluginCommunityRead,
		UpdateContext: resourcePluginCommunityUpdate,
		DeleteContext: resourcePluginCommunityDelete,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the plugin",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "If the plugin is enabled",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the plugin",
			},
			"require": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Required version of RabbitMQ",
			},
			"sleep": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Configurable sleep time in seconds between retries for plugins",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1800,
				Description: "Configurable timeout time in seconds for plugins",
			},
		},
	}
}

func resourcePluginCommunityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	data, err := api.ReadPluginCommunity(ctx, instanceID, name, sleep, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = api.InstallPluginCommunity(ctx, instanceID, name, sleep, timeout)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(name)

	for k, v := range data {
		if validateCommunityPluginSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return nil
}

func resourcePluginCommunityRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	// Support for importing resource
	if strings.Contains(d.Id(), ",") {
		tflog.Info(ctx, fmt.Sprintf("import of resource with identifiers: %s", d.Id()))
		s := strings.Split(d.Id(), ",")
		name = s[0]
		d.SetId(name)
		d.Set("name", name)
		instanceID, _ = strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
		// Set default values for optional arguments
		d.Set("sleep", 10)
		d.Set("timeout", 1800)
	}
	if instanceID == 0 {
		return diag.Errorf("missing instance identifier: {resource_id},{instance_id}")
	}

	data, err := api.ReadPluginCommunity(ctx, instanceID, name, sleep, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range data {
		if validateCommunityPluginSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	return diag.Diagnostics{}
}

func resourcePluginCommunityUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		enabled    = d.Get("enabled").(bool)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	_, err := api.UpdatePluginCommunity(ctx, instanceID, name, enabled, sleep, timeout)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourcePluginCommunityRead(ctx, d, meta)
}

func resourcePluginCommunityDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	if enableFasterInstanceDestroy {
		tflog.Debug(ctx, "delete will skip calling backend.")
		return nil
	}

	if _, err := api.UninstallPluginCommunity(ctx, instanceID, name, sleep, timeout); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func validateCommunityPluginSchemaAttribute(key string) bool {
	switch key {
	case "name",
		"require",
		"description":
		return true
	}
	return false
}

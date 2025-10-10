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
)

func resourcePlugin() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePluginCreate,
		ReadContext:   resourcePluginRead,
		UpdateContext: resourcePluginUpdate,
		DeleteContext: resourcePluginDelete,
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
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the plugin",
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

type Plugin struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func resourcePluginCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	_, err := api.EnablePlugin(ctx, instanceID, name, sleep, timeout)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourcePluginRead(ctx, d, meta)
}

func resourcePluginRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

	log.Printf("[DEBUG] import plugin instanceID: %v, name: %v, sleep: %v, timeout: %v",
		instanceID, name, sleep, timeout)
	data, err := api.ReadPlugin(ctx, instanceID, name, sleep, timeout)
	if err != nil {
		// If instance not found (404), return nil to indicate resource not found
		// This allows Terraform to recreate the resource when the instance is recreated
		if strings.Contains(err.Error(), "instance not found") || strings.Contains(err.Error(), "status=404") {
			tflog.Info(ctx, fmt.Sprintf("instance not found, plugin resource will be recreated: %s", name))
			return nil
		}
		return diag.FromErr(err)
	}

	// If no data returned (instance not found), return nil to indicate resource not found
	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("plugin not found, resource will be recreated: %s", name))
		return nil
	}

	for k, v := range data {
		if validatePluginSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	return diag.Diagnostics{}
}

func resourcePluginUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		enabled    = d.Get("enabled").(bool)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	_, err := api.UpdatePlugin(ctx, instanceID, name, enabled, sleep, timeout)
	if err != nil {
		return diag.Errorf("[Failed to update pluign: %v", err)
	}
	return resourcePluginRead(ctx, d, meta)
}

func resourcePluginDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	if enableFasterInstanceDestroy {
		log.Printf("[DEBUG] cloudamqp::resource::plugin::delete skip calling backend.")
		return diag.Diagnostics{}
	}

	if err := api.DeletePlugin(ctx, instanceID, name, sleep, timeout); err != nil {
		// If instance not found (404), consider deletion successful
		if strings.Contains(err.Error(), "instance not found") || strings.Contains(err.Error(), "status=404") {
			tflog.Info(ctx, fmt.Sprintf("instance not found during plugin deletion, considering successful: %s", name))
			return nil
		}
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func validatePluginSchemaAttribute(key string) bool {
	switch key {
	case "name",
		"enabled",
		"description",
		"version":
		return true
	}
	return false
}

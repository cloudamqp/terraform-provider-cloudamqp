package cloudamqp

import (
	"context"
	"fmt"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePlugins() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluginsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"plugins": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
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

func dataSourcePluginsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	data, err := api.ListPlugins(ctx, instanceID, sleep, timeout)
	d.SetId(fmt.Sprintf("%v.plugins", instanceID))
	if err != nil {
		return diag.FromErr(err)
	}

	plugins := make([]map[string]any, len(data))
	for k, v := range data {
		plugins[k] = readPlugin(v)
	}

	if err = d.Set("plugins", plugins); err != nil {
		return diag.Errorf("error setting plugins for resource %s: %s", d.Id(), err)
	}
	return nil
}

func readPlugin(data map[string]any) map[string]any {
	plugin := make(map[string]any)
	for k, v := range data {
		if validatePluginsSchemaAttribute(k) {
			plugin[k] = v
		}
	}
	return plugin
}

func validatePluginsSchemaAttribute(key string) bool {
	switch key {
	case "name",
		"version",
		"description",
		"enabled":
		return true
	}
	return false
}

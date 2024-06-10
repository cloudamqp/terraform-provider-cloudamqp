package cloudamqp

import (
	"fmt"
	"log"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePlugins() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePluginsRead,

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
				},
			},
		},
	}
}

func dataSourcePluginsRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	log.Printf("[DEBUG] cloudamqp::data_source::plugins::read instance id: %v", instanceID)
	data, err := api.ListPlugins(instanceID, sleep, timeout)
	d.SetId(fmt.Sprintf("%v.plugins", instanceID))
	if err != nil {
		return err
	}

	plugins := make([]map[string]interface{}, len(data))
	for k, v := range data {
		plugins[k] = readPlugin(v)
	}

	if err = d.Set("plugins", plugins); err != nil {
		return fmt.Errorf("error setting plugins for resource %s: %s", d.Id(), err)
	}
	return nil
}

func readPlugin(data map[string]interface{}) map[string]interface{} {
	plugin := make(map[string]interface{})
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

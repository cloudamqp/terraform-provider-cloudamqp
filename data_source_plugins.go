package main

import (
	"fmt"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform/helper/schema"
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
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePluginsRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data, err := api.ReadPlugins(d.Get("instance_id").(int))
	d.SetId(fmt.Sprintf("%v.plugins", d.Get("instance_id").(int)))
	if err != nil {
		return err
	}
	d.Set("plugins", data)
	return nil
}

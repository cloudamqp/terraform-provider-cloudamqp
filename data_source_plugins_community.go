package main

import (
	"fmt"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourcePluginsCommunity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePluginsCommunityRead,

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
						"require": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePluginsCommunityRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data, err := api.ReadPluginsCommunity(d.Get("instance_id").(int))
	d.SetId(fmt.Sprintf("%v.plugins_community", d.Get("instance_id").(int)))
	if err != nil {
		return err
	}
	d.Set("plugins", data)
	return nil
}

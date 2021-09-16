package cloudamqp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePluginCommunity() *schema.Resource {
	return &schema.Resource{
		Create: resourcePluginCommunityCreate,
		Read:   resourcePluginCommunityRead,
		Update: resourcePluginCommunityUpdate,
		Delete: resourcePluginCommunityDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
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
		},
	}
}

func resourcePluginCommunityCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	_, err := api.EnablePluginCommunity(d.Get("instance_id").(int), d.Get("name").(string))
	if err != nil {
		return err
	}
	d.SetId(d.Get("name").(string))
	return resourcePluginCommunityRead(d, meta)
}

func resourcePluginCommunityRead(d *schema.ResourceData, meta interface{}) error {
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
	data, err := api.ReadPluginCommunity(d.Get("instance_id").(int), d.Get("name").(string))
	if err != nil {
		return err
	}

	for k, v := range data {
		if validateCommunityPluginSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	return nil
}

func resourcePluginCommunityUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"name", "enabled"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}
	_, err := api.UpdatePluginCommunity(d.Get("instance_id").(int), params)
	if err != nil {
		return err
	}
	return resourcePluginCommunityRead(d, meta)
}

func resourcePluginCommunityDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	_, err := api.DisablePluginCommunity(d.Get("instance_id").(int), d.Get("name").(string))
	return err
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

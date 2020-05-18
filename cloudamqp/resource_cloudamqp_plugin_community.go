package cloudamqp

import (
	"errors"
	"fmt"
	"log"
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
	log.Printf("[DEBUG] cloudamqp::resource::plugin_community::create instance id: %v, name: %v", d.Get("instance_id"), d.Get("name"))
	data, err := api.EnablePluginCommunity(d.Get("instance_id").(int), d.Get("name").(string))
	log.Printf("[DEBUG] cloudamqp::resource::plugin_community::create data: %v", data)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s", d.Get("name").(string)))
	log.Printf("[DEBUG] cloudamqp::resource::plugin::create id set: %v", d.Id())
	for k, v := range data {
		d.Set(k, v)
	}
	return nil
}

func resourcePluginCommunityRead(d *schema.ResourceData, meta interface{}) error {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		d.Set("name", s[0])
		instance_id, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instance_id)
	}
	if d.Get("instance_id").(int) == 0 {
		return errors.New("Missing instance identifier: {resource_id},{instance_id}")
	}

	api := meta.(*api.API)
	log.Printf("[DEBUG] cloudamqp::resource::plugin_community::read instance id: %v, name: %v", d.Get("instance_id"), d.Get("name"))
	data, err := api.ReadPluginCommunity(d.Get("instance_id").(int), d.Get("name").(string))
	log.Printf("[DEBUG] cloudamqp::resource::plugin::read data: %v", data)
	if err != nil {
		return err
	}

	for k, v := range data {
		d.Set(k, v)
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
	log.Printf("[DEBUG] cloudamqp::resource::plugin::update instance id: %v, params: %v", d.Get("instance_id"), params)
	_, err := api.UpdatePluginCommunity(d.Get("instance_id").(int), params)
	return err
}

func resourcePluginCommunityDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	log.Printf("[DEBUG] cloudamqp::resource::plugin::delete instance id: %v, name: %v", d.Get("instance_id"), d.Get("name"))
	_, err := api.DisablePluginCommunity(d.Get("instance_id").(int), d.Get("name").(string))
	return err
}

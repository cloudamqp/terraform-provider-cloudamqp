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

func resourcePlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourcePluginCreate,
		Read:   resourcePluginRead,
		Update: resourcePluginUpdate,
		Delete: resourcePluginDelete,
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

func resourcePluginCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	log.Printf("[DEBUG] cloudamqp::resource::plugin::create instance id: %v, name: %v", d.Get("instance_id"), d.Get("name"))
	data, err := api.EnablePlugin(d.Get("instance_id").(int), d.Get("name").(string))
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s", d.Get("name").(string)))
	log.Printf("[DEBUG] cloudamqp::resource::plugin::create id set: %v", d.Id())
	for k, v := range data {
		if k == "id" || k == "version" || k == "description" {
			continue
		}
		d.Set(k, v)
	}
	return nil
}

func resourcePluginRead(d *schema.ResourceData, meta interface{}) error {
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
	log.Printf("[DEBUG] cloudamqp::resource::plugin::read instance id: %v, name: %v", d.Get("instance_id"), d.Get("name"))
	data, err := api.ReadPlugin(d.Get("instance_id").(int), d.Get("name").(string))
	log.Printf("[DEBUG] cloudamqp::resource::plugin::read data: %v", data)
	if err != nil {
		return err
	}

	for k, v := range data {
		if k == "id" || k == "version" || k == "description" {
			continue
		}
		d.Set(k, v)
	}

	return nil
}

func resourcePluginUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"name", "enabled"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}
	log.Printf("[DEBUG] cloudamqp::resource::plugin::update instance id: %v, params: %v", d.Get("instance_id"), params)
	_, err := api.UpdatePlugin(d.Get("instance_id").(int), params)
	if err != nil {
		log.Printf("[ERROR] cloudamqp::resource::plugin::update Failed to update pluign: %v", err)
	}
	return resourcePluginRead(d, meta)
}

func resourcePluginDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	log.Printf("[DEBUG] cloudamqp::resource::plugin::delete instance id: %v, name: %v", d.Get("instance_id"), d.Get("name"))
	err := api.DeletePlugin(d.Get("instance_id").(int), d.Get("name").(string))
	return err
}

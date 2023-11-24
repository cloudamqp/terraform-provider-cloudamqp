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

func resourcePluginCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	_, err := api.EnablePlugin(instanceID, name, sleep, timeout)
	if err != nil {
		return err
	}
	d.SetId(name)
	return resourcePluginRead(d, meta)
}

func resourcePluginRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	// Support for importing resources
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		d.Set("name", s[0])
		instanceID, _ = strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
	}
	if instanceID == 0 {
		return errors.New("missing instance identifier: {resource_id},{instance_id}")
	}

	data, err := api.ReadPlugin(instanceID, name, sleep, timeout)
	if err != nil {
		return err
	}

	for k, v := range data {
		if validatePluginSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	return nil
}

func resourcePluginUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		keys       = []string{"name", "enabled"}
		params     map[string]interface{}
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	_, err := api.UpdatePlugin(instanceID, params, sleep, timeout)
	if err != nil {
		return fmt.Errorf("[Failed to update pluign: %v", err)
	}
	return resourcePluginRead(d, meta)
}

func resourcePluginDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	if enableFasterInstanceDestroy {
		log.Printf("[DEBUG] cloudamqp::resource::plugin::delete skip calling backend.")
		return nil
	}

	return api.DeletePlugin(instanceID, name, sleep, timeout)
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

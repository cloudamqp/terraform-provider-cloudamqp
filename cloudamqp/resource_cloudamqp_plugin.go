package cloudamqp

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

type Plugin struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func resourcePluginCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	log.Printf("[DEBUG] create plugin instanceID: %v, name: %v, sleep: %v, timeout: %v",
		instanceID, name, sleep, timeout)
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

	// Support for importing resource
	if strings.Contains(d.Id(), ",") {
		log.Printf("[DEBUG] import plugin instanceID: %v, id: %v", instanceID, d.Id())
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
		return errors.New("missing instance identifier: {resource_id},{instance_id}")
	}

	log.Printf("[DEBUG] import plugin instanceID: %v, name: %v, sleep: %v, timeout: %v",
		instanceID, name, sleep, timeout)
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
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		enabled    = d.Get("enabled").(bool)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	_, err := api.UpdatePlugin(instanceID, name, enabled, sleep, timeout)
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

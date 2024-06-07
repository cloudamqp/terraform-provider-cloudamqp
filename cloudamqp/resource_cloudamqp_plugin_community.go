package cloudamqp

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePluginCommunity() *schema.Resource {
	return &schema.Resource{
		Create: resourcePluginCommunityCreate,
		Read:   resourcePluginCommunityRead,
		Update: resourcePluginCommunityUpdate,
		Delete: resourcePluginCommunityDelete,
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
			"require": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Required version of RabbitMQ",
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

func resourcePluginCommunityCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	data, err := api.ReadPluginCommunity(instanceID, name, sleep, timeout)
	if err != nil {
		return err
	}

	_, err = api.InstallPluginCommunity(instanceID, name, sleep, timeout)
	if err != nil {
		return err
	}
	d.SetId(name)

	for k, v := range data {
		if validateCommunityPluginSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return nil
}

func resourcePluginCommunityRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	// Support for importing resource
	if strings.Contains(d.Id(), ",") {
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

	data, err := api.ReadPluginCommunity(instanceID, name, sleep, timeout)
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
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		enabled    = d.Get("enabled").(bool)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	_, err := api.UpdatePluginCommunity(instanceID, name, enabled, sleep, timeout)
	if err != nil {
		return err
	}
	return resourcePluginCommunityRead(d, meta)
}

func resourcePluginCommunityDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		name       = d.Get("name").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	if enableFasterInstanceDestroy {
		log.Printf("[DEBUG] cloudamqp::resource::plugin-community::delete skip calling backend.")
		return nil
	}

	_, err := api.UninstallPluginCommunity(instanceID, name, sleep, timeout)
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

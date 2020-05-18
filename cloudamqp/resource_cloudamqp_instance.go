package cloudamqp

import (
	"fmt"
	"log"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreate,
		Read:   resourceRead,
		Update: resourceUpdate,
		Delete: resourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the instance",
			},
			"plan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the plan, valid options are: lemur, tiger, bunny, rabbit, panda, ape, hippo, lion",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the region you want to create your instance in",
			},
			"vpc_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Dedicated VPC subnet, shouldn't overlap with your current VPC's subnet",
			},
			"nodes": {
				Type:        schema.TypeInt,
				Default:     1,
				Optional:    true,
				Description: "Number of nodes in cluster (plan must support it)",
			},
			"rmq_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "RabbitMQ version",
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "URL of the CloudAMQP instance",
			},
			"apikey": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "API key for the CloudAMQP instance",
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Tag the instances with optional tags",
			},
			"host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Host name for the CloudAMQP instance",
			},
			"vhost": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The virtual host",
			},
			"ready": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag describing if the resource is ready",
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("plan", func(old, new, meta interface{}) bool {
				// Recreate instance if changing plan type (from dedicated to shared or vice versa)
				return !(getPlanType(old.(string)) == getPlanType(new.(string)))
			}),
		),
	}
}

func resourceCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"name", "plan", "region", "nodes", "tags", "rmq_version", "vpc_subnet"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		} else if k == "rmq_version" {
			version, _ := api.DefaultRmqVersion()
			params[k] = version["default_rmq_version"]
		}
		if k == "vpc_subnet" {
			if d.Get(k) == "" {
				delete(params, "vpc_subnet")
			}
		}
	}
	log.Printf("[DEBUG] cloudamqp::resource::instance::create params: %v", params)

	data, err := api.CreateInstance(params)
	log.Printf("[DEBUG] cloudamqp::resource::instance::create data: %v", data)

	if err != nil {
		return err
	}
	d.SetId(data["id"].(string))
	log.Printf("[DEBUG] cloudamqp::resource::instance::create id set: %v", d.Id())
	for k, v := range data {
		if validateInstanceSchemaAttribute(k) {
			if k == "vpc" {
				d.Set("vpc_subnet", v.(map[string]interface{})["subnet"])
			} else {
				d.Set(k, v)
			}
		}
	}

	data = api.UrlInformation(data["url"].(string))
	if err == nil {
		d.Set("host", data["host"])
		d.Set("vhost", data["vhost"])
	}
	return nil
}

func resourceRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	log.Printf("[DEBUG] cloudamqp::resource::instance::read id: %v", d.Id())
	data, err := api.ReadInstance(d.Id())
	log.Printf("[DEBUG] cloudamqp::resource::instance::read data: %v", data)

	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%v", data["id"]))
	for k, v := range data {
		if validateInstanceSchemaAttribute(k) {
			if k == "vpc" {
				d.Set("vpc_subnet", v.(map[string]interface{})["subnet"])
			} else {
				d.Set(k, v)
			}
		}
	}

	data = api.UrlInformation(data["url"].(string))
	if err == nil {
		d.Set("host", data["host"])
		d.Set("vhost", data["vhost"])
	}
	return nil
}

func resourceUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"name", "plan", "nodes", "tags"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = d.Get(k)
		}
	}
	log.Printf("[DEBUG] cloudamqp::resource::instance::update params: %v", params)
	return api.UpdateInstance(d.Id(), params)
}

func resourceDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	log.Printf("[DEBUG] cloudamqp::resource::instance::delete id: %v", d.Id())
	return api.DeleteInstance(d.Id())
}

func validateInstanceSchemaAttribute(key string) bool {
	switch key {
	case "name",
		"plan",
		"region",
		"vpc",
		"nodes",
		"rmq_version",
		"url",
		"apikey",
		"tags":
		return true
	}
	return false
}

func getPlanType(plan string) string {
	switch plan {
	case "lemur", "tiger":
		return "shared"
	case "bunny", "rabbit", "panda", "ape", "hippo", "lion":
		return "dedicated"
	default:
		return "unknown" // This shouldn't happen. However we shouldn't break if a new instance type gets implemented
	}
}

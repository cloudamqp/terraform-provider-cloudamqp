package cloudamqp

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the plan, see documentation for valid plans",
				ValidateFunc: validatePlanName(),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the region you want to create your instance in",
			},
			"vpc_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the VPC to create your instance in",
			},
			"vpc_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Dedicated VPC subnet, shouldn't overlap with your current VPC's subnet",
			},
			"nodes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Number of nodes in cluster (plan must support it)",
			},
			"rmq_version": {
				Type:        schema.TypeString,
				Computed:    true,
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
				Description: "External hostname for the CloudAMQP instance",
			},
			"host_internal": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal hostname for the CloudAMQP instance",
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
			"dedicated": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the instance hosted on a dedicated server",
			},
			"no_default_alarms": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Set to true to not create default alarms",
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("plan", func(old, new, meta interface{}) bool {
				// Recreate instance if changing plan type (from dedicated to shared or vice versa)
				oldPlanType, _ := getPlanType(old.(string))
				newPlanType, _ := getPlanType(new.(string))
				return !(oldPlanType == newPlanType)
			}),
		),
	}
}

func resourceCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"name", "plan", "region", "nodes", "tags", "rmq_version", "vpc_id", "vpc_subnet", "no_default_alarms"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil && v != "" {
			params[k] = v
		} else if k == "rmq_version" {
			version, _ := api.DefaultRmqVersion()
			params[k] = version["default_rmq_version"]
		} else if k == "no_default_alarms" {
			params[k] = false
		}

		if k == "nodes" {
			plan := d.Get("plan").(string)
			if is2020Plan(plan) {
				nodes := numberOfNodes(plan)
				params[k] = nodes
			}
		}

		if k == "vpc_id" {
			if d.Get(k) == "" {
				delete(params, "vpc_id")
			}
		}

		if k == "vpc_subnet" {
			if d.Get(k) == "" {
				delete(params, "vpc_subnet")
			}
		}
	}

	data, err := api.CreateInstance(params)
	if err != nil {
		return err
	}

	d.SetId(data["id"].(string))
	return resourceRead(d, meta)
}

func resourceRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data, err := api.ReadInstance(d.Id())

	if err != nil {
		return err
	}

	for k, v := range data {
		if validateInstanceSchemaAttribute(k) {
			if k == "vpc" {
				err = d.Set("vpc_id", v.(map[string]interface{})["id"])
				err = d.Set("vpc_subnet", v.(map[string]interface{})["subnet"])
			} else if k == "nodes" {
				plan := d.Get("plan").(string)
				if is2020Plan(plan) {
					nodes := numberOfNodes(plan)
					err = d.Set(k, nodes)
				} else {
					err = d.Set(k, v)
				}
			} else {
				err = d.Set(k, v)
			}

			if err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	if err = d.Set("host", data["hostname_external"].(string)); err != nil {
		return fmt.Errorf("error setting host for resource %s: %s", d.Id(), err)
	}

	if err = d.Set("host_internal", data["hostname_internal"].(string)); err != nil {
		return fmt.Errorf("error setting host for resource %s: %s", d.Id(), err)
	}

	planType, _ := getPlanType(d.Get("plan").(string))
	dedicated := planType == "dedicated"
	if err = d.Set("dedicated", dedicated); err != nil {
		return fmt.Errorf("error setting dedicated for resource %s: %s", d.Id(), err)
	}

	data = api.UrlInformation(data["url"].(string))
	for k, v := range data {
		if validateInstanceSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
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
		if k == "nodes" {
			plan := d.Get("plan").(string)
			if is2020Plan(plan) {
				nodes := numberOfNodes(plan)
				params[k] = nodes
			}
		}
	}

	if err := api.UpdateInstance(d.Id(), params); err != nil {
		return err
	}

	return resourceRead(d, meta)
}

func resourceDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
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
		"tags",
		"vhost",
		"no_default_alarms",
		"ready":
		return true
	}
	return false
}

func getPlanType(plan string) (string, error) {
	switch plan {
	case "lemur", "tiger":
		return "shared", nil
	// Legacy plans
	case "bunny", "rabbit", "panda", "ape", "hippo", "lion",
		// 2020 plans
		"squirrel-1",
		"bunny-1", "bunny-3",
		"rabbit-1", "rabbit-3", "rabbit-5",
		"panda-1", "panda-3", "panda-5",
		"ape-1", "ape-3", "ape-5",
		"hippo-1", "hippo-3", "hippo-5",
		"lion-1", "lion-3", "lion-5",
		"rhino-1":
		return "dedicated", nil
	}
	return "", fmt.Errorf("couldn't find a matching plan type for: %s", plan)
}

func validatePlanName() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"lemur", "tiger",
		"bunny", "rabbit", "panda", "ape", "hippo", "lion",
		"squirrel-1", "bunny-1", "bunny-3",
		"rabbit-1", "rabbit-3", "rabbit-5",
		"panda-1", "panda-3", "panda-5",
		"ape-1", "ape-3", "ape-5",
		"hippo-1", "hippo-3", "hippo-5",
		"lion-1", "lion-3", "lion-5",
		"rhino-1",
	}, true)
}

func is2020Plan(plan string) bool {
	switch plan {
	case
		"squirrel-1",
		"bunny-1", "bunny-3",
		"rabbit-1", "rabbit-3", "rabbit-5",
		"panda-1", "panda-3", "panda-5",
		"ape-1", "ape-3", "ape-5",
		"hippo-1", "hippo-3", "hippo-5",
		"lion-1", "lion-3", "lion-5",
		"rhino-1":
		return true
	}
	return false
}

func numberOfNodes(plan string) int {
	r := regexp.MustCompile("[135]")
	match := r.FindString(plan)
	nodes, _ := strconv.Atoi(match)
	return nodes
}

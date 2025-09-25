package cloudamqp

import (
	"context"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description: "Name of the plan, see documentation for valid plans",
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
			"keep_associated_vpc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Keep associated VPC when deleting instance",
			},
			"backend": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Software backend used, determined by subscription plan",
			},
			"preferred_az": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Preferred availability zone for the instance(s), one entry per node in the cluster",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"copy_settings": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscription_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Instance identifier of the CloudAMQP instance to copy settings from",
						},
						"settings": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validateCopySettings(),
							},
							Description: "Settings to be copied. [alarms, config, definitions, firewall, logs, metrics, plugins]",
						},
					},
				},
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("plan", func(ctx context.Context, old, new, meta any) bool {
				// Recreate instance if changing plan type (from dedicated to shared or vice versa)
				oldPlanType := isSharedPlan(old.(string))
				newPlanType := isSharedPlan(new.(string))
				return !(oldPlanType == newPlanType)
			}),
			customdiff.ValidateChange("plan", func(ctx context.Context, old, new, meta any) error {
				if old == new {
					return nil
				}
				api := meta.(*api.API)
				return api.ValidatePlan(new.(string))
			}),
			customdiff.ValidateChange("region", func(ctx context.Context, old, new, meta any) error {
				if old == new {
					return nil
				}
				api := meta.(*api.API)
				return api.ValidateRegion(new.(string))
			}),
		),
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	keys := instanceCreateAttributeKeys()
	params := make(map[string]any)
	for _, k := range keys {
		if v := d.Get(k); v != nil && v != "" {
			params[k] = v
		} else if k == "no_default_alarms" {
			params[k] = false
		}

		// Remove keys from params
		switch k {
		case "nodes":
			plan := d.Get("plan").(string)
			if isSharedPlan(plan) || !isLegacyPlan(plan) {
				delete(params, k)
			}
		case "vpc_id":
			if d.Get(k).(int) == 0 {
				delete(params, k)
			}
		case "vpc_subnet":
			if d.Get(k) == "" {
				delete(params, k)
			}
		case "copy_settings":
			if d.Get(k).(*schema.Set).Len() == 0 {
				delete(params, k)
			} else {
				for _, v := range d.Get(k).(*schema.Set).List() {
					params[k] = v.(map[string]any)
				}
			}
		}
	}

	data, err := api.CreateInstance(ctx, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(data["id"].(string))
	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	data, err := api.ReadInstance(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Handle resource drift and trigger re-creation if resource been deleted outside the provider
	if data == nil {
		d.SetId("")
		return nil
	}

	for k, v := range data {
		if validateInstanceSchemaAttribute(k) {
			if k == "vpc" {
				err = d.Set("vpc_id", v.(map[string]any)["id"])
			} else {
				err = d.Set(k, v)
			}

			if err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	if v, ok := d.Get("nodes").(int); ok && v > 0 {
		d.Set("dedicated", true)
	} else {
		d.Set("dedicated", false)
	}

	if err = d.Set("host", data["hostname_external"].(string)); err != nil {
		return diag.Errorf("error setting host for resource %s: %s", d.Id(), err)
	}

	if err = d.Set("host_internal", data["hostname_internal"].(string)); err != nil {
		return diag.Errorf("error setting host for resource %s: %s", d.Id(), err)
	}

	data = api.UrlInformation(data["url"].(string))
	for k, v := range data {
		if validateInstanceSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return diag.Diagnostics{}
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	keys := []string{"name", "plan", "nodes", "tags"}
	params := make(map[string]any)

	if !d.HasChanges("name", "plan", "nodes", "tags") {
		return nil
	}

	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = d.Get(k)
		}
		if k == "nodes" {
			plan := d.Get("plan").(string)
			if isSharedPlan(plan) || !isLegacyPlan(plan) {
				delete(params, k)
			}
		}
	}

	if err := api.UpdateInstance(ctx, d.Id(), params); err != nil {
		return diag.FromErr(err)
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	if err := api.DeleteInstance(ctx, d.Id(), d.Get("keep_associated_vpc").(bool)); err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
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
		"ready",
		"backend",
		"vpc_id":
		return true
	}
	return false
}

func isSharedPlan(plan string) bool {
	switch plan {
	case
		"lemur",
		"tiger",
		"lemming",
		"ermine":
		return true
	}
	return false
}

func isLegacyPlan(plan string) bool {
	switch plan {
	case
		"bunny", "rabbit", "panda", "ape", "hippo", "lion":
		return true
	}
	return false
}

func instanceCreateAttributeKeys() []string {
	return []string{
		"name",
		"plan",
		"region",
		"nodes",
		"tags",
		"rmq_version",
		"vpc_id",
		"vpc_subnet",
		"no_default_alarms",
		"copy_settings",
		"preferred_az",
	}
}

func validateCopySettings() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"alarms",
		"config",
		"definitions",
		"firewall",
		"logs",
		"metrics",
		"plugins",
	}, true))
}

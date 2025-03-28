package cloudamqp

import (
	"context"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Identifier for the instance",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the instance",
			},
			"plan": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the plan, see documentation for valid plans",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the region you want to create your instance in",
			},
			"vpc_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the VPC to create your instance in",
			},
			"vpc_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dedicated VPC subnet, shouldn't overlap with your current VPC's subnet",
			},
			"nodes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of nodes in cluster (plan must support it)",
			},
			"rmq_version": {
				Type:        schema.TypeString,
				Computed:    true,
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
				Computed: true,
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
				Description: "If default alarms set or not for the instance",
			},
			"backend": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Software backend used, determined by subscription plan",
			},
		},
	}
}

func dataSourceInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = strconv.Itoa(d.Get("instance_id").(int))
	)

	data, err := api.ReadInstance(ctx, instanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceID)
	for k, v := range data {
		if validateInstanceSchemaAttribute(k) {
			if k == "vpc" {
				err = d.Set("vpc_id", v.(map[string]interface{})["id"])
				err = d.Set("vpc_subnet", v.(map[string]interface{})["subnet"])
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

	if data["no_default_alarms"] == nil {
		d.Set("no_default_alarms", false)
	}

	data = api.UrlInformation(data["url"].(string))
	for k, v := range data {
		if validateInstanceSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return nil
}

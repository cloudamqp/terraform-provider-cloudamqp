package cloudamqp

import (
	"fmt"
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInstanceRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Description: "Name of the plan, valid options are: lemur, tiger, bunny, rabbit, panda, ape, hippo, lion",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
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
				Description: "Host name for the CloudAMQP instance",
			},
			"vhost": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The virtual host",
			},
		},
	}
}

func dataSourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	instance_id := strconv.Itoa(d.Get("instance_id").(int))
	data, err := api.ReadInstance(instance_id)

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

package cloudamqp

import (
	"context"
	"errors"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVpcGcpInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcGcpInfoRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Instance identifier",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPC instance identifier",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC name",
			},
			"vpc_subnet": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC subnet",
			},
			"network": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC network uri",
			},
			"sleep": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Configurable sleep in seconds between retries when reading peering",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1800,
				Description: "Configurable timeout time (seconds) before retries times out",
			},
		},
	}
}

func dataSourceVpcGcpInfoRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api         = meta.(*api.API)
		data        = make(map[string]any)
		err         = errors.New("")
		instance_id = d.Get("instance_id").(int)
		vpc_id      = d.Get("vpc_id").(string)
		sleep       = d.Get("sleep").(int)
		timeout     = d.Get("timeout").(int)
	)

	if instance_id == 0 && vpc_id == "" {
		return diag.Errorf("you need to specify either instance_id or vpc_id")
	} else if instance_id != 0 {
		data, err = api.ReadVpcGcpInfo(ctx, instance_id, sleep, timeout)
	} else if d.Get("vpc_id") != nil {
		data, err = api.ReadVpcGcpInfoWithVpcId(ctx, vpc_id, sleep, timeout)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(data["name"].(string))

	for k, v := range data {
		if validateVpcGcpInfoSchemaAttribute(k) {
			if k == "subnet" {
				err = d.Set("vpc_subnet", v)
			} else {
				err = d.Set(k, v)
			}

			if err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return diag.Diagnostics{}
}

func validateVpcGcpInfoSchemaAttribute(key string) bool {
	switch key {
	case "name",
		"subnet",
		"vpc_subnet",
		"network":
		return true
	}
	return false
}

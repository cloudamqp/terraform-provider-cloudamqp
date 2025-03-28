package cloudamqp

import (
	"context"
	"errors"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVpcInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcInfoRead,

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
			"owner_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Owner identifier",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The security group identifier",
			},
		},
	}
}

func dataSourceVpcInfoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		vpcID      = d.Get("vpc_id").(string)
		data       map[string]interface{}
		err        = errors.New("")
	)

	if instanceID == 0 && vpcID == "" {
		return diag.Errorf("you need to specify either instance_id or vpc_id")
	} else if instanceID != 0 {
		data, err = api.ReadVpcInfo(ctx, instanceID)
	} else if d.Get("vpc_id") != nil {
		data, err = api.ReadVpcInfoWithVpcId(ctx, vpcID)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if data["id"] == nil {
		return diag.Errorf("failed to find external VPC identifier. Data source used for AWS VPC")
	}

	for k, v := range data {
		switch k {
		case "id":
			d.SetId(v.(string))
		case "name":
			d.Set("name", v)
		case "owner_id":
			d.Set("owner_id", v)
		case "subnet":
			d.Set("vpc_subnet", v)
		case "security_group":
			sg := data[k].(map[string]interface{})
			d.Set("security_group_id", sg["id"])
		}
	}

	return nil
}

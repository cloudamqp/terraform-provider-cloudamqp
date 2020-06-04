package cloudamqp

import (
	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVpcInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcInfoRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
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

func dataSourceVpcInfoRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data, err := api.ReadVpcInfo(d.Get("instance_id").(int))

	if err != nil {
		return err
	}

	d.SetId(data["id"].(string))
	for k, v := range data {
		if validateVpcInfoSchemaAttribute(k) {
			if k == "security_group" {
				sg := data[k].(map[string]interface{})
				d.Set("security_group_id", sg["id"])
			} else if k == "subnet" {
				d.Set("vpc_subnet", v)
			} else {

				d.Set(k, v)
			}
		}
	}
	return nil
}

func validateVpcInfoSchemaAttribute(key string) bool {
	switch key {
	case "name",
		"subnet",
		"vpc_subnet",
		"owner_id",
		"security_group",
		"security_group_id":
		return true
	}
	return false
}

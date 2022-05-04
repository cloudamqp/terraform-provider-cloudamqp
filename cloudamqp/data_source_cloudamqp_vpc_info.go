package cloudamqp

import (
	"errors"
	"fmt"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVpcInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcInfoRead,

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

func dataSourceVpcInfoRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)

	data := make(map[string]interface{})
	err := errors.New("")
	if d.Get("instance_id") == 0 && d.Get("vpc_id") == nil {
		return errors.New("You need to specify either instance_id or vpc_id")
	} else if d.Get("instance_id") != 0 {
		data, err = api.ReadVpcInfo(d.Get("instance_id").(int))
	} else if d.Get("vpc_id") != nil {
		data, err = api.ReadVpcInfoWithVpcId(d.Get("vpc_id").(string))
	}

	if err != nil {
		return err
	}

	d.SetId(data["id"].(string))
	for k, v := range data {
		if validateVpcInfoSchemaAttribute(k) {
			if k == "security_group" {
				sg := data[k].(map[string]interface{})
				err = d.Set("security_group_id", sg["id"])
			} else if k == "security_group_id" {
				continue
			} else if k == "subnet" {
				err = d.Set("vpc_subnet", v)
			} else {
				err = d.Set(k, v)
			}

			if err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
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

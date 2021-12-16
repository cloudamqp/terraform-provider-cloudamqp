package cloudamqp

import (
	"fmt"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVpcGcpInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcGcpInfoRead,

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
			"network": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC network uri",
			},
		},
	}
}

func dataSourceVpcGcpInfoRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data, err := api.ReadVpcGcpInfo(d.Get("instance_id").(int))

	if err != nil {
		return err
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
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return nil
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

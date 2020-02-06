package cloudamqp

import (
	"log"

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
	log.Printf("[DEBUG] cloudamqp::data_source::vpc_info::read instance id: %v", d.Get("instance_id"))
	data, err := api.ReadVpcInfo(d.Get("instance_id").(int))
	log.Printf("[DEBUG] cloudamqp::data_source::vpc_info::read data: %v", data)

	if err != nil {
		return err
	}
	d.SetId(data["id"].(string))
	d.Set("name", data["name"])
	d.Set("vpc_subnet", data["subnet"])
	d.Set("owner_id", data["owner_id"])
	sg := data["security_group"].(map[string]interface{})
	d.Set("security_group_id", sg["id"])
	return nil
}

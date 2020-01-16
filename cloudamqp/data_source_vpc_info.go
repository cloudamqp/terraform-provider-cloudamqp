package cloudamqp

import (
	"fmt"
	"log"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform/helper/schema"
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
				Optional:    true,
				Description: "VPC name",
			},
			"vpc_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VPC subnet",
			},
			"owner_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Owner identifier",
			},
			"security_group": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"owner_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
	setSecurityGroup(d, data["security_group"].(map[string]interface{}))
	return nil
}

func setSecurityGroup(d *schema.ResourceData, data map[string]interface{}) error {
	if data == nil {
		return fmt.Errorf("Unexpected nil pointer in: %s", data)
	}

	securityGroup := make([]map[string]interface{}, 0, len(data))
	securityGroup = append(securityGroup, map[string]interface{}{
		"id":          data["id"],
		"name":        data["name"],
		"description": data["description"],
		"owner_id":    data["owner_id"],
	})

	return d.Set("security_group", securityGroup)
}

package cloudamqp

import (
	"errors"
	"fmt"
	"log"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceVpcGcpInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVpcGcpInfoRead,

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

func dataSourceVpcGcpInfoRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api         = meta.(*api.API)
		data        = make(map[string]interface{})
		err         = errors.New("")
		instance_id = d.Get("instance_id").(int)
		vpc_id      = d.Get("vpc_id").(string)
		sleep       = d.Get("sleep").(int)
		timeout     = d.Get("timeout").(int)
	)

	log.Printf("[DEBUG] cloudamqp::data::vpc_gcp_info::request instance_id: %v, vpc_id: %v",
		instance_id, vpc_id)
	if instance_id == 0 && vpc_id == "" {
		return errors.New("you need to specify either instance_id or vpc_id")
	} else if instance_id != 0 {
		data, err = api.ReadVpcGcpInfo(instance_id, sleep, timeout)
	} else if d.Get("vpc_id") != nil {
		data, err = api.ReadVpcGcpInfoWithVpcId(vpc_id, sleep, timeout)
	}

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

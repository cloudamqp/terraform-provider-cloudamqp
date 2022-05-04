package cloudamqp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVpcPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcPeeringAccept,
		Read:   resourceVpcPeeringRead,
		Update: resourceVpcPeeringAccept,
		Delete: resourceVpcPeeringDelete,
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
			"peering_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC peering identifier",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC peering status",
			},
		},
	}
}

func resourceVpcPeeringAccept(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	err := errors.New("")
	if d.Get("instance_id") == 0 && d.Get("vpc_id") == nil {
		return errors.New("You need to specify either instance_id or vpc_id")
	} else if d.Get("instance_id") != 0 {
		_, err = api.AcceptVpcPeering(d.Get("instance_id").(int), d.Get("peering_id").(string))
	} else if d.Get("vpc_id") != nil {
		_, err = api.AcceptVpcPeeringTemp(d.Get("vpc_id").(string), d.Get("peering_id").(string))
	}

	if err != nil {
		return err
	}
	if d.Get("peering_id") != nil {
		d.SetId(d.Get("peering_id").(string))
	}

	return nil
}

func resourceVpcPeeringRead(d *schema.ResourceData, meta interface{}) error {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		if s[1] != "" {
			instanceID, _ := strconv.Atoi(s[1])
			d.Set("instance_id", instanceID)
		} else if s[2] != "" {
			vpcID, _ := strconv.Atoi(s[2])
			d.Set("vpc_id", vpcID)
		}
	}
	if d.Get("instance_id").(int) == 0 && d.Get("vpc_id").(string) == "" {
		return errors.New("Missing instance identifier: {resource_id},{instance_id}")
	}

	api := meta.(*api.API)
	err := errors.New("")
	data := make(map[string]interface{})
	if d.Get("instance_id") == 0 && d.Get("vpc_id") == nil {
		return errors.New("You need to specify either instance_id or vpc_id")
	} else if d.Get("instance_id") != 0 {
		data, err = api.ReadVpcPeeringRequest(d.Get("instance_id").(int), d.Get("peering_id").(string))
	} else if d.Get("vpc_id") != nil {
		data, err = api.ReadVpcPeeringRequestTemp(d.Get("vpc_id").(string), d.Get("peering_id").(string))
	}

	if err != nil {
		return err
	}
	for k, v := range data {
		if k == "status" {
			status := v.(map[string]interface{})
			if err = d.Set(k, status["code"]); err != nil {
				return fmt.Errorf("error setting status for resource %s: %s", d.Id(), err)
			}
		}
	}

	return nil
}

func resourceVpcPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	if d.Get("instance_id") == 0 && d.Get("vpc_id") == nil {
		return errors.New("You need to specify either instance_id or vpc_id")
	} else if d.Get("instance_id") != 0 {
		return api.RemoveVpcPeering(d.Get("instance_id").(int), d.Get("peering_id").(string))
	} else if d.Get("vpc_id") != nil {
		return api.RemoveVpcPeeringTemp(d.Get("vpc_id").(string), d.Get("peering_id").(string))
	}
	return errors.New("")
}

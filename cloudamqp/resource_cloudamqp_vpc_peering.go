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
				Required:    true,
				Description: "Instance identifier",
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
	_, err := api.AcceptVpcPeering(d.Get("instance_id").(int), d.Get("peering_id").(string))

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
		instanceID, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
	}
	if d.Get("instance_id").(int) == 0 {
		return errors.New("Missing instance identifier: {resource_id},{instance_id}")
	}

	api := meta.(*api.API)
	data, err := api.ReadVpcPeeringRequest(d.Get("instance_id").(int), d.Get("peering_id").(string))

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
	return api.RemoveVpcPeering(d.Get("instance_id").(int), d.Get("peering_id").(string))
}

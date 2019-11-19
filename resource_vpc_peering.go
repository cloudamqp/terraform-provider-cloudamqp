package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform/helper/schema"
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
		instance_id, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instance_id)
	}
	if d.Get("instance_id").(int) == 0 {
		return errors.New("Missing instance identifier: {resource_id},{instance_id}")
	}

	api := meta.(*api.API)
	_, err := api.ReadVpcPeeringRequest(d.Get("instance_id").(int), d.Get("peering_id").(string))

	if err != nil {
		return err
	}

	return nil
}

func resourceVpcPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	err := api.RemoveVpcPeering(d.Get("instance_id").(int), d.Get("peering_id").(string))

	if err != nil {
		return err
	}

	return nil
}

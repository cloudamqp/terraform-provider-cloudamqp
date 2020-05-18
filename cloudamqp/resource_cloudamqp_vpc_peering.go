package cloudamqp

import (
	"errors"
	"log"
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
	log.Printf("[DEBUG] cloudamqp::resource::vpc_peering::accept instance id: %v, peering id: %v", d.Get("instance_id"), d.Get("peering_id"))
	_, err := api.AcceptVpcPeering(d.Get("instance_id").(int), d.Get("peering_id").(string))

	if err != nil {
		return err
	}
	if d.Get("peering_id") != nil {
		d.SetId(d.Get("peering_id").(string))
		log.Printf("[DEBUG] cloudamqp::resource::vpc_peering::accept id set: %v", d.Id())
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
	log.Printf("[DEBUG] cloudamqp::resource::vpc_peering::read instance id: %v, peering id: %v", d.Get("instance_id"), d.Get("peering_id"))
	data, err := api.ReadVpcPeeringRequest(d.Get("instance_id").(int), d.Get("peering_id").(string))

	if err != nil {
		return err
	}
	for k, v := range data {
		if k == "status" {
			status := v.(map[string]interface{})
			d.Set(k, status["code"])
		}
	}

	return nil
}

func resourceVpcPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	log.Printf("[DEBUG] cloudamqp::resource::vpc_peering::delete instance id: %v, peering id: %v", d.Get("instance_id"), d.Get("peering_id"))
	err := api.RemoveVpcPeering(d.Get("instance_id").(int), d.Get("peering_id").(string))

	if err != nil {
		return err
	}

	return nil
}

package cloudamqp

import (
	"errors"
	"log"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVpcGcpPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateVpcGcpPeering,
		Read:   resourceReadVpcGcpPeering,
		Update: resourceUpdateVpcGcpPeering,
		Delete: resourceDeleteVpcGcpPeering,
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
			"peer_network_uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC network uri",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC peering state",
			},
			"state_details": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC peering state details",
			},
			"auto_create_routes": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "VPC peering auto created routes",
			},
		},
	}
}

func resourceCreateVpcGcpPeering(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"peer_network_uri"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil && v != "" {
			params[k] = v
		}
	}

	// Todo: Create if/else if/else to check either instance_id or vpc_id. Different calls!
	// data, err := api.RequestVpcGcpPeering(d.Get("instance_id").(int), params)
	log.Printf("[DEBUG] cloudamqp::vpc_gcp_peering::create instance_id: %v, vpc_id: %v", d.Get("instance_id"), d.Get("vpc_id"))
	data := make(map[string]interface{})
	err := errors.New("")
	if d.Get("instance_id") == 0 && d.Get("vpc_id") == nil {
		return errors.New("You need to specify either instance_id or vpc_id")
	} else if d.Get("instance_id") != 0 {
		data, err = api.RequestVpcGcpPeering(d.Get("instance_id").(int), params)
	} else if d.Get("vpc_id") != nil {
		data, err = api.RequestVpcGcpPeeringTemp(d.Get("vpc_id").(string), params)
	}

	if err != nil {
		return err
	}

	for k, v := range data {
		if k == "peering" {
			d.SetId(v.(string))
		}
	}

	return resourceReadVpcGcpPeering(d, meta)
}

func resourceReadVpcGcpPeering(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	// Todo: Create if/else if/else to check either instance_id or vpc_id. Different calls!
	// data, err := api.ReadVpcGcpPeering(d.Get("instance_id").(int), d.Id())
	data := make(map[string]interface{})
	err := errors.New("")
	if d.Get("instance_id") == 0 && d.Get("vpc_id") == nil {
		return errors.New("You need to specify either instance_id or vpc_id")
	} else if d.Get("instance_id") != 0 {
		data, err = api.ReadVpcGcpPeering(d.Get("instance_id").(int), d.Id())
	} else if d.Get("vpc_id") != nil {
		data, err = api.ReadVpcGcpPeeringTemp(d.Get("vpc_id").(string), d.Id())
	}
	log.Printf("[DEBUG] cloudamqp::vpc_gcp_peering::read data: %v", data)

	if err != nil {
		return err
	}
	// Todo: Needed?
	if data["rows"] == nil {
		return errors.New("No peering data available")
	}

	rows := data["rows"].([]interface{})
	if len(rows) > 0 {
		for _, row := range rows {
			tempRow := row.(map[string]interface{})
			if tempRow["name"] != d.Id() {
				continue
			}
			for k, v := range tempRow {
				if validateGcpPeeringSchemaAttribute(k) {
					if k == "stateDetails" {
						d.Set("state_details", v.(string))
					} else if k == "autoCreateRoutes" {
						d.Set("auto_create_routes", v.(bool))
					} else {
						d.Set(k, v.(string))
					}
				}
			}
		}
	}

	return nil
}

func resourceUpdateVpcGcpPeering(d *schema.ResourceData, meta interface{}) error {
	return resourceReadVpcGcpPeering(d, meta)
}

func resourceDeleteVpcGcpPeering(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	// Todo: Create if/else if/else to check either instance_id or vpc_id. Different calls!
	// return api.RemoveVpcGcpPeering(d.Get("instance_id").(int), d.Id())
	if d.Get("instance_id") == 0 && d.Get("vpc_id") == nil {
		return errors.New("You need to specify either instance_id or vpc_id")
	} else if d.Get("instance_id") != 0 {
		return api.RemoveVpcGcpPeering(d.Get("instance_id").(int), d.Id())
	} else if d.Get("vpc_id") != nil {
		return api.RemoveVpcGcpPeeringTemp(d.Get("vpc_id").(string), d.Id())
	}
	return errors.New("Failed to remove VPC peering")
}

func validateGcpPeeringSchemaAttribute(key string) bool {
	switch key {
	case "state",
		"stateDetails",
		"autoCreateRoutes":
		return true
	}
	return false
}

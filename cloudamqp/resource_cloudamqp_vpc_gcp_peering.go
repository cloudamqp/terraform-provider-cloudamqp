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
			"wait_on_peering_status": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Wait until peering status change to 'connected'",
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
			"sleep": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Configurable sleep in seconds between retries when requesting or reading peering",
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

func resourceCreateVpcGcpPeering(d *schema.ResourceData, meta interface{}) error {
	var (
		api          = meta.(*api.API)
		keys         = []string{"peer_network_uri"}
		params       = make(map[string]interface{})
		instanceID   = d.Get("instance_id").(int)
		vpcID        = d.Get("vpc_id").(string)
		waitOnStatus = d.Get("wait_on_peering_status").(bool)
		data         map[string]interface{}
		err          = errors.New("")
		sleep        = d.Get("sleep").(int)
		timeout      = d.Get("timeout").(int)
	)

	for _, k := range keys {
		if v := d.Get(k); v != nil && v != "" {
			params[k] = v
		}
	}

	log.Printf("[DEBUG] cloudamqp::vpc_gcp_peering::create instance_id: %d, vpc_id: %s, "+
		"waitOnStatus: %v", instanceID, vpcID, waitOnStatus)
	if instanceID == 0 && vpcID == "" {
		return errors.New("you need to specify either instance_id or vpc_id")
	} else if instanceID != 0 {
		data, err = api.RequestVpcGcpPeering(instanceID, params, waitOnStatus, sleep, timeout)
	} else if vpcID != "" {
		data, err = api.RequestVpcGcpPeeringWithVpcId(vpcID, params, waitOnStatus, sleep, timeout)
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
	var (
		api     = meta.(*api.API)
		data    map[string]interface{}
		err     = errors.New("")
		sleep   = d.Get("sleep").(int)
		timeout = d.Get("timeout").(int)
	)

	if d.Get("instance_id") == 0 && d.Get("vpc_id") == nil {
		return errors.New("you need to specify either instance_id or vpc_id")
	} else if d.Get("instance_id") != 0 {
		data, err = api.ReadVpcGcpPeering(d.Get("instance_id").(int), sleep, timeout)
	} else if d.Get("vpc_id") != nil {
		data, err = api.ReadVpcGcpPeeringWithVpcId(d.Get("vpc_id").(string), sleep, timeout)
	}
	log.Printf("[DEBUG] cloudamqp::vpc_gcp_peering::read data: %v", data)

	if err != nil {
		return err
	}

	if data["rows"] == nil {
		return errors.New("no peering data available")
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
	if d.Get("instance_id") == 0 && d.Get("vpc_id") == nil {
		return errors.New("you need to specify either instance_id or vpc_id")
	} else if d.Get("instance_id") != 0 {
		return api.RemoveVpcGcpPeering(d.Get("instance_id").(int), d.Id())
	} else if d.Get("vpc_id") != nil {
		return api.RemoveVpcGcpPeeringWithVpcId(d.Get("vpc_id").(string), d.Id())
	}
	return errors.New("failed to remove VPC peering")
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

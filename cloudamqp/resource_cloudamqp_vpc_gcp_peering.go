package cloudamqp

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVpcGcpPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreateVpcGcpPeering,
		Read:   resourceReadVpcGcpPeering,
		Update: resourceUpdateVpcGcpPeering,
		Delete: resourceDeleteVpcGcpPeering,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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

	log.Printf("[DEBUG] cloudamqp::vpc_gcp_peering#create instance_id: %d, vpc_id: %s, "+
		"waitOnStatus: %v", instanceID, vpcID, waitOnStatus)

	for _, k := range keys {
		if v := d.Get(k); v != nil && v != "" {
			params[k] = v
		}
	}

	if instanceID == 0 && vpcID == "" {
		return errors.New("you need to specify either instance_id or vpc_id")
	} else if instanceID != 0 {
		data, err = api.RequestVpcGcpPeering(instanceID, params, waitOnStatus, sleep, timeout)
	} else if vpcID != "" {
		data, err = api.RequestVpcGcpPeeringWithVpcId(vpcID, params, waitOnStatus, sleep, timeout)
	}

	log.Printf("[DEBUG] cloudamqp::vpc_gcp_peering#create data: %v", data)

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
		api        = meta.(*api.API)
		data       map[string]interface{}
		err        = errors.New("")
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
		instanceID = d.Get("instance_id").(int)
		vpcID      = d.Get("vpc_id").(string)
	)

	// Check to determine if the resource should be imported.
	if strings.Contains(d.Id(), ",") {
		log.Printf("[DEBUG] cloudamqp::resource::vpc_gcp_peering#read import detected")
		return resourceImportVpcGcpPeering(d, meta)
	}

	if instanceID == 0 && vpcID == "" {
		return errors.New("you need to specify either instance_id or vpc_id")
	} else if instanceID != 0 {
		data, err = api.ReadVpcGcpPeering(instanceID, sleep, timeout)
	} else if vpcID != "" {
		data, err = api.ReadVpcGcpPeeringWithVpcId(vpcID, sleep, timeout)
	}

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] cloudamqp::vpc_gcp_peering#read data: %v", data)

	rows := data["rows"].([]interface{})
	if len(rows) > 0 {
		for _, row := range rows {
			tempRow := row.(map[string]interface{})
			if tempRow["name"] != d.Id() {
				continue
			}
			for k, v := range tempRow {
				switch k {
				case "autoCreateRoutes":
					d.Set("auto_create_routes", v.(bool))
				case "network":
					d.Set("peer_network_uri", v.(string))
				case "state":
					d.Set("state", v.(string))
				case "stateDetails":
					d.Set("state_details", v.(string))
				}
			}
		}
	} else {
		return errors.New("no peering data available")
	}

	return nil
}

func resourceImportVpcGcpPeering(d *schema.ResourceData, meta interface{}) error {
	var (
		api     = meta.(*api.API)
		data    map[string]interface{}
		err     = errors.New("")
		sleep   = 10
		timeout = 1800
	)

	// Set default values to arguments
	d.Set("sleep", 10)
	d.Set("timeout", 1800)
	d.Set("wait_on_peering_status", false)

	log.Printf("[DEBUG] cloudamqp::resource::vpc_gcp_peering#import id: %v", d.Id())
	importValues := strings.Split(d.Id(), ",")
	if len(importValues) < 3 {
		return errors.New("wrong number of import argument, need all three <type>,<id>,<peer_network_uri>")
	}
	if importValues[0] == "instance" {
		instanceID, _ := strconv.Atoi(importValues[1])
		d.Set("instance_id", instanceID)
		data, err = api.ReadVpcGcpPeering(instanceID, sleep, timeout)
	} else if importValues[0] == "vpc" {
		vpcID := importValues[1]
		d.Set("vpc_id", vpcID)
		data, err = api.ReadVpcGcpPeeringWithVpcId(vpcID, sleep, timeout)
	}

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] cloudamqp::vpc_gcp_peering#import data: %v", data)

	rows := data["rows"].([]interface{})
	if len(rows) > 0 {
		for _, row := range rows {
			tempRow := row.(map[string]interface{})
			if tempRow["network"] != importValues[2] {
				continue
			}
			for k, v := range tempRow {
				switch k {
				case "autoCreateRoutes":
					d.Set("auto_create_routes", v.(bool))
				case "name":
					d.SetId(v.(string))
				case "network":
					d.Set("peer_network_uri", v.(string))
				case "state":
					d.Set("state", v.(string))
				case "stateDetails":
					d.Set("state_details", v.(string))
				}
			}
		}
	} else {
		return errors.New("no peering data available")
	}

	return nil
}

func resourceUpdateVpcGcpPeering(d *schema.ResourceData, meta interface{}) error {
	return resourceReadVpcGcpPeering(d, meta)
}

func resourceDeleteVpcGcpPeering(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		vpcID      = d.Get("vpc_id").(string)
	)

	if instanceID == 0 && vpcID == "" {
		return errors.New("you need to specify either instance_id or vpc_id")
	} else if instanceID != 0 {
		return api.RemoveVpcGcpPeering(instanceID, d.Id())
	} else if vpcID != "" {
		return api.RemoveVpcGcpPeeringWithVpcId(vpcID, d.Id())
	}
	return errors.New("failed to remove VPC peering")
}

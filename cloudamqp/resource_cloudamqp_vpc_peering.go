package cloudamqp

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVpcPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcPeeringAccept,
		Read:   resourceVpcPeeringRead,
		Update: resourceVpcPeeringAccept,
		Delete: resourceVpcPeeringDelete,
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
			"sleep": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "Configurable sleep time in seconds between retries for accepting or removing peering",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "Configurable timeout time in seconds for accepting or removing peering",
			},
		},
	}
}

func resourceVpcPeeringAccept(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		err        = errors.New("")
		instanceID = d.Get("instance_id").(int)
		vpcID      = d.Get("vpc_id").(string)
		peeringID  = d.Get("peering_id").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	log.Printf("[DEBUG] cloudamqp::resource::vpc_aws_peering#accept instance_id: %d, vpc_id: %s, "+
		"peering_id: %s", instanceID, vpcID, peeringID)

	if instanceID == 0 && vpcID == "" {
		return errors.New("you need to specify either instance_id or vpc_id")
	} else if instanceID != 0 {
		_, err = api.AcceptVpcPeering(instanceID, peeringID, sleep, timeout)
	} else if vpcID != "" {
		_, err = api.AcceptVpcPeeringWithVpcId(vpcID, peeringID, sleep, timeout)
	}

	if err != nil {
		return err
	}

	return resourceVpcPeeringRead(d, meta)
}

func resourceVpcPeeringRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		vpcID      = d.Get("vpc_id").(string)
		peeringID  = d.Get("peering_id").(string)
		data       map[string]interface{}
		err        = errors.New("")
	)

	// Check to determine if the resource should be imported.
	if strings.Contains(d.Id(), ",") {
		log.Printf("[DEBUG] cloudamqp::resource::vpc_aws_peering#read import detected")
		return resourceVpcPeeringImport(d, meta)
	}

	if instanceID == 0 && vpcID == "" {
		return errors.New("you need to specify either instance_id or vpc_id")
	} else if instanceID != 0 {
		data, err = api.ReadVpcPeeringRequest(instanceID, peeringID)
	} else if d.Get("vpc_id") != nil {
		data, err = api.ReadVpcPeeringRequestWithVpcId(vpcID, peeringID)
	}

	if err != nil {
		return err
	}

	for k, v := range data {
		switch k {
		case "vpc_peering_connection_id":
			d.SetId(v.(string))
		case "status":
			status := v.(map[string]interface{})
			if err = d.Set(k, status["code"]); err != nil {
				return fmt.Errorf("error setting status for resource %s: %s", d.Id(), err)
			}
		}
	}

	return nil
}

func resourceVpcPeeringImport(d *schema.ResourceData, meta interface{}) error {
	var (
		api  = meta.(*api.API)
		data map[string]interface{}
		err  = errors.New("")
	)

	// Set default values to arguments
	d.Set("sleep", 60)
	d.Set("timeout", 3600)

	log.Printf("[DEBUG] cloudamqp::resource::vpc_aws_peering#import id: %v", d.Id())
	importValues := strings.Split(d.Id(), ",")
	if len(importValues) < 3 {
		return errors.New("wrong number of import argument, need all three <type>,<id>,<peering_id>")
	}
	peeringID := importValues[2]
	d.Set("peering_id", peeringID)
	if importValues[0] == "instance" {
		instanceID, _ := strconv.Atoi(importValues[1])
		d.Set("instance_id", instanceID)
		data, err = api.ReadVpcPeeringRequest(instanceID, peeringID)
	} else if importValues[0] == "vpc" {
		vpcID := importValues[1]
		d.Set("vpc_id", vpcID)
		data, err = api.ReadVpcPeeringRequestWithVpcId(vpcID, peeringID)
	}

	if err != nil {
		return err
	}

	for k, v := range data {
		switch k {
		case "vpc_peering_connection_id":
			d.SetId(v.(string))
		case "status":
			status := v.(map[string]interface{})
			if err = d.Set(k, status["code"]); err != nil {
				return fmt.Errorf("error setting status for resource %s: %s", d.Id(), err)
			}
		}
	}

	return nil
}

func resourceVpcPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		vpcID      = d.Get("vpc_id").(string)
		peeringID  = d.Get("peering_id").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	if instanceID == 0 && vpcID == "" {
		return errors.New("you need to specify either instance_id or vpc_id")
	} else if instanceID != 0 {
		return api.RemoveVpcPeering(instanceID, peeringID, sleep, timeout)
	} else if vpcID != "" {
		return api.RemoveVpcPeeringWithVpcId(vpcID, peeringID, sleep, timeout)
	}

	return nil
}

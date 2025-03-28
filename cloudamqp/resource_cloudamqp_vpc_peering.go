package cloudamqp

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVpcPeering() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcPeeringAccept,
		ReadContext:   resourceVpcPeeringRead,
		UpdateContext: resourceVpcPeeringAccept,
		DeleteContext: resourceVpcPeeringDelete,
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

func resourceVpcPeeringAccept(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		err        = errors.New("")
		instanceID = d.Get("instance_id").(int)
		vpcID      = d.Get("vpc_id").(string)
		peeringID  = d.Get("peering_id").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	if instanceID == 0 && vpcID == "" {
		return diag.Errorf("you need to specify either instance_id or vpc_id")
	} else if instanceID != 0 {
		_, err = api.AcceptVpcPeering(ctx, instanceID, peeringID, sleep, timeout)
	} else if vpcID != "" {
		_, err = api.AcceptVpcPeeringWithVpcId(ctx, vpcID, peeringID, sleep, timeout)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceVpcPeeringRead(ctx, d, meta)
}

func resourceVpcPeeringRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		vpcID      = d.Get("vpc_id").(string)
		peeringID  = d.Get("peering_id").(string)
		data       map[string]any
		err        = errors.New("")
	)

	// Check to determine if the resource should be imported.
	if strings.Contains(d.Id(), ",") {
		tflog.Info(ctx, fmt.Sprintf("import of resource with identifiers: %s", d.Id()))
		if diag := resourceVpcPeeringImport(ctx, d, meta); diag != nil {
			return diag
		}
	}

	if instanceID == 0 && vpcID == "" {
		return diag.Errorf("you need to specify either instance_id or vpc_id")
	} else if instanceID != 0 {
		data, err = api.ReadVpcPeeringRequest(ctx, instanceID, peeringID)
	} else if d.Get("vpc_id") != nil {
		data, err = api.ReadVpcPeeringRequestWithVpcId(ctx, vpcID, peeringID)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range data {
		switch k {
		case "vpc_peering_connection_id":
			d.SetId(v.(string))
		case "status":
			status := v.(map[string]any)
			if err = d.Set(k, status["code"]); err != nil {
				return diag.Errorf("error setting status for resource %s: %s", d.Id(), err)
			}
		}
	}

	return diag.Diagnostics{}
}

func resourceVpcPeeringImport(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api  = meta.(*api.API)
		data map[string]any
		err  = errors.New("")
	)

	// Set default values to arguments
	d.Set("sleep", 60)
	d.Set("timeout", 3600)

	importValues := strings.Split(d.Id(), ",")
	if len(importValues) < 3 {
		return diag.Errorf("wrong number of import argument, need all three <type>,<id>,<peering_id>")
	}
	peeringID := importValues[2]
	d.Set("peering_id", peeringID)
	if importValues[0] == "instance" {
		instanceID, _ := strconv.Atoi(importValues[1])
		d.Set("instance_id", instanceID)
		data, err = api.ReadVpcPeeringRequest(ctx, instanceID, peeringID)
	} else if importValues[0] == "vpc" {
		vpcID := importValues[1]
		d.Set("vpc_id", vpcID)
		data, err = api.ReadVpcPeeringRequestWithVpcId(ctx, vpcID, peeringID)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range data {
		switch k {
		case "vpc_peering_connection_id":
			d.SetId(v.(string))
		case "status":
			status := v.(map[string]any)
			if err = d.Set(k, status["code"]); err != nil {
				return diag.Errorf("error setting status for resource %s: %s", d.Id(), err)
			}
		}
	}

	return diag.Diagnostics{}
}

func resourceVpcPeeringDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		vpcID      = d.Get("vpc_id").(string)
		peeringID  = d.Get("peering_id").(string)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
	)

	if instanceID == 0 && vpcID == "" {
		return diag.Errorf("you need to specify either instance_id or vpc_id")
	} else if instanceID != 0 {
		if err := api.RemoveVpcPeering(ctx, instanceID, peeringID, sleep, timeout); err != nil {
			return diag.FromErr(err)
		}
	} else if vpcID != "" {
		if err := api.RemoveVpcPeeringWithVpcId(ctx, vpcID, peeringID, sleep, timeout); err != nil {
			return diag.FromErr(err)
		}
	}

	return diag.Diagnostics{}
}

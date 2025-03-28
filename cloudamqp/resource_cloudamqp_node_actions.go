package cloudamqp

import (
	"context"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNodeAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeActionRequest,
		UpdateContext: resourceNodeActionRequest,
		ReadContext:   resourceNodeActionRead,
		DeleteContext: resourceNodeActionRemove,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"node_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the node",
			},
			"action": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The action to perform on the node",
				ValidateDiagFunc: validateNodeAction(),
			},
			"running": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If the node is running",
			},
		},
	}
}

func resourceNodeActionRequest(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	nodeName := d.Get("node_name").(string)
	data, err := api.PostAction(ctx, d.Get("instance_id").(int), nodeName, d.Get("action").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(nodeName)
	d.Set("running", data["running"])
	return diag.Diagnostics{}
}

func resourceNodeActionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	data, err := api.ReadNode(ctx, d.Get("instance_id").(int), d.Get("node_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("running", data["running"])
	return diag.Diagnostics{}
}

func resourceNodeActionRemove(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return diag.Diagnostics{}
}

func validateNodeAction() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"start",
		"stop",
		"restart",
		"reboot",
		"mgmt.restart",
	}, true))
}

package cloudamqp

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAccountAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountActionRequest,
		UpdateContext: resourceAccountActionRequest,
		ReadContext:   resourceAccountActionRead,
		DeleteContext: resourceAccountActionRemove,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"action": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The action to perform on the node",
				ValidateDiagFunc: validateAccountAction(),
			},
		},
	}
}

func resourceAccountActionRequest(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {

	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		action     = d.Get("action")
		err        = errors.New("")
	)

	switch action {
	case "rotate-password":
		err = api.RotatePassword(ctx, instanceID)
	case "rotate-apikey":
		err = api.RotateApiKey(ctx, instanceID)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", instanceID))
	return diag.Diagnostics{}
}

func resourceAccountActionRead(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {

	return diag.Diagnostics{}
}

func resourceAccountActionRemove(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {

	return diag.Diagnostics{}
}

func validateAccountAction() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"rotate-password",
		"rotate-apikey",
	}, true))
}

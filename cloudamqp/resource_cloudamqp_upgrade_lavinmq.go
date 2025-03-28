package cloudamqp

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUpgradeLavinMQ() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUpgradeLavinMQInvoke,
		ReadContext:   resourceUpgradeLavinMQRead,
		UpdateContext: resourceUpgradeLavinMQUpdate,
		DeleteContext: resourceUpgradeLavinMQRemove,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The CloudAMQP instance identifier",
			},
			"new_version": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The new version to upgrade to",
			},
		},
	}
}

func resourceUpgradeLavinMQInvoke(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		api         = meta.(*api.API)
		instanceID  = d.Get("instance_id").(int)
		new_version = d.Get("new_version").(string)
	)

	response, err := api.UpgradeLavinMQ(ctx, instanceID, new_version)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(instanceID))

	if len(response) > 0 {
		tflog.Info(ctx, fmt.Sprintf("LavinMQ update result: %s", response))
	}

	return diag.Diagnostics{}
}

func resourceUpgradeLavinMQRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceUpgradeLavinMQUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceUpgradeLavinMQRemove(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Diagnostics{}
}

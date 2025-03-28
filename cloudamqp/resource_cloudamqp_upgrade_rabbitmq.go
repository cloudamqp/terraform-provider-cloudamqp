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

func resourceUpgradeRabbitMQ() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUpgradeRabbitMQInvoke,
		ReadContext:   resourceUpgradeRabbitMQRead,
		UpdateContext: resourceUpgradeRabbitMQUpdate,
		DeleteContext: resourceUpgradeRabbitMQRemove,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The CloudAMQP instance identifier",
			},
			"current_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Helper argument to change upgrade behaviour to latest possible version",
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

func resourceUpgradeRabbitMQInvoke(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api             = meta.(*api.API)
		instanceID      = d.Get("instance_id").(int)
		current_version = d.Get("current_version").(string)
		new_version     = d.Get("new_version").(string)
	)

	response, err := api.UpgradeRabbitMQ(ctx, instanceID, current_version, new_version)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(instanceID))

	if len(response) > 0 {
		tflog.Info(ctx, fmt.Sprintf("RabbitMQ update result: %s", response))
	}

	return diag.Diagnostics{}
}

func resourceUpgradeRabbitMQRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceUpgradeRabbitMQUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return diag.Diagnostics{}
}

func resourceUpgradeRabbitMQRemove(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	return diag.Diagnostics{}
}

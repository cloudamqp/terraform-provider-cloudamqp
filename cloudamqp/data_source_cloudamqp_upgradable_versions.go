package cloudamqp

import (
	"context"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUpgradableVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUpgradableVersionRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"new_rabbitmq_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latest possible upgradable RabbitMQ version",
			},
			"new_erlang_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latest possible upgradable Erlang version",
			},
		},
	}
}

func dataSourceUpgradableVersionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	data, err := api.ReadVersions(ctx, d.Get("instance_id").(int))
	if err != nil {
		return diag.FromErr(err)
	}
	instanceID := strconv.Itoa(d.Get("instance_id").(int))
	d.SetId(instanceID)

	for k, v := range data {
		if validateVersionsSchemaAttribute(k) {
			d.Set(k, v)
		}
	}

	return diag.Diagnostics{}
}

func validateVersionsSchemaAttribute(key string) bool {
	switch key {
	case "new_rabbitmq_version",
		"new_erlang_version":
		return true
	}
	return false
}

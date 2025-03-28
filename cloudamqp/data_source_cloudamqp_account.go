package cloudamqp

import (
	"context"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountRead,

		Schema: map[string]*schema.Schema{
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The instance identifier",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the instance",
						},
						"plan": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subscription plan used for the instance",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region were the instanece is located in",
						},
						"tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Tag for the instance",
						},
					},
				},
			},
		},
	}
}

func dataSourceAccountRead(ctx context.Context, d *schema.ResourceData,
	meta any) diag.Diagnostics {

	api := meta.(*api.API)
	data, err := api.ListInstances(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("noId")
	instances := make([]map[string]any, len(data))
	for k, v := range data {
		instances[k] = readAccount(v)
	}

	if err = d.Set("instances", instances); err != nil {
		return diag.Errorf("error setting instances for resource %s, %s", d.Id(), err)
	}

	return diag.Diagnostics{}
}

func readAccount(data map[string]any) map[string]any {
	instance := make(map[string]any)
	for k, v := range data {
		if validateAccountSchemaAttribute(k) {
			instance[k] = v
		}
	}
	return instance
}

func validateAccountSchemaAttribute(key string) bool {
	switch key {
	case "id",
		"name",
		"plan",
		"region",
		"tags":
		return true
	}
	return false
}

package cloudamqp

import (
	"context"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccountVpcs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountVpcsRead,

		Schema: map[string]*schema.Schema{
			"vpcs": {
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
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region were the instanece is located in",
						},
						"subnet": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The VPC subnet",
						},
						"tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Tag the VPC instance with optional tags",
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC name given when hosted at the cloud provider",
						},
					},
				},
			},
		},
	}
}

func dataSourceAccountVpcsRead(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {

	api := meta.(*api.API)
	data, err := api.ListVpcs(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("noId")
	vpcs := make([]map[string]interface{}, len(data))
	for k, v := range data {
		vpcs[k] = readAccountVpc(v)
	}

	if err = d.Set("vpcs", vpcs); err != nil {
		return diag.Errorf("error setting vpcs for resource %s, %s", d.Id(), err)
	}

	return diag.Diagnostics{}
}

func readAccountVpc(data map[string]interface{}) map[string]interface{} {
	vpc := make(map[string]interface{})
	for k, v := range data {
		if validateAccountVpcsSchemaAttribute(k) {
			vpc[k] = v
		}
	}
	return vpc
}

func validateAccountVpcsSchemaAttribute(key string) bool {
	switch key {
	case "id",
		"name",
		"region",
		"subnet",
		"tags",
		"vpc_name":
		return true
	}
	return false
}

package cloudamqp

import (
	"context"
	"time"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Note: Cannot yet be migrated to framework while using "vpcs: schema.ListNestedAttribute" to build
// up the data source schema.
// Error: Failed to load plugin schema: AttributeName("vpcs"): protocol version 5 cannot have Attributes set..
// Makes the provider crash when loading the provider.
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
	meta any) diag.Diagnostics {

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	api := meta.(*api.API)
	data, err := api.ListVpcs(timeoutCtx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("noId")
	vpcs := make([]map[string]any, len(data))
	for k, vpcData := range data {
		vpc := map[string]any{
			"id":       vpcData.ID,
			"name":     vpcData.Name,
			"region":   vpcData.Region,
			"subnet":   vpcData.Subnet,
			"tags":     vpcData.Tags,
			"vpc_name": vpcData.VpcName,
		}
		vpcs[k] = vpc
	}

	if err = d.Set("vpcs", vpcs); err != nil {
		return diag.Errorf("error setting vpcs for resource %s, %s", d.Id(), err)
	}

	return diag.Diagnostics{}
}

package cloudamqp

import (
	"fmt"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccountVpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAccountVpcsRead,

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

func dataSourceAccountVpcsRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data, err := api.ListVpcs()
	if err != nil {
		return err
	} else if len(data) == 0 {
		return fmt.Errorf("no vpcs found for resoruce %s", d.Id())
	}
	d.SetId("noId")
	vpcs := make([]map[string]interface{}, len(data))
	for k, v := range data {
		vpcs[k] = readAccountVpc(v)
	}

	if err = d.Set("vpcs", vpcs); err != nil {
		return fmt.Errorf("error setting vpcs for resource %s, %s", d.Id(), err)
	}

	return nil
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

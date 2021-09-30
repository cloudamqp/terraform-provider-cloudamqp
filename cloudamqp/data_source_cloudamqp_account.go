package cloudamqp

import (
	"fmt"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAccountRead,

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

func dataSourceAccountRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data, err := api.ListInstances()
	if err != nil {
		return err
	} else if len(data) == 0 {
		return fmt.Errorf("no instances found for resoruce %s", d.Id())
	}
	d.SetId("na")
	instances := make([]map[string]interface{}, len(data))
	for k, v := range data {
		instances[k] = readAccount(v)
	}

	if err = d.Set("instances", instances); err != nil {
		return fmt.Errorf("error setting instances for resource %s, %s", d.Id(), err)
	}

	return nil
}

func readAccount(data map[string]interface{}) map[string]interface{} {
	instance := make(map[string]interface{})
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

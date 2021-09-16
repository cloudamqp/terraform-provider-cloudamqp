package cloudamqp

import (
	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"fmt"
)

func dataSourceCredentials() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCredentialsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"username": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceCredentialsRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data, err := api.ReadCredentials(d.Get("instance_id").(int))
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%v.%s", d.Get("instance_id").(int), data["username"]))
	for k, v := range data {
		if validateCredentialsSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return nil
}

func validateCredentialsSchemaAttribute(key string) bool {
	switch key {
	case "username",
		"password":
		return true
	}
	return false
}

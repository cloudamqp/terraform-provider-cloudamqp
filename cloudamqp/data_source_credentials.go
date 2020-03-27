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
				Optional:  true,
				Sensitive: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
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
	d.Set("username", data["username"])
	d.Set("password", data["password"])
	return nil
}

package cloudamqp

import (
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUpgradableVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUpgradableVersionRead,

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

func dataSourceUpgradableVersionRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data, err := api.ReadVersions(d.Get("instance_id").(int))
	if err != nil {
		return err
	}
	instanceID := strconv.Itoa(d.Get("instance_id").(int))
	d.SetId(instanceID)

	for k, v := range data {
		if validateVersionsSchemaAttribute(k) {
			d.Set(k, v)
		}
	}

	return nil
}

func validateVersionsSchemaAttribute(key string) bool {
	switch key {
	case "new_rabbitmq_version",
		"new_erlang_version":
		return true
	}
	return false
}

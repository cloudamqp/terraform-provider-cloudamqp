package cloudamqp

import (
	"fmt"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNodesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"running": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"rabbitmq_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"erlang_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hipe": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"configured": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"additional_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNodesRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data, err := api.ListNodes(d.Get("instance_id").(int))
	if err != nil {
		return err
	}
	instanceID := strconv.Itoa(d.Get("instance_id").(int))
	d.SetId(instanceID)

	nodes := make([]map[string]interface{}, len(data))
	for k, v := range data {
		nodes[k] = readNode(v)
	}

	if err = d.Set("nodes", nodes); err != nil {
		return fmt.Errorf("error setting nodes for resource %s, %s", d.Id(), err)
	}

	return nil
}

func readNode(data map[string]interface{}) map[string]interface{} {
	node := make(map[string]interface{})
	for k, v := range data {
		if validateNodesSchemaAttribute(k) {
			node[k] = v
		}
	}
	return node
}

func validateNodesSchemaAttribute(key string) bool {
	switch key {
	case "hostname",
		"name",
		"running",
		"rabbitmq_version",
		"erlang_version",
		"hipe",
		"configured",
		"rmq_version",
		"disk_size",
		"additional_disk_size":
		return true
	}
	return false
}

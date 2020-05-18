package cloudamqp

import (
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"running": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"rabbitmq_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"erlang_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"hipe": {
							Type:     schema.TypeBool,
							Optional: true,
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
	data, err := api.ReadNodes(d.Get("instance_id").(int))
	if err != nil {
		return err
	}
	id := strconv.Itoa(d.Get("instance_id").(int))
	d.SetId(id)
	d.Set("nodes", data)
	return nil
}

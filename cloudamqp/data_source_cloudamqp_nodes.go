package cloudamqp

import (
	"context"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodesRead,

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
						"hostname_internal": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNodesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	data, err := api.ListNodes(ctx, int64(d.Get("instance_id").(int)))
	if err != nil {
		return diag.FromErr(err)
	}
	instanceID := strconv.Itoa(d.Get("instance_id").(int))
	d.SetId(instanceID)

	nodes := make([]map[string]any, len(data))
	for k, v := range data {
		node := map[string]any{}
		node["additional_disk_size"] = v.AdditionalDiskSize
		node["availability_zone"] = v.AvailabilityZone
		node["configured"] = v.Configured
		node["disk_size"] = v.DiskSize
		node["erlang_version"] = v.ErlangVersion
		node["hostname"] = v.Hostname
		node["hostname_internal"] = v.HostnameInternal
		node["hipe"] = v.Hipe
		node["name"] = v.Name
		node["rabbitmq_version"] = v.RabbitMqVersion
		node["running"] = v.Running
		nodes[k] = node
	}

	if err = d.Set("nodes", nodes); err != nil {
		return diag.Errorf("error setting nodes for resource %s, %s", d.Id(), err)
	}

	return diag.Diagnostics{}
}

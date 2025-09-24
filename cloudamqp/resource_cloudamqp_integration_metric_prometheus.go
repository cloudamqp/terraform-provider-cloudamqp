package cloudamqp

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIntegrationMetricPrometheus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationMetricPrometheusCreate,
		ReadContext:   resourceIntegrationMetricPrometheusRead,
		UpdateContext: resourceIntegrationMetricPrometheusUpdate,
		DeleteContext: resourceIntegrationMetricPrometheusDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Instance identifier",
			},
			"newrelic_v3": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tags": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tags. E.g. env=prod,region=europe",
						},
					},
				},
			},
		},
	}
}

func resourceIntegrationMetricPrometheusCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api     = meta.(*api.API)
		intName string
		params  = make(map[string]any)
	)

	if d.Get("newrelic_v3") != nil {
		intName = "newrelic_v3"
		params["api_key"] = d.Get("newrelic_v3").(*schema.Set).List()[0].(map[string]any)["api_key"]
		params["tags"] = d.Get("newrelic_v3").(*schema.Set).List()[0].(map[string]any)["tags"]
	}

	data, err := api.CreateIntegration(ctx, d.Get("instance_id").(int), "metrics", intName, params)
	if err != nil {
		return diag.FromErr(err)
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
	}

	return nil
}

func resourceIntegrationMetricPrometheusRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if strings.Contains(d.Id(), ",") {
		tflog.Info(ctx, fmt.Sprintf("import resource with identifier: %s", d.Id()))
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		instanceID, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
	}
	if d.Get("instance_id").(int) == 0 {
		return diag.Errorf("missing instance identifier: {resource_id},{instance_id}")
	}

	api := meta.(*api.API)
	data, err := api.ReadIntegration(ctx, d.Get("instance_id").(int), "metrics", d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	// Handle resource drift and trigger re-creation if resource been deleted outside the provider
	if data == nil {
		d.SetId("")
		return nil
	}

	name := strings.ToLower(data["type"].(string))
	d.Set("name", name)
	if name == "newrelic_v3" {
		if _, ok := data["api_key"]; ok {
			newRelicV3 := []map[string]any{
				{
					"api_key": data["api_key"],
				},
			}
			if tags, ok := data["tags"]; ok {
				newRelicV3[0]["tags"] = tags
			}
			if err := d.Set("newrelic_v3", newRelicV3); err != nil {
				return diag.Errorf("error setting newrelic_v3 for resource %s: %s", d.Id(), err)
			}
		}
	}

	return nil
}

func resourceIntegrationMetricPrometheusUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api    = meta.(*api.API)
		params = make(map[string]any)
	)

	if d.Get("newrelic_v3") != nil {
		params["api_key"] = d.Get("newrelic_v3").(*schema.Set).List()[0].(map[string]any)["api_key"]
		params["tags"] = d.Get("newrelic_v3").(*schema.Set).List()[0].(map[string]any)["tags"]
	}

	err := api.UpdateIntegration(ctx, d.Get("instance_id").(int), "metrics", d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceIntegrationMetricPrometheusDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	if err := api.DeleteIntegration(ctx, d.Get("instance_id").(int), "metrics", d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"datadog_v3", "azure_monitor"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"tags": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tags. E.g. env=prod,service=web",
						},
					},
				},
			},
			"datadog_v3": {
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"newrelic_v3", "azure_monitor"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"region": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Datadog region; us1, us3, us5, or eu1",
							ValidateFunc: validation.StringInSlice([]string{"us1", "us3", "us5", "eu1"}, true),
						},
						"tags": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tags. E.g. env=prod,service=web",
						},
					},
				},
			},
			"azure_monitor": {
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"newrelic_v3", "datadog_v3"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_string": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Azure Application Insights Connection String",
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

	// Check which integration type is configured
	if newrelicList := d.Get("newrelic_v3").(*schema.Set).List(); len(newrelicList) > 0 {
		intName = "newrelic_v3"
		newrelicConfig := newrelicList[0].(map[string]any)
		params["api_key"] = newrelicConfig["api_key"]
		if tags := newrelicConfig["tags"]; tags != nil && tags != "" {
			params["tags"] = tags
		}
	} else if datadogList := d.Get("datadog_v3").(*schema.Set).List(); len(datadogList) > 0 {
		intName = "datadog_v3"
		datadogConfig := datadogList[0].(map[string]any)
		params["api_key"] = datadogConfig["api_key"]
		if region := datadogConfig["region"]; region != nil && region != "" {
			params["region"] = region
		}
		if tags := datadogConfig["tags"]; tags != nil && tags != "" {
			params["tags"] = tags
		}
	} else if azureMonitorList := d.Get("azure_monitor").(*schema.Set).List(); len(azureMonitorList) > 0 {
		intName = "azure_monitor"
		azureMonitorConfig := azureMonitorList[0].(map[string]any)
		params["connection_string"] = azureMonitorConfig["connection_string"]
	}

	if intName == "" {
		return diag.Errorf("no integration configuration provided")
	}

	data, err := api.CreateIntegration(ctx, d.Get("instance_id").(int), "metrics", intName, params)
	if err != nil {
		return diag.FromErr(err)
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
	}

	return resourceIntegrationMetricPrometheusRead(ctx, d, meta)
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

	// Clear all blocks first
	d.Set("newrelic_v3", nil)
	d.Set("datadog_v3", nil)

	name := strings.ToLower(data["type"].(string))
	if name == "newrelic_v3" {
		newRelicV3 := []map[string]any{{}}
		if _, ok := data["api_key"]; ok {
			newRelicV3[0]["api_key"] = data["api_key"]
		}
		if tags, ok := data["tags"]; ok {
			newRelicV3[0]["tags"] = tags
		}
		if err := d.Set("newrelic_v3", newRelicV3); err != nil {
			return diag.Errorf("error setting newrelic_v3 for resource %s: %s", d.Id(), err)
		}
	} else if name == "datadog_v3" {
		datadogV3 := []map[string]any{{}}
		if _, ok := data["api_key"]; ok {
			datadogV3[0]["api_key"] = data["api_key"]
		}
		if region, ok := data["region"]; ok {
			datadogV3[0]["region"] = region
		}
		if tags, ok := data["tags"]; ok {
			datadogV3[0]["tags"] = tags
		}
		if err := d.Set("datadog_v3", datadogV3); err != nil {
			return diag.Errorf("error setting datadog_v3 for resource %s: %s", d.Id(), err)
		}
	} else if name == "azure_monitor" {
		azureMonitor := []map[string]any{{}}
		if _, ok := data["connection_string"]; ok {
			azureMonitor[0]["connection_string"] = data["connection_string"]
		}
		if err := d.Set("azure_monitor", azureMonitor); err != nil {
			return diag.Errorf("error setting azure_monitor for resource %s: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceIntegrationMetricPrometheusUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api    = meta.(*api.API)
		params = make(map[string]any)
	)

	// Check which integration type is configured
	if newrelicList := d.Get("newrelic_v3").(*schema.Set).List(); len(newrelicList) > 0 {
		newrelicConfig := newrelicList[0].(map[string]any)
		params["api_key"] = newrelicConfig["api_key"]
		if tags := newrelicConfig["tags"]; tags != nil && tags != "" {
			params["tags"] = tags
		}
	} else if datadogList := d.Get("datadog_v3").(*schema.Set).List(); len(datadogList) > 0 {
		datadogConfig := datadogList[0].(map[string]any)
		params["api_key"] = datadogConfig["api_key"]
		if region := datadogConfig["region"]; region != nil && region != "" {
			params["region"] = region
		}
		if tags := datadogConfig["tags"]; tags != nil && tags != "" {
			params["tags"] = tags
		}
	} else if azureMonitorList := d.Get("azure_monitor").(*schema.Set).List(); len(azureMonitorList) > 0 {
		azureMonitorConfig := azureMonitorList[0].(map[string]any)
		params["connection_string"] = azureMonitorConfig["connection_string"]
	}

	err := api.UpdateIntegration(ctx, d.Get("instance_id").(int), "metrics", d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIntegrationMetricPrometheusRead(ctx, d, meta)
}

func resourceIntegrationMetricPrometheusDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	api := meta.(*api.API)
	if err := api.DeleteIntegration(ctx, d.Get("instance_id").(int), "metrics", d.Id()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

package cloudamqp

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
			"metrics_filter": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of metrics to include. If not specified, default metrics are used. See https://www.cloudamqp.com/docs/monitoring_metrics_splunk_v2.html#metrics-filtering for more information",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"newrelic_v3": {
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"datadog_v3", "azure_monitor", "splunk_v2", "dynatrace", "cloudwatch_v3", "stackdriver_v2"},
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
							Description:  "New Relic region; eu or us",
							ValidateFunc: validation.StringInSlice([]string{"eu", "us"}, true),
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
				ConflictsWith: []string{"newrelic_v3", "azure_monitor", "splunk_v2", "dynatrace", "cloudwatch_v3", "stackdriver_v2"},
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
							Description:  "Datadog region; us1, us3, us5, eu1, or ap2",
							ValidateFunc: validation.StringInSlice([]string{"us1", "us3", "us5", "eu1", "ap2"}, true),
						},
						"tags": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tags. E.g. env=prod,service=web",
						},
						"rabbitmq_dashboard_metrics_format": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable metric name transformation to match Datadog's RabbitMQ dashboard format",
						},
					},
				},
			},
			"azure_monitor": {
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"newrelic_v3", "datadog_v3", "splunk_v2", "dynatrace", "cloudwatch_v3", "stackdriver_v2"},
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
			"splunk_v2": {
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"newrelic_v3", "datadog_v3", "azure_monitor", "dynatrace", "cloudwatch_v3", "stackdriver_v2"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"token": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Splunk HEC token",
						},
						"endpoint": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Splunk HEC endpoint. E.g. https://your-instance-id.splunkcloud.com:8088/services/collector",
						},
						"tags": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tags. E.g. env=prod,service=web",
						},
					},
				},
			},
			"dynatrace": {
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"newrelic_v3", "datadog_v3", "azure_monitor", "splunk_v2", "cloudwatch_v3", "stackdriver_v2"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"environment_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Dynatrace environment ID",
						},
						"access_token": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Dynatrace access token with 'Ingest metrics' permission",
						},
						"tags": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tags. E.g. env=prod,service=web",
						},
					},
				},
			},
			"cloudwatch_v3": {
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"newrelic_v3", "datadog_v3", "azure_monitor", "splunk_v2", "dynatrace", "stackdriver_v2"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iam_role": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AWS IAM role ARN with PutMetricData permission",
						},
						"iam_external_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "External identifier that matches the role you created.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "AWS region",
						},
						"tags": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tags. E.g. env=prod,service=web",
						},
					},
				},
			},
			"stackdriver_v2": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"newrelic_v3", "datadog_v3", "azure_monitor", "splunk_v2", "dynatrace", "cloudwatch_v3"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credentials_file": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: "Base64-encoded Google service account key JSON file",
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// Only suppress for existing resources
								if d.Id() == "" {
									return false
								}

								if stackdriver := d.Get("stackdriver_v2").([]any); len(stackdriver) > 0 {
									config := stackdriver[0].(map[string]any)
									newCredentials, err := extractStackdriverCredentials(new)
									if err != nil {
										return false
									}
									// Suppress diff if new credentials match current state
									return newCredentials["project_id"] == config["project_id"] &&
										newCredentials["client_email"] == config["client_email"] &&
										newCredentials["private_key_id"] == config["private_key_id"] &&
										newCredentials["private_key"] == config["private_key"]
								}
								return false
							},
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Google Cloud project ID (computed from credentials file)",
						},
						"client_email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Google service account client email (computed from credentials file)",
						},
						"private_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Google service account private key (computed from credentials file)",
						},
						"private_key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Google service account private key ID (computed from credentials file)",
						},
						"tags": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tags. E.g. env=prod,service=web",
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

	if newrelicList := d.Get("newrelic_v3").(*schema.Set).List(); len(newrelicList) > 0 {
		intName = "newrelic_v3"
		newrelicConfig := newrelicList[0].(map[string]any)
		params["api_key"] = newrelicConfig["api_key"]
		if region := newrelicConfig["region"]; region != nil && region != "" {
			params["region"] = region
		}
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
		if format, ok := datadogConfig["rabbitmq_dashboard_metrics_format"].(bool); ok {
			if format {
				params["rabbitmq_dashboard_metrics_format"] = "true"
			} else {
				params["rabbitmq_dashboard_metrics_format"] = "false"
			}
		}
	} else if azureMonitorList := d.Get("azure_monitor").(*schema.Set).List(); len(azureMonitorList) > 0 {
		intName = "azure_monitor"
		azureMonitorConfig := azureMonitorList[0].(map[string]any)
		params["connection_string"] = azureMonitorConfig["connection_string"]
	} else if splunkList := d.Get("splunk_v2").(*schema.Set).List(); len(splunkList) > 0 {
		intName = "splunk_v2"
		splunkConfig := splunkList[0].(map[string]any)
		params["token"] = splunkConfig["token"]
		params["endpoint"] = splunkConfig["endpoint"]
		if tags := splunkConfig["tags"]; tags != nil && tags != "" {
			params["tags"] = tags
		}
	} else if dynatraceList := d.Get("dynatrace").(*schema.Set).List(); len(dynatraceList) > 0 {
		intName = "dynatrace"
		dynatraceConfig := dynatraceList[0].(map[string]any)
		params["environment_id"] = dynatraceConfig["environment_id"]
		params["access_token"] = dynatraceConfig["access_token"]
		if tags := dynatraceConfig["tags"]; tags != nil && tags != "" {
			params["tags"] = tags
		}
	} else if cloudwatchList := d.Get("cloudwatch_v3").(*schema.Set).List(); len(cloudwatchList) > 0 {
		intName = "cloudwatch_v3"
		cloudwatchConfig := cloudwatchList[0].(map[string]any)
		params["iam_role"] = cloudwatchConfig["iam_role"]
		params["iam_external_id"] = cloudwatchConfig["iam_external_id"]
		params["region"] = cloudwatchConfig["region"]
		if tags := cloudwatchConfig["tags"]; tags != nil && tags != "" {
			params["tags"] = tags
		}
	} else if stackdriverList := d.Get("stackdriver_v2").([]any); len(stackdriverList) > 0 {
		intName = "stackdriver_v2"
		stackdriverConfig := stackdriverList[0].(map[string]any)
		credentials := stackdriverConfig["credentials_file"].(string)

		extractedCredentials, err := extractStackdriverCredentials(credentials)
		if err != nil {
			return diag.FromErr(err)
		}

		for key, value := range extractedCredentials {
			params[key] = value
		}

		if tags := stackdriverConfig["tags"]; tags != nil && tags != "" {
			params["tags"] = tags
		}
	}

	if intName == "" {
		return diag.Errorf("no integration configuration provided")
	}

	if metricsFilter := d.Get("metrics_filter").([]any); len(metricsFilter) > 0 {
		filters := make([]string, len(metricsFilter))
		for i, v := range metricsFilter {
			filters[i] = v.(string)
		}
		params["metrics_filter"] = filters
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

func extractStackdriverCredentials(credentials string) (map[string]string, error) {
	decoded, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		return nil, fmt.Errorf("failed to decode stackdriver credentials: %s", err)
	}

	var jsonMap map[string]any
	if err := json.Unmarshal(decoded, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to parse stackdriver credentials JSON: %s", err)
	}

	requiredFields := []string{"client_email", "private_key_id", "private_key", "project_id"}
	for _, field := range requiredFields {
		if jsonMap[field] == nil || jsonMap[field] == "" {
			return nil, fmt.Errorf("required field '%s' is missing from credentials JSON", field)
		}
	}

	return map[string]string{
		"client_email":   jsonMap["client_email"].(string),
		"private_key_id": jsonMap["private_key_id"].(string),
		"private_key":    jsonMap["private_key"].(string),
		"project_id":     jsonMap["project_id"].(string),
	}, nil
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

	// Resource drift: instance or resource not found, trigger re-creation
	if data == nil {
		tflog.Info(ctx, fmt.Sprintf("prometheus metric integration not found, resource will be recreated: %s", d.Id()))
		d.SetId("")
		return nil
	}

	d.Set("newrelic_v3", nil)
	d.Set("datadog_v3", nil)
	d.Set("azure_monitor", nil)
	d.Set("splunk_v2", nil)
	d.Set("dynatrace", nil)
	d.Set("cloudwatch_v3", nil)
	d.Set("stackdriver_v2", nil)

	if metricsFilter, ok := data["metrics_filter"]; ok && metricsFilter != nil {
		if filterSlice, ok := metricsFilter.([]any); ok {
			d.Set("metrics_filter", filterSlice)
		}
	}

	name := strings.ToLower(data["type"].(string))
	if name == "newrelic_v3" {
		newRelicV3 := []map[string]any{{}}
		if _, ok := data["api_key"]; ok {
			newRelicV3[0]["api_key"] = data["api_key"]
		}
		if region, ok := data["region"]; ok {
			newRelicV3[0]["region"] = region
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
		if format, ok := data["rabbitmq_dashboard_metrics_format"]; ok {
			datadogV3[0]["rabbitmq_dashboard_metrics_format"] = format == "true"
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
	} else if name == "splunk_v2" {
		splunkV2 := []map[string]any{{}}
		if _, ok := data["token"]; ok {
			splunkV2[0]["token"] = data["token"]
		}
		if _, ok := data["endpoint"]; ok {
			splunkV2[0]["endpoint"] = data["endpoint"]
		}
		if tags, ok := data["tags"]; ok {
			splunkV2[0]["tags"] = tags
		}
		if err := d.Set("splunk_v2", splunkV2); err != nil {
			return diag.Errorf("error setting splunk_v2 for resource %s: %s", d.Id(), err)
		}
	} else if name == "dynatrace" {
		dynatrace := []map[string]any{{}}
		if _, ok := data["environment_id"]; ok {
			dynatrace[0]["environment_id"] = data["environment_id"]
		}
		if _, ok := data["access_token"]; ok {
			dynatrace[0]["access_token"] = data["access_token"]
		}
		if tags, ok := data["tags"]; ok {
			dynatrace[0]["tags"] = tags
		}
		if err := d.Set("dynatrace", dynatrace); err != nil {
			return diag.Errorf("error setting dynatrace for resource %s: %s", d.Id(), err)
		}
	} else if name == "cloudwatch_v3" {
		cloudwatchV3 := []map[string]any{{}}
		if _, ok := data["iam_role"]; ok {
			cloudwatchV3[0]["iam_role"] = data["iam_role"]
		}
		if _, ok := data["iam_external_id"]; ok {
			cloudwatchV3[0]["iam_external_id"] = data["iam_external_id"]
		}
		if _, ok := data["region"]; ok {
			cloudwatchV3[0]["region"] = data["region"]
		}
		if tags, ok := data["tags"]; ok {
			cloudwatchV3[0]["tags"] = tags
		}
		if err := d.Set("cloudwatch_v3", cloudwatchV3); err != nil {
			return diag.Errorf("error setting cloudwatch_v3 for resource %s: %s", d.Id(), err)
		}
	} else if name == "stackdriver_v2" {
		stackdriverV2 := []map[string]any{{}}

		if project_id, ok := data["project_id"]; ok {
			stackdriverV2[0]["project_id"] = project_id
		}
		if client_email, ok := data["client_email"]; ok {
			stackdriverV2[0]["client_email"] = client_email
		}
		if private_key, ok := data["private_key"]; ok {
			stackdriverV2[0]["private_key"] = private_key
		}
		if private_key_id, ok := data["private_key_id"]; ok {
			stackdriverV2[0]["private_key_id"] = private_key_id
		}
		if tags, ok := data["tags"]; ok {
			stackdriverV2[0]["tags"] = tags
		}

		if err := d.Set("stackdriver_v2", stackdriverV2); err != nil {
			return diag.Errorf("error setting stackdriver_v2 for resource %s: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceIntegrationMetricPrometheusUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var (
		api    = meta.(*api.API)
		params = make(map[string]any)
	)

	if newrelicList := d.Get("newrelic_v3").(*schema.Set).List(); len(newrelicList) > 0 {
		newrelicConfig := newrelicList[0].(map[string]any)
		params["api_key"] = newrelicConfig["api_key"]
		if region := newrelicConfig["region"]; region != nil && region != "" {
			params["region"] = region
		}
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
		if format, ok := datadogConfig["rabbitmq_dashboard_metrics_format"].(bool); ok {
			if format {
				params["rabbitmq_dashboard_metrics_format"] = "true"
			} else {
				params["rabbitmq_dashboard_metrics_format"] = "false"
			}
		}
	} else if azureMonitorList := d.Get("azure_monitor").(*schema.Set).List(); len(azureMonitorList) > 0 {
		azureMonitorConfig := azureMonitorList[0].(map[string]any)
		params["connection_string"] = azureMonitorConfig["connection_string"]
	} else if splunkList := d.Get("splunk_v2").(*schema.Set).List(); len(splunkList) > 0 {
		splunkConfig := splunkList[0].(map[string]any)
		params["token"] = splunkConfig["token"]
		params["endpoint"] = splunkConfig["endpoint"]
		if tags := splunkConfig["tags"]; tags != nil && tags != "" {
			params["tags"] = tags
		}
	} else if dynatraceList := d.Get("dynatrace").(*schema.Set).List(); len(dynatraceList) > 0 {
		dynatraceConfig := dynatraceList[0].(map[string]any)
		params["environment_id"] = dynatraceConfig["environment_id"]
		params["access_token"] = dynatraceConfig["access_token"]
		if tags := dynatraceConfig["tags"]; tags != nil && tags != "" {
			params["tags"] = tags
		}
	} else if cloudwatchList := d.Get("cloudwatch_v3").(*schema.Set).List(); len(cloudwatchList) > 0 {
		cloudwatchConfig := cloudwatchList[0].(map[string]any)
		params["iam_role"] = cloudwatchConfig["iam_role"]
		params["iam_external_id"] = cloudwatchConfig["iam_external_id"]
		params["region"] = cloudwatchConfig["region"]
		if tags := cloudwatchConfig["tags"]; tags != nil && tags != "" {
			params["tags"] = tags
		}
	} else if stackdriverList := d.Get("stackdriver_v2").([]interface{}); len(stackdriverList) > 0 {
		stackdriverConfig := stackdriverList[0].(map[string]any)

		credentials := stackdriverConfig["credentials_file"].(string)
		extractedCreds, err := extractStackdriverCredentials(credentials)
		if err != nil {
			return diag.FromErr(err)
		}

		for key, value := range extractedCreds {
			params[key] = value
		}

		if tags := stackdriverConfig["tags"]; tags != nil && tags != "" {
			params["tags"] = tags
		}
	}

	if d.HasChange("metrics_filter") {
		metricsFilter := d.Get("metrics_filter").([]any)
		filters := make([]string, len(metricsFilter))
		for i, v := range metricsFilter {
			filters[i] = v.(string)
		}
		params["metrics_filter"] = filters
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

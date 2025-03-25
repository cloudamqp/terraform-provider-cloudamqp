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

func resourceIntegrationLog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIntegrationLogCreate,
		ReadContext:   resourceIntegrationLogRead,
		UpdateContext: resourceIntegrationLogUpdate,
		DeleteContext: resourceIntegrationLogDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Instance identifier used to make proxy calls",
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The name of log integration",
				ValidateDiagFunc: validateIntegrationLogName(),
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL to push the logs to. (Papertrail)",
			},
			"host_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Destination to send the logs. (Splunk)",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The token used for authentication. (Loggly, Logentries, Splunk, Scalyr)",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region hosting integration service. (Cloudwatch, Datadog)",
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS access key identifier. (Cloudwatch)",
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS secret access key. (Cloudwatch)",
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The API key for the integration service. (Datadog)",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "(optional) tags. E.g. env=prod,region=europe. (Datadog)",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Project ID. (Stackdriver)",
			},
			"private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "The private API key used for authentication. (Stackdriver, Coralogix)",
			},
			"client_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The client email. (Stackdriver)",
			},
			"host": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The host information. (Scalyr)",
				ValidateDiagFunc: validateIntegrationLogScalyrHost(),
			},
			"sourcetype": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Assign source type to the data exported, eg. generic_single_line. (Splunk)",
			},
			"private_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "Private key identifier. (Stackdriver)",
			},
			"credentials": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Base64Encoded credentials. (Stackdriver)",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The syslog destination to send the logs to. (Coralogix)",
			},
			"application": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The 'Application Name'. (Coralogix)",
			},
			"subsystem": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The 'Subsystem Name'. (Coralogix)",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The tenant ID. (Azure Monitor)",
			},
			"application_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The application ID. (Azure Monitor)",
			},
			"application_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The application secret. (Azure Monitor)",
			},
			"dce_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The data collection endpoint. (Azure Monitor)",
			},
			"table": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The table name. (Azure Monitor)",
			},
			"dcr_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of data collection rule that your DCE is linked to. (Azure Monitor)",
			},
		},
	}
}

func resourceIntegrationLogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		api     = meta.(*api.API)
		intName = strings.ToLower(d.Get("name").(string))
		keys    = integrationLogKeys(intName)
		params  = make(map[string]interface{})
	)

	v := d.Get("credentials")
	if intName == "stackdriver" && v != "" {
		uDec, err := base64.URLEncoding.DecodeString(v.(string))
		if err != nil {
			return diag.Errorf("log integration failed, error decoding private_key: %s ", err.Error())
		}
		var jsonMap map[string]interface{}
		json.Unmarshal([]byte(uDec), &jsonMap)
		fmt.Printf("jsonMap: %v", jsonMap)
		for _, k := range keys {
			params[k] = jsonMap[k]
		}
	} else {
		for _, k := range keys {
			v := d.Get(k)
			if k == "tags" && v == "" {
				delete(params, k)
			} else if v != nil {
				params[k] = v
			}
		}
	}

	data, err := api.CreateIntegration(ctx, d.Get("instance_id").(int), "logs", intName, params)

	if err != nil {
		return diag.FromErr(err)
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
	}

	return resourceIntegrationLogRead(ctx, d, meta)
}

func resourceIntegrationLogRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	data, err := api.ReadIntegration(ctx, d.Get("instance_id").(int), "logs", d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range data {
		if k == "type" {
			d.Set("name", v)
		}
		if validateIntegrationLogsSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return diag.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return nil
}

func resourceIntegrationLogUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		intName = strings.ToLower(d.Get("name").(string))
		keys    = integrationLogKeys(intName)
		params  = make(map[string]interface{})
	)

	v := d.Get("credentials")
	if intName == "stackdriver" && v != "" {
		uDec, err := base64.URLEncoding.DecodeString(v.(string))
		if err != nil {
			return diag.Errorf("log integration failed, error decoding private_key: %s ", err.Error())
		}
		var jsonMap map[string]interface{}
		json.Unmarshal([]byte(uDec), &jsonMap)
		for _, k := range keys {
			params[k] = jsonMap[k]
		}
	} else {
		for _, k := range keys {
			v := d.Get(k)
			if k == "tags" && v == "" {
				delete(params, k)
			} else if v != nil {
				params[k] = v
			}
		}
	}

	api := meta.(*api.API)
	err := api.UpdateIntegration(ctx, d.Get("instance_id").(int), "logs", d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceIntegrationLogRead(ctx, d, meta)
}

func resourceIntegrationLogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	api := meta.(*api.API)
	if err := api.DeleteIntegration(ctx, d.Get("instance_id").(int), "logs", d.Id()); err != nil {
		diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func validateIntegrationLogName() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"azure_monitor",
		"cloudwatchlog",
		"coralogix",
		"datadog",
		"logentries",
		"loggly",
		"papertrail",
		"scalyr",
		"splunk",
		"stackdriver",
	}, true))
}

func validateIntegrationLogScalyrHost() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice([]string{
		"app.scalyr.com",
		"app.eu.scalyr.com",
	}, true))
}

func validateIntegrationLogsSchemaAttribute(key string) bool {
	switch key {
	case "access_key_id",
		"api_key",
		"application",
		"application_id",
		"application_secret",
		"client_email",
		"dce_uri",
		"dcr_id",
		"endpoint",
		"host",
		"host_port",
		"private_key",
		"private_key_id",
		"project_id",
		"region",
		"secret_access_key",
		"sourcetype",
		"subsystem",
		"table",
		"tags",
		"tenant_id",
		"token",
		"url":
		return true
	}
	return false
}

func integrationLogKeys(intName string) []string {
	switch intName {
	case "azure_monitor":
		return []string{"tenant_id", "application_id", "application_secret", "dce_uri", "table", "dcr_id"}
	case "cloudwatchlog":
		return []string{"region", "access_key_id", "secret_access_key"}
	case "coralogix":
		return []string{"private_key", "endpoint", "application", "subsystem"}
	case "datadog":
		return []string{"region", "api_key", "tags"}
	case "logentries":
		return []string{"token"}
	case "loggly":
		return []string{"token"}
	case "papertrail":
		return []string{"url"}
	case "scalyr":
		return []string{"token", "host"}
	case "splunk":
		return []string{"host_port", "token", "sourcetype"}
	case "stackdriver":
		return []string{"client_email", "private_key_id", "private_key", "project_id"}
	default:
		return []string{"url", "host_port", "token", "region", "access_key_id", "secret_access_key"}
	}
}

package cloudamqp

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceIntegrationLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceIntegrationLogCreate,
		Read:   resourceIntegrationLogRead,
		Update: resourceIntegrationLogUpdate,
		Delete: resourceIntegrationLogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier used to make proxy calls",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of log integration",
				ValidateFunc: validateIntegrationLogName(),
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
				Description: "The private key. (Stackdriver)",
			},
			"client_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The client email. (Stackdriver)",
			},
			"host": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The host information. (Scalyr)",
				ValidateFunc: validateIntegrationLogScalyrHost(),
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
		},
	}
}

func resourceIntegrationLogCreate(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("Log integration failed, error decoding private_key: %s ", err.Error())
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

	data, err := api.CreateIntegration(d.Get("instance_id").(int), "logs", intName, params)

	if err != nil {
		return err
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
	}

	return resourceIntegrationLogRead(d, meta)
}

func resourceIntegrationLogRead(d *schema.ResourceData, meta interface{}) error {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		instanceID, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
	}
	if d.Get("instance_id").(int) == 0 {
		return errors.New("Missing instance identifier: {resource_id},{instance_id}")
	}

	api := meta.(*api.API)
	data, err := api.ReadIntegration(d.Get("instance_id").(int), "logs", d.Id())

	if err != nil {
		return err
	}

	for k, v := range data {
		if k == "type" {
			d.Set("name", v)
		}
		if validateIntegrationLogsSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return nil
}

func resourceIntegrationLogUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		intName = strings.ToLower(d.Get("name").(string))
		keys    = integrationLogKeys(intName)
		params  = make(map[string]interface{})
	)

	v := d.Get("credentials")
	if intName == "stackdriver" && v != "" {
		uDec, err := base64.URLEncoding.DecodeString(v.(string))
		if err != nil {
			return fmt.Errorf("Log integration failed, error decoding private_key: %s ", err.Error())
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
	err := api.UpdateIntegration(d.Get("instance_id").(int), "logs", d.Id(), params)
	if err != nil {
		return err
	}
	return resourceIntegrationLogRead(d, meta)
}

func resourceIntegrationLogDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	return api.DeleteIntegration(d.Get("instance_id").(int), "logs", d.Id())
}

func validateIntegrationLogName() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"papertrail",
		"loggly",
		"logentries",
		"splunk",
		"cloudwatchlog",
		"datadog",
		"stackdriver",
		"scalyr",
	}, true)
}

func validateIntegrationLogScalyrHost() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"app.scalyr.com",
		"app.eu.scalyr.com",
	}, true)
}

func validateIntegrationLogsSchemaAttribute(key string) bool {
	switch key {
	case "url",
		"host_port",
		"token",
		"region",
		"access_key_id",
		"secret_access_key",
		"api_key",
		"tags",
		"project_id",
		"client_email",
		"private_key",
		"private_key_id",
		"host",
		"sourcetype":
		return true
	}
	return false
}

func integrationLogKeys(intName string) []string {
	switch intName {
	case "papertrail":
		return []string{"url"}
	case "loggly":
		return []string{"token"}
	case "logentries":
		return []string{"token"}
	case "splunk":
		return []string{"host_port", "token", "sourcetype"}
	case "cloudwatchlog":
		return []string{"region", "access_key_id", "secret_access_key"}
	case "datadog":
		return []string{"region", "api_key", "tags"}
	case "stackdriver":
		return []string{"client_email", "private_key_id", "private_key", "project_id"}
	case "scalyr":
		return []string{"token", "host"}
	default:
		return []string{"url", "host_port", "token", "region", "access_key_id", "secret_access_key"}
	}
}

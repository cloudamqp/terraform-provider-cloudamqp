package cloudamqp

import (
	"errors"
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
				Description: "The token used for authentication. (Loggly, Logentries, Splunk)",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region hosting integration service. (Cloudwatch)",
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
		},
	}
}

func resourceIntegrationLogCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := integrationLogKeys(d.Get("name").(string))
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	data, err := api.CreateIntegration(d.Get("instance_id").(int), "logs", d.Get("name").(string), params)

	if err != nil {
		return err
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
	}

	for k, v := range data {
		if validateIntegrationLogsSchemaAttribute(k) {
			d.Set(k, v)
		}
	}
	return nil
}

func resourceIntegrationLogRead(d *schema.ResourceData, meta interface{}) error {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(s[0])
		instance_id, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instance_id)
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
		if validateIntegrationLogsSchemaAttribute(k) {
			d.Set(k, v)
		}
	}
	return nil
}

func resourceIntegrationLogUpdate(d *schema.ResourceData, meta interface{}) error {
	keys := integrationLogKeys(d.Get("name").(string))
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	api := meta.(*api.API)
	err := api.UpdateIntegration(d.Get("instance_id").(int), "logs", d.Id(), params)
	return err
}

func resourceIntegrationLogDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	params := make(map[string]interface{})
	params["id"] = d.Id()
	return api.DeleteIntegration(d.Get("instance_id").(int), "logs", d.Id())
}

func validateIntegrationLogName() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"papertrail",
		"loggly",
		"logentries",
		"splunk",
		"cloudwatchlog",
	}, true)
}

func validateIntegrationLogsSchemaAttribute(key string) bool {
	switch key {
	case "type",
		"url",
		"host_port",
		"token",
		"region",
		"access_key_id",
		"secret_access_key":
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
		return []string{"host_port", "token"}
	case "cloudwatchlog":
		return []string{"region", "access_key_id", "secret_access_key"}
	default:
		return []string{"url", "host_port", "token", "region", "access_key_id", "secret_access_key"}
	}
}

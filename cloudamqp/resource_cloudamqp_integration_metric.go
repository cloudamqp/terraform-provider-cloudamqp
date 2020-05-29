package cloudamqp

import (
	"errors"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceIntegrationMetric() *schema.Resource {
	return &schema.Resource{
		Create: resourceIntegrationMetricCreate,
		Read:   resourceIntegrationMetricRead,
		Update: resourceIntegrationMetricUpdate,
		Delete: resourceIntegrationMetricDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Instance identifier",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of metrics integration",
				ValidateFunc: validateIntegrationMetricName(),
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS region for Cloudwatch and [US/EU] for Data dog/New relic. (Cloudwatch, Data Dog, New Relic)",
			},
			"access_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AWS access key identifier. (Cloudwatch)",
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "AWS secret key. (Cloudwatch)",
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The API key for the integration service. (Librato)",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email address registred for the integration service. (Librato)",
			},
			"license_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The license key registred for the integration service. (New Relic)",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "(optional) tags. E.g. env=prod,region=europe",
			},
			"queue_whitelist": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "(optional) whitelist using regular expression",
			},
			"vhost_whitelist": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "(optional) whitelist using regular expression",
			},
		},
	}
}

func resourceIntegrationMetricCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := integrationMetricKeys(d.Get("name").(string))
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	data, err := api.CreateIntegration(d.Get("instance_id").(int), "metrics", d.Get("name").(string), params)

	if err != nil {
		return err
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
	}

	for k, v := range data {
		if validateIntegrationMetricSchemaAttribute(k) {
			d.Set(k, v)
		}
	}
	return nil
}

func resourceIntegrationMetricRead(d *schema.ResourceData, meta interface{}) error {
	if d.Get("instance_id").(int) == 0 {
		return errors.New("Missing instance identifier: {instance_id}")
	}
	if len(d.Get("name").(string)) == 0 {
		return errors.New("Missing type representation: {name}")
	}

	api := meta.(*api.API)
	data, err := api.ReadIntegration(d.Get("instance_id").(int), "metrics", d.Id())

	if err != nil {
		return err
	}

	for k, v := range data {
		if validateIntegrationMetricSchemaAttribute(k) {
			d.Set(k, v)
		}
	}
	return nil
}

func resourceIntegrationMetricUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := integrationMetricKeys(d.Get("name").(string))
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}
	err := api.UpdateIntegration(d.Get("instance_id").(int), "metrics", d.Id(), params)
	if err != nil {
		return err
	}
	return resourceIntegrationMetricRead(d, meta)
}

func resourceIntegrationMetricDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	params := make(map[string]interface{})
	params["id"] = d.Id()
	return api.DeleteIntegration(d.Get("instance_id").(int), "metrics", d.Id())
}

func validateIntegrationMetricName() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"cloudwatch",
		"cloudwatch_v2",
		"librato",
		"datadog",
		"datadog_v2",
		"newrelic",
		"newrelic_v2",
		"stackdriver",
	}, true)
}

func validateIntegrationMetricSchemaAttribute(key string) bool {
	switch key {
	case "region",
		"access_key_id",
		"secret_access_key",
		"tags",
		"queue_whitelist",
		"vhost_whitelist",
		"api_key",
		"email",
		"license_key":
		return true
	default:
		return false
	}
}

func integrationMetricKeys(intName string) []string {
	keys := []string{"tags", "queue_whitelist", "vhost_whitelist"}
	switch intName {
	case "cloudwatch":
		return append(keys, "region", "access_key_id", "secret_access_key")
	case "cloudwatch_v2":
		return append(keys, "region", "access_key_id", "secret_access_key")
	case "librato":
		return append(keys, "email", "api_key")
	case "datadog":
		return append(keys, "api_key", "region")
	case "datadog_v2":
		return append(keys, "api_key", "region")
	case "newrelic":
		return append(keys, "license_key")
	case "newrelic_v2":
		return append(keys, "api_key", "region")
	default:
		return append(keys, "region", "access_keys", "secret_access_keys", "email", "api_key", "license_key")
	}
}

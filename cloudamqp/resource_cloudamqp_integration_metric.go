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
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "**Deprecated**",
				Deprecated:    "use queue_allowlist instead",
				ConflictsWith: []string{"queue_allowlist"},
			},
			"vhost_whitelist": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "**Deprecated**",
				Deprecated:    "use vhost_allowlist instead",
				ConflictsWith: []string{"vhost_allowlist"},
			},
			"queue_allowlist": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "(optional) allowlist using regular expression",
				ConflictsWith: []string{"queue_whitelist"},
			},
			"vhost_allowlist": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "(optional) allowlist using regular expression",
				ConflictsWith: []string{"vhost_whitelist"},
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

func resourceIntegrationMetricCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		intName    = strings.ToLower(d.Get("name").(string))
		commonKeys = []string{"tags", "queue_allowlist", "vhost_allowlist"}
		keys       = integrationMetricKeys(commonKeys, intName)
		params     = make(map[string]interface{})
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
			if contains(commonKeys, k) {
				if v := d.Get(k); v == "" || v == nil {
					delete(params, k)
				} else {
					params[k] = v
				}
			} else {
				params[k] = jsonMap[k]
			}
		}
	} else {
		for _, k := range keys {
			v := d.Get(k)
			if contains(commonKeys, k) && v == "" {
				delete(params, k)
			} else if v != nil {
				params[k] = v
			}
		}
	}

	data, err := api.CreateIntegration(d.Get("instance_id").(int), "metrics", intName, params)

	if err != nil {
		return err
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
	}

	return resourceIntegrationMetricRead(d, meta)
}

func resourceIntegrationMetricRead(d *schema.ResourceData, meta interface{}) error {
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
	data, err := api.ReadIntegration(d.Get("instance_id").(int), "metrics", d.Id())

	if err != nil {
		return err
	}

	for k, v := range data {
		if k == "type" {
			d.Set("name", v)
		}
		if validateIntegrationMetricSchemaAttribute(k) {
			if err = d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}
	return nil
}

func resourceIntegrationMetricUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		intName    = strings.ToLower(d.Get("name").(string))
		commonKeys = []string{"tags", "queue_allowlist", "vhost_allowlist"}
		keys       = integrationMetricKeys(commonKeys, intName)
		params     = make(map[string]interface{})
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
			if contains(commonKeys, k) {
				if v := d.Get(k); v == "" || v == nil {
					delete(params, k)
				} else {
					params[k] = v
				}
			} else {
				params[k] = jsonMap[k]
			}
		}
	} else {
		for _, k := range keys {
			v := d.Get(k)
			if contains(commonKeys, k) && v == "" {
				delete(params, k)
			} else if v != nil {
				params[k] = v
			}
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
		"queue_allowlist",
		"vhost_allowlist",
		"api_key",
		"email",
		"license_key",
		"project_id",
		"private_key",
		"client_email",
		"private_key_id":
		return true
	default:
		return false
	}
}

func integrationMetricKeys(commonKeys []string, intName string) []string {
	keys := commonKeys
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
	case "stackdriver":
		return append(keys, "client_email", "private_key_id", "private_key", "project_id")
	default:
		return append(keys, "region", "access_keys", "secret_access_keys", "email", "api_key", "license_key")
	}
}

func contains(s []string, searchString string) bool {
	for i := range s {
		if searchString == s[i] {
			return true
		}
	}
	return false
}

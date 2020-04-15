package cloudamqp

import (
	"errors"
	"log"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceIntegrationMetric() *schema.Resource {
	return &schema.Resource{
		Create: resourceIntegrationMetricCreate,
		Read:   resourceIntegrationMetricRead,
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
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The type of metrics integration",
				ValidateFunc: validateIntegrationMetricType(),
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CloudWatch/CloudWatch v2/Data dog(US/EU)/Data dog v2(US/EU)/New relic(US/EU) - region",
			},
			"aws_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CloudWatch/CloudWatch v2 - aws access key id",
			},
			"aws_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "CloudWatch/CloudWatch v2 - aws secret key",
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CloudWatch/CloudWatch v2/Librato/Data dog/Data dog v2/New relic(US/EU) - optional tags. E.g. env=prod,region=europe",
			},
			"queue_whitelist": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CloudWatch/CloudWatch v2/Librato/Data dog/Data dog v2/New relic(US/EU) - whitelist using regular expression",
			},
			"vhost_whitelist": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CloudWatch/CloudWatch v2/Librato/Data dog/Data dog v2/New relic(US/EU) - whitelist using regular expression",
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Librato - email",
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Librato/Data dog/Data dog v2/New relic(US/EU) - api access key",
			},
		},
	}
}

func resourceIntegrationMetricCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"region", "aws_key", "aws_secret", "tags", "queue_whitelist", "vhost_whitelist", "api_key", "email"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}
	log.Printf("[DEBUG] cloudamqp::resource::metric_integration::create params: %v", params)

	data, err := api.CreateIntegration(d.Get("instance_id").(int), "metrics", d.Get("type").(string), params)

	if err != nil {
		return err
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
		log.Printf("[DEBUG] cloudamqp::resource::metric_integration::create id set: %v", d.Id())
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
	if len(d.Get("type").(string)) == 0 {
		return errors.New("Missing type representation: {type}")
	}

	log.Printf("[DEBUG] cloudamqp::resource::metric_integration::read instance id: %v, id: %v", d.Get("instance_id"), d.Id())
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

func resourceIntegrationMetricDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	params := make(map[string]interface{})
	params["id"] = d.Id()
	log.Printf("[DEBUG] cloudamqp::resource::metric_integration::delete instance_id: %v, id: %s", d.Get("instance_id"), d.Id())
	return api.DeleteIntegration(d.Get("instance_id").(int), "metrics", d.Id())
}

func validateIntegrationMetricType() schema.SchemaValidateFunc {
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
		"aws_key",
		"aws_secret",
		"tags",
		"queue_whitelist",
		"vhost_whitelist",
		"api_key",
		"email":
		return true
	}
	return false
}

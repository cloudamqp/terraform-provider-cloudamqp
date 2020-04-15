package cloudamqp

import (
	"errors"
	"log"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceIntegrationLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceIntegrationLogCreate,
		Read:   resourceIntegrationLogRead,
		Delete: resourceIntegrationLogDelete,
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
				Description:  "The type of log integration",
				ValidateFunc: validateIntegrationLogType(),
			},
			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Papertrail/Splunk - address",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Loggly/Logentries/Splunk - token",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CloudWatch - region",
			},
			"aws_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CloudWatch - aws access key id",
			},
			"aws_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "CloudWatch - aws secret key",
			},
		},
	}
}

func resourceIntegrationLogCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"address", "token", "region", "aws_key", "aws_secrect"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}
	log.Printf("[DEBUG] cloudamqp::resource::log_integration::create params: %v", params)

	data, err := api.CreateIntegration(d.Get("instance_id").(int), "logs", d.Get("type").(string), params)

	if err != nil {
		return err
	}
	if data["id"] != nil {
		d.SetId(data["id"].(string))
		log.Printf("[DEBUG] cloudamqp::resource::log_integration::create id set: %v", d.Id())
	}

	for k, v := range data {
		if validateIntegrationLogsSchemaAttribute(k) {
			d.Set(k, v)
		}
	}
	return nil
}

func resourceIntegrationLogRead(d *schema.ResourceData, meta interface{}) error {
	if d.Get("instance_id").(int) == 0 {
		return errors.New("Missing instance identifier: {instance_id}")
	}
	if len(d.Get("type").(string)) == 0 {
		return errors.New("Missing type representation: {type}")
	}

	log.Printf("[DEBUG] cloudamqp::resource::log_integration::read instance id: %v, id: %v", d.Get("instance_id"), d.Id())
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

func resourceIntegrationLogDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	params := make(map[string]interface{})
	params["id"] = d.Id()
	log.Printf("[DEBUG] cloudamqp::resource::log_integration::delete instance_id: %v, type: %s, id: %s", d.Get("instance_id"), d.Get("type"), d.Id())
	return api.DeleteIntegration(d.Get("instance_id").(int), "logs", d.Id())
}

func validateIntegrationLogType() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"papertrail",
		"loggly",
		"logentries",
		"splunk",
		"cloudwatchlog",
		"stackdriver",
	}, true)
}

func validateIntegrationLogsSchemaAttribute(key string) bool {
	switch key {
	case "address",
		"token",
		"region",
		"aws_key",
		"aws_secret":
		return true
	}
	return false
}

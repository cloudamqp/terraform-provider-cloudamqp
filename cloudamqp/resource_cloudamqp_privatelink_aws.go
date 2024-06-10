package cloudamqp

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePrivateLinkAws() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrivateLinkAwsCreate,
		Read:   resourcePrivateLinkAwsRead,
		Update: resourcePrivateLinkAwsUpdate,
		Delete: resourcePrivateLinkAwsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The CloudAMQP instance identifier",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the PrivateLink [enabled, pending, disabled]",
			},
			"service_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service name of the PrivateLink, needed when creating the endpoint",
			},
			"allowed_principals": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Allowed principals to access the endpoint service",
			},
			"active_zones": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Covering availability zones used when creating an Endpoint from other VPC",
			},
			"sleep": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "Configurable sleep in seconds between retries when enable PrivateLink",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1800,
				Description: "Configurable timeout in seconds when enable PrivateLink",
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ValidateValue("allowed_principals", func(ctx context.Context, value, meta interface{}) error {
				for _, v := range value.([]interface{}) {
					re := regexp.MustCompile(`^arn:aws:iam::\d{12}:(root|user/.+|role/.+)$`)
					if !re.MatchString(v.(string)) {
						return fmt.Errorf("invalid ARN : %v", v)
					}
				}
				return nil
			}),
		),
	}
}

func resourcePrivateLinkAwsCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
		params     = make(map[string][]interface{})
	)

	params["allowed_principals"] = d.Get("allowed_principals").([]interface{})
	err := api.EnablePrivatelink(instanceID, params, sleep, timeout)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", instanceID))
	return resourcePrivateLinkAwsRead(d, meta)
}

func resourcePrivateLinkAwsRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api           = meta.(*api.API)
		instanceID, _ = strconv.Atoi(d.Id()) // Uses d.Id() to allow import
		sleep         = d.Get("sleep").(int)
		timeout       = d.Get("timeout").(int)
	)

	// Set arguments during import
	if d.Get("instance_id").(int) == 0 {
		d.Set("instance_id", instanceID)
	}
	if sleep == 0 && timeout == 0 {
		sleep = 10
		d.Set("sleep", 10)
		timeout = 1800
		d.Set("timeout", 1800)
	}

	data, err := api.ReadPrivatelink(instanceID, sleep, timeout)
	if err != nil {
		return err
	}

	for k, v := range data {
		if validatePrivateLinkAwsSchemaAttribute(k) {
			d.Set(k, v)
		}
	}
	return nil
}

func resourcePrivateLinkAwsUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		params     = make(map[string][]interface{})
	)

	params["allowed_principals"] = d.Get("allowed_principals").([]interface{})
	err := api.UpdatePrivatelink(instanceID, params)
	if err != nil {
		return err
	}
	return nil
}

func resourcePrivateLinkAwsDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
	)

	err := api.DisablePrivatelink(instanceID)
	if err != nil {
		return err
	}
	return nil
}

func validatePrivateLinkAwsSchemaAttribute(key string) bool {
	switch key {
	case "status",
		"service_name",
		"allowed_principals",
		"active_zones":
		return true
	}
	return false
}

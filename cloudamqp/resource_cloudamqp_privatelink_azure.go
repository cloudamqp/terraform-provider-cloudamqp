package cloudamqp

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePrivateLinkAzure() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrivateLinkAzureCreate,
		Read:   resourcePrivateLinkAzureRead,
		Update: resourcePrivateLinkAzureUpdate,
		Delete: resourcePrivateLinkAzureDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
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
				Description: "Service name (alias) of the PrivateLink, needed when creating the endpoint",
			},
			"server_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the server having the PrivateLink enabled",
			},
			"approved_subscriptions": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Approved subscriptions to access the endpoint service",
			},
			"sleep": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     60,
				Description: "Configurable sleep in seconds between retries when enable PrivateLink",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "Configurable timeout in seconds when enable PrivateLink",
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ValidateValue("approved_subscriptions", func(value, meta interface{}) error {
				for _, v := range value.([]interface{}) {
					re := regexp.MustCompile(`/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/`)
					if !re.MatchString(v.(string)) {
						return fmt.Errorf("Invalid ARN : %v", v)
					}
				}
				return nil
			}),
		),
	}
}

func resourcePrivateLinkAzureCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
		params     = make(map[string][]interface{})
	)

	err := api.EnablePrivatelink(instanceID, sleep, timeout)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", instanceID))
	params["approved_subscriptions"] = d.Get("approved_subscriptions").([]interface{})
	if len(params) > 0 {
		err := api.UpdatePrivatelink(instanceID, params)
		if err != nil {
			return err
		}
	}

	return resourcePrivateLinkAwsRead(d, meta)
}

func resourcePrivateLinkAzureRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api           = meta.(*api.API)
		instanceID, _ = strconv.Atoi(d.Id()) // Uses d.Id() to allow import
	)

	data, err := api.ReadPrivatelink(instanceID)
	if err != nil {
		return err
	}

	for k, v := range data {
		if validatePrivateLinkAzureSchemaAttribute(k) {
			if k == "alias" {
				d.Set("service_name", v)
			} else {
				d.Set(k, v)
			}
		}
	}
	return nil
}

func resourcePrivateLinkAzureUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		params     = make(map[string][]interface{})
	)

	params["approved_subscriptions"] = d.Get("approved_subscriptions").([]interface{})
	err := api.UpdatePrivatelink(instanceID, params)
	if err != nil {
		return err
	}
	return nil
}

func resourcePrivateLinkAzureDelete(d *schema.ResourceData, meta interface{}) error {
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

func validatePrivateLinkAzureSchemaAttribute(key string) bool {
	switch key {
	case "status",
		"server_name",
		"alias",
		"approved_subscriptions":
		return true
	}
	return false
}

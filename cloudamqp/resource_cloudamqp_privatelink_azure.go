package cloudamqp

import (
	"fmt"
	"strconv"

	"github.com/84codes/go-api/api"
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
				Description: "Service name of the PrivateLink, needed when creating the endpoint",
			},
			"alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "...",
			},
			"approved_subscriptions": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "...",
			},
		},
	}
}

func resourcePrivateLinkAzureCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
	)
	err := api.EnablePrivatelink(instanceID)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%d", instanceID))
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
			d.Set(k, v)
		}
	}
	return nil
}

func resourcePrivateLinkAzureUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		params     = d.Get("allowed_principals").(map[string]interface{})
	)
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
		"service_name",
		"alias",
		"approved_subscriptions":
		return true
	}
	return false
}

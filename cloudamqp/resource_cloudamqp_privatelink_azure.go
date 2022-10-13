package cloudamqp

import (
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

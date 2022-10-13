package cloudamqp

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePrivateLinkAws() *schema.Resource {
	return &schema.Resource{
		Create: resourcePrivateLinkAwsCreate,
		Read:   resourcePrivateLinkAwsRead,
		Update: resourcePrivateLinkAwsUpdate,
		Delete: resourcePrivateLinkAwsDelete,
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
			"allowed_principals": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Optional:    true,
				Description: "...",
			},
			"active_zones": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "...",
			},
		},
	}
}

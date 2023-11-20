package cloudamqp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVpcConnect() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcConnectCreate,
		Read:   resourceVpcConnectRead,
		Update: resourceVpcConnectUpdate,
		Delete: resourceVpcConnectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The CloudAMQP instance identifier",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region where the CloudAMQP instance is hosted",
			},
			"allowlist": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "List of allowed prinicpals, projects or subscriptions depending on platform provider",
			},
			"active_zones": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Covering availability zones used when creating an endpoint from other VPC. [AWS]",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the Private Service Connect [enabled, pending, disabled]",
			},
			"service_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service name (alias for Azure) of the PrivateLink.",
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
	}
}

func resourceVpcConnectCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		region     = d.Get("region").(string)
		allowlist  = d.Get("allowlist").([]interface{})
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
		params     = make(map[string][]interface{})
	)

	switch getPlatform(region) {
	case "amazon":
		params["allowed_principals"] = allowlist
	case "azure":
		params["approved_subscriptions"] = allowlist
	case "google":
		params["allowed_projects"] = allowlist
	default:
		return fmt.Errorf("invalid region")
	}

	err := api.EnablePrivatelink(instanceID, params, sleep, timeout)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", instanceID))
	return resourceVpcConnectRead(d, meta)
}

func resourceVpcConnectRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api           = meta.(*api.API)
		instanceID, _ = strconv.Atoi(d.Id()) // Uses d.Id() to allow import
	)

	data, err := api.ReadPrivatelink(instanceID)
	if err != nil {
		return err
	}

	for k, v := range data {
		if validateVpcConnectSchemaAttribute(k) {
			if k == "alias" {
				d.Set("service_name", v)
			} else {
				d.Set(k, v)
			}
		}
	}
	return nil
}

func resourceVpcConnectUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		region     = d.Get("region").(string)
		allowlist  = d.Get("allowlist").([]interface{})
		params     = make(map[string][]interface{})
	)

	switch getPlatform(region) {
	case "amazon":
		params["allowed_principals"] = allowlist
	case "azure":
		params["approved_subscriptions"] = allowlist
	case "google":
		params["allowed_projects"] = allowlist
	default:
		return fmt.Errorf("invalid region")
	}

	err := api.UpdatePrivatelink(instanceID, params)
	if err != nil {
		return err
	}
	return nil
}

func resourceVpcConnectDelete(d *schema.ResourceData, meta interface{}) error {
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

func validateVpcConnectSchemaAttribute(key string) bool {
	switch key {
	case "active_zones",
		"alias",
		"allowed_principals",
		"allowed_projects",
		"approved_subscriptions",
		"service_name",
		"status":
		return true
	}
	return false
}

func getPlatform(region string) string {
	regionSplit := strings.Split(region, "::")
	switch regionSplit[0] {
	case "amazon-web-services":
		return "amazon"
	case "azure-arm":
		return "azure"
	case "google-compute-engine":
		return "google"
	default:
		return ""
	}
}

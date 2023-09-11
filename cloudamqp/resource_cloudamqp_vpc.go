package cloudamqp

import (
	"fmt"
	"log"
	"net"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcCreate,
		Read:   resourceVpcRead,
		Update: resourceVpcUpdate,
		Delete: resourceVpcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the VPC instance",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The hosted region for the standalone VPC instance",
			},
			"subnet": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, _, err := net.ParseCIDR(v)
					if err != nil {
						errs = append(errs, fmt.Errorf("subnet: %v", err))
					}
					return
				},
				Description: "The VPC subnet",
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Tag the VPC instance with optional tags",
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VPC name given when hosted at the cloud provider",
			},
		},
	}
}

func resourceVpcCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"name", "region", "subnet", "tags"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil && v != "" {
			params[k] = v
		}
	}

	data, err := api.CreateVpcInstance(params)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] cloudamqp::vpc::create data: %v", data)
	d.SetId(data["id"].(string))
	return resourceVpcRead(d, meta)
}

func resourceVpcRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data, err := api.ReadVpcInstance(d.Id())

	if err != nil {
		return err
	}

	for k, v := range data {
		if validateVpcSchemaAttribute(k) {
			err = d.Set(k, v)
			if err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
	}

	return nil
}

func resourceVpcUpdate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	keys := []string{"name", "tags"}
	params := make(map[string]interface{})
	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = d.Get(k)
		}
	}

	if err := api.UpdateVpcInstance(d.Id(), params); err != nil {
		return err
	}

	return resourceVpcRead(d, meta)
}

func resourceVpcDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	return api.DeleteVpcInstance(d.Id())
}

func validateVpcSchemaAttribute(key string) bool {
	switch key {
	case "name",
		"region",
		"subnet",
		"tags",
		"vpc_name":
		return true
	}
	return false
}

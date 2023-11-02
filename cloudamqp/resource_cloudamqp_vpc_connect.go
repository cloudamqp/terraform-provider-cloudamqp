package cloudamqp

import (
	"fmt"
	"strconv"

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
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the Private Service Connect [enabled, pending, disabled]",
			},
			"service_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service name of the Private Service Connect",
			},
			"allowed_projects": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Only give access to allowed GCP projects",
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
		sleep      = d.Get("sleep").(int)
		timeout    = d.Get("timeout").(int)
		params     = make(map[string][]interface{})
	)

	params["allowed_projects"] = d.Get("allowed_projects").([]interface{})
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
			d.Set(k, v)
		}
	}
	return nil
}

func resourceVpcConnectUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
		params     = make(map[string][]interface{})
	)

	params["allowed_projects"] = d.Get("allowed_projects").([]interface{})
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
	case "status",
		"service_name",
		"allowed_projects":
		return true
	}
	return false
}

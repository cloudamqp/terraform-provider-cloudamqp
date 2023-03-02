package cloudamqp

import (
	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAwsEventBridge() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsEventBridgeCreate,
		Read:   resourceAwsEventBridgeRead,
		Delete: resourceAwsEventBridgeDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "Instance identifier",
			},
			"aws_account_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "",
			},
			"aws_region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "",
			},
			"vhost": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "",
			},
			"queue": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "",
			},
			"with_headers": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Required:    true,
				Description: "",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
		},
	}
}

func resourceAwsEventBridgeCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		keys       = alarmAttributeKeys()
		params     map[string]interface{}
		instanceID = d.Get("instance_id").(int)
	)

	for _, k := range keys {
		if v := d.Get(k); v != nil {
			params[k] = v
		}
	}

	data, err := api.CreateAwsEventBridge(instanceID, params)
	if err != nil {
		return err
	}

	d.SetId(data["id"].(string))
	return resourceAlarmRead(d, meta)
}

func resourceAwsEventBridgeRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api           = meta.(*api.API)
		instanceID    = d.Get("instance_id").(int)
		eventbridgeID = d.Get("id").(string)
	)

	data, err := api.ReadAwsEventBridge(instanceID, eventbridgeID)
	if err != nil {
		return err
	}
	if data["status"] != nil {
		d.Set("status", data["status"].(string))
	}
	return nil
}

func resourceAwsEventBridgeUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAwsEventBridgeDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		api           = meta.(*api.API)
		instanceID    = d.Get("instance_id").(int)
		eventbridgeID = d.Get("id").(string)
	)

	return api.DeleteAwsEventBridge(instanceID, eventbridgeID)
}

func awsEventbridgeAttributeKeys() []string {
	return []string{
		"aws_account_id",
		"aws_region",
		"vhost",
		"queue",
		"with_headers",
	}
}

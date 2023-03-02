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
				Description: "The 12 digit AWS Account ID where you want the events to be sent to.",
			},
			"aws_region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The AWS region where you the events to be sent to. (e.g. us-west-1, us-west-2, ..., etc.)",
			},
			"vhost": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The VHost the queue resides in.",
			},
			"queue": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "A (durable) queue on your RabbitMQ instance.",
			},
			"with_headers": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Required:    true,
				Description: "Include message headers in the event data.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status for the EventBridge.",
			},
		},
	}
}

func resourceAwsEventBridgeCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		keys       = awsEventbridgeAttributeKeys()
		params     = make(map[string]interface{})
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
	return nil
}

func resourceAwsEventBridgeRead(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
	)

	data, err := api.ReadAwsEventBridge(instanceID, d.Id())
	if err != nil {
		return err
	}
	if data["status"] != nil {
		d.Set("status", data["status"].(string))
	}
	return nil
}

func resourceAwsEventBridgeDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
	)

	return api.DeleteAwsEventBridge(instanceID, d.Id())
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

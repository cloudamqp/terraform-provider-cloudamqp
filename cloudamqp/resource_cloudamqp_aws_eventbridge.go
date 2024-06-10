package cloudamqp

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAwsEventBridge() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsEventBridgeCreate,
		Read:   resourceAwsEventBridgeRead,
		Delete: resourceAwsEventBridgeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
				Description: "Always set to null, unless there is an error starting the EventBridge",
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
	if strings.Contains(d.Id(), ",") {
		log.Printf("[DEBUG] cloudamqp::resource::aws-eventbridge::read id contains : %v", d.Id())
		s := strings.Split(d.Id(), ",")
		log.Printf("[DEBUG] cloudamqp::resource::aws-eventbridge::read split ids: %v, %v", s[0], s[1])
		d.SetId(s[0])
		instanceID, _ := strconv.Atoi(s[1])
		d.Set("instance_id", instanceID)
	}
	if d.Get("instance_id").(int) == 0 {
		return fmt.Errorf("missing instance identifier: {resource_id},{instance_id}")
	}

	var (
		api        = meta.(*api.API)
		instanceID = d.Get("instance_id").(int)
	)

	log.Printf("[DEBUG] cloudamqp::resource::aws-eventbridge::read ID: %v, instanceID %v", d.Id(), instanceID)
	data, err := api.ReadAwsEventBridge(instanceID, d.Id())
	if err != nil {
		return err
	}

	for k, v := range data {
		if validateAwsEventBridgeSchemaAttribute(k) {
			if v == nil {
				continue
			}
			if err = d.Set(k, v); err != nil {
				return fmt.Errorf("error setting %s for resource %s: %s", k, d.Id(), err)
			}
		}
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

func validateAwsEventBridgeSchemaAttribute(key string) bool {
	switch key {
	case "aws_account_id",
		"aws_region",
		"vhost",
		"queue",
		"with_headers",
		"status":
		return true
	}
	return false
}

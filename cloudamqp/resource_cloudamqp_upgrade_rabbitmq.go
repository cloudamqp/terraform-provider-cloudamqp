package cloudamqp

import (
	"log"
	"strconv"

	"github.com/cloudamqp/terraform-provider-cloudamqp/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUpgradeRabbitMQ() *schema.Resource {
	return &schema.Resource{
		Create: resourceUpgradeRabbitMQInvoke,
		Read:   resourceUpgradeRabbitMQRead,
		Update: resourceUpgradeRabbitMQUpdate,
		Delete: resourceUpgradeRabbitMQRemove,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The CloudAMQP instance identifier",
			},
			"current_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Helper argument to change upgrade behaviour to latest possible version",
			},
			"new_version": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "The new version to upgrade to",
			},
		},
	}
}

func resourceUpgradeRabbitMQInvoke(d *schema.ResourceData, meta interface{}) error {
	var (
		api             = meta.(*api.API)
		instanceID      = d.Get("instance_id").(int)
		current_version = d.Get("current_version").(string)
		new_version     = d.Get("new_version").(string)
	)

	log.Printf("[DEBUG] - Upgrading RabbitMQ instance %d to version %s", instanceID, new_version)
	response, err := api.UpgradeRabbitMQ(instanceID, current_version, new_version)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(instanceID))

	if len(response) > 0 {
		log.Println("[INFO] - ", response)
	}

	return nil
}

func resourceUpgradeRabbitMQRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceUpgradeRabbitMQUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceUpgradeRabbitMQRemove(d *schema.ResourceData, meta interface{}) error {
	return nil
}

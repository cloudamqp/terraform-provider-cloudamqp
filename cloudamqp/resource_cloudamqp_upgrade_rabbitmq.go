package cloudamqp

import (
	"log"
	"strconv"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUpgradeRabbitMQ() *schema.Resource {
	return &schema.Resource{
		Create: resourceUpgradeRabbitMQInvoke,
		Read:   resourceUpgradeRabbitMQRead,
		Delete: resourceUpgradeRabbitMQRemove,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
				Description: "The CloudAMQP instance identifier",
			},
		},
	}
}

func resourceUpgradeRabbitMQInvoke(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	response, err := api.UpgradeRabbitMQ(d.Get("instance_id").(int))
	if err != nil {
		return err
	}
	id := strconv.Itoa(d.Get("instance_id").(int))
	d.SetId(id)

	if len(response) > 0 {
		log.Printf("[INFO] - " + response)
	}

	return nil
}

func resourceUpgradeRabbitMQRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceUpgradeRabbitMQRemove(d *schema.ResourceData, meta interface{}) error {
	return nil
}
